package service

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ljcbaby/hdu-wiki-qa/conf"
	"github.com/sirupsen/logrus"
)

func Serve() {
	config := conf.Service

	r := SetupRouter()

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	err := r.Run(addr)
	if err != nil {
		logrus.WithField("module", "service").Errorf("Failed to run service: %v", err)
	}
}

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.Default()

	r.Use(LoggerMiddleware())
	r.Use(CORS())

	r.GET("/", Index)
	r.POST("/v1/chat/completions", Chat)

	return r
}

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, HDU Wiki QA!",
	})
}

func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Expose-Headers", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
