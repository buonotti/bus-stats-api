package v1

import (
	"time"

	"github.com/buonotti/bus-stats-api/middleware"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

func MapRoutes(router *gin.RouterGroup, store *persist.MemoryStore) {
	router.POST("/login", LoginUser)
	router.POST("/register", RegisterUser)
	router.POST("/refresh", RefreshUserToken)

	router.Use(middleware.Auth())
	{
		router.POST("/profile/:id", cache.CacheByRequestPath(store, 1*time.Minute), UploadUserProfilePicture)
		router.GET("/profile/:id", GetUserProfile)
	}
}
