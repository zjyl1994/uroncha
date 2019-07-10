package uroncha

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/caasiu"
	"github.com/zjyl1994/uroncha/utils"
)

type Uroncha struct{}

var Logger *logrus.Logger
var router *httprouter.Router
var srv *http.Server
var debugMode bool
var loggerUseColor bool

func init() {
	debugMode = utils.MustGetenvBool("URONCHA_DEBUG", false)
	loggerUseColor = !utils.MustGetenvBool("URONCHA_DISABLE_COLOR", false)
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.000",
		DisableTimestamp: false,
		FullTimestamp:    true,
		ForceColors:      loggerUseColor,
	})
	Logger.SetReportCaller(debugMode)
	router = httprouter.New()
	timeoutHandler := http.TimeoutHandler(router, time.Duration(utils.MustGetenvInt("URONCHA_HANDLER_TIMEOUT", 8))*time.Second, "")
	srv = &http.Server{
		Addr:         utils.MustGetenvStr("URONCHA_PORT", ":8080"),
		Handler:      timeoutHandler,
		ReadTimeout:  time.Duration(utils.MustGetenvInt("URONCHA_READ_TIMEOUT", 5)) * time.Second,
		WriteTimeout: time.Duration(utils.MustGetenvInt("URONCHA_WRITE_TIMEOUT", 10)) * time.Second,
	}
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) { w.WriteHeader(http.StatusOK) })
}

func Run() error {
	return srv.ListenAndServe()
}

func IsDebug() bool {
	return debugMode
}

func Handle(method, url string, validRules Rules, handler HandleFunc) {
	router.Handle(method, url, func(w http.ResponseWriter, r *http.Request, hrps httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				stack := stack(3)
				httpRequest, _ := httputil.DumpRequest(r, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				var errMsg string
				if brokenPipe {
					errMsg = fmt.Sprintf("\n%s\n%s", err, string(httpRequest))
				} else if debugMode {
					errMsg = fmt.Sprintf("\n%s panic recovered:\n%s%s\n%s",
						time.Now().Format("2006-01-02 15:04:05.000"), strings.Join(headers, "\r\n"), err, stack)
				} else {
					errMsg = fmt.Sprintf("\n%s panic recovered:\n%s%s",
						time.Now().Format("2006-01-02 15:04:05.000"), err, stack)
				}
				if loggerUseColor {
					errMsg = "\x1b[31m" + errMsg + "\x1b[0m"
				}
				fmt.Println(errMsg)
				if !brokenPipe {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()
		startTime := time.Now()
		timestamp := startTime.Unix()
		cv, err := caasiu.New(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			Logger.WithFields(logrus.Fields{
				"method":  method,
				"url":     url,
				"status":  http.StatusInternalServerError,
				"elapsed": time.Since(startTime).String(),
				"error":   err.Error(),
			}).Error("failed to init caasiu")
			return
		}
		validated, validMsg := cv.Valid(validRules)
		var result H
		if !validated {
			result = H{
				"success":   false,
				"code":      1,
				"message":   "param not validated",
				"result":    validMsg,
				"timestamp": timestamp,
			}
		} else {
			pathParam := make(map[string]string)
			for _, p := range hrps {
				pathParam[p.Key] = p.Value
			}
			c := &Context{
				Req:         r,
				Res:         w,
				QueryString: cv.QueryString().Data(),
				Body:        cv.Body().Data(),
				PathParam:   pathParam,
			}
			ret, uerr := handler(c)
			result = H{
				"success":   uerr.Success,
				"code":      uerr.Code,
				"message":   uerr.Message,
				"result":    ret,
				"timestamp": timestamp,
			}
		}
		bjson, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			Logger.WithFields(logrus.Fields{
				"method":  method,
				"url":     url,
				"status":  http.StatusInternalServerError,
				"elapsed": time.Since(startTime).String(),
				"error":   err.Error(),
			}).Error("failed to marshal json")
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write(bjson)
			if err != nil {
				Logger.WithFields(logrus.Fields{
					"method":  method,
					"url":     url,
					"status":  http.StatusInternalServerError,
					"elapsed": time.Since(startTime).String(),
					"error":   err.Error(),
				}).Error("failed to write response")
			} else {
				Logger.WithFields(logrus.Fields{
					"method":  method,
					"url":     url,
					"status":  http.StatusOK,
					"elapsed": time.Since(startTime).String(),
				}).Info(url)
			}
		}
	})
}
