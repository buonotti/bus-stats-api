package hello

import (
	"github.com/buonotti/bus-stats-api/services/v1/hello"
	"github.com/gin-gonic/gin"
)

func Endpoint(c *gin.Context) {
	data := hello.Result()

	c.JSON(200, gin.H{"message": data})
}
