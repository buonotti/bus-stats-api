package v1

import (
	"time"

	"github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

func MapRoutes(router *gin.RouterGroup, store *persist.MemoryStore) {
	router.POST("/login",
		cache.CacheByRequestURI(store, 2*time.Second),
		LoginUser)
	router.POST("/register", RegisterUser)
}
