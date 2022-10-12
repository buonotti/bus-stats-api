package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func LogrusLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.FullPath()
		if url == "" {
			url = "404"
		}
		t := time.Now()
		c.Next()
		elapsed := time.Since(t)
		ApiLogger.WithFields(logrus.Fields{"time":fmt.Sprintf("%dms", elapsed.Milliseconds()), "status":c.Writer.Status()}).Info(url)
	}
}