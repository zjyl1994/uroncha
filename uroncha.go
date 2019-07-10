package uroncha

import (
	"encoding/json"
	"net/http"
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

func init() {
	debugMode = strings.EqualFold(utils.MustGetenvStr("URONCHA_DEBUG", "false"), "true")
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.000",
		DisableTimestamp: false,
		FullTimestamp:    true,
	})
	router = httprouter.New()
	timeoutHandler := http.TimeoutHandler(router, time.Duration(utils.MustGetenvInt("URONCHA_HANDLER_TIMEOUT", 3)), "")
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
		timestamp := time.Now().Unix()
		cv, err := caasiu.New(r)
		if err != nil {
			Logger.Errorln("failed to init caasiu", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
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
			Logger.Errorln("failed to marshal json", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			n, err := w.Write(bjson)
			if err != nil {
				Logger.Errorln("failed to write response", err.Error(), n)
			}
		}
	})
}
