package middlewares

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(c *gin.Context) {
	logger := log.New()

	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel
	}

	logger.SetLevel(logLevel)
	logger.SetFormatter(&log.JSONFormatter{})

	startTime := time.Now()

	c.Next()

	latency := time.Since(startTime)

	logger.WithFields(log.Fields{
		"METHOD":    c.Request.Method,
		"URI":       c.Request.RequestURI,
		"STATUS":    c.Writer.Status(),
		"LATENCY":   latency,
		"CLIENT_IP": c.ClientIP(),
	}).Info("HTTP REQUEST")
}
