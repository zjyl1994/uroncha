package main

import (
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05.000",
		DisableTimestamp: false,
		FullTimestamp:    true,
	})
	logger.Println("Uroncha Test Server")
	router := gin.Default()
	router.GET("/test", TestHandler)
	err := router.Run(":8080")
	logrus.Println("Server", err)
}

func TestHandler(c *gin.Context) {
	c.JSON(200, os.Getenv("TEST"))
}
