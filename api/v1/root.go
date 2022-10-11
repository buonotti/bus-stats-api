package v1

import (
	"time"

	"github.com/buonotti/bus-stats-api/api/v1/hello"
	"github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

func MapRoutes(router *gin.RouterGroup, store *persist.MemoryStore) {
	router.GET("/hello",
		cache.CacheByRequestURI(store, 2*time.Second),
		hello.Endpoint)
}
