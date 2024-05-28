package service

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms", stopTime.Nanoseconds()/1000000)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method
		url := c.Request.RequestURI
		Log := logrus.WithFields(logrus.Fields{
			"SpendTime": spendTime,
			"path":      url,
			"Method":    method,
			"status":    statusCode,
			"Ip":        clientIP,
		})
		if len(c.Errors) > 0 {
			Log.Error(c.Errors.ByType(gin.ErrorTypePrivate))
		}
		if statusCode >= 500 {
			Log.Error()
		} else if statusCode >= 400 {
			Log.Warn()
		} else {
			Log.Info()
		}
	}
}
