package uroncha

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"github.com/zjyl1994/caasiu"
)

type Uroncha struct{}

var Logger *logrus.Logger
var router *gin.Engine
var debugMode bool
var httpPort = ":8080"

func init() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.000",
		DisableTimestamp: false,
		FullTimestamp:    true,
	})
	debugMode = strings.EqualFold(os.Getenv("URONCHA_DEBUG"), "true")
	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.Default()
	router.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
}

func Run() error {
	return router.Run(httpPort)
}

func IsDebug() bool {
	return debugMode
}

func Handle(method, url string, validRules Rules, handler HandleFunc) {
	router.Handle(method, url, func(c *gin.Context) {
		timestamp := time.Now().Unix()
		cv, err := caasiu.New(c.Request)
		if err != nil {
			Logger.Errorln("failed to init caasiu", err.Error())
			c.AbortWithStatus(500)
			return
		}
		params := Datas{
			QueryString: cv.QueryString().Data(),
			Body:        cv.Body().Data(),
		}
		ret, uerr := handler(c, params)
		c.JSON(200, gin.H{
			"success":   uerr.Success,
			"code":      uerr.Code,
			"message":   uerr.Message,
			"result":    ret,
			"timestamp": timestamp,
		})
	})
}
