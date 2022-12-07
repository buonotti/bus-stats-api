package v1

import (
	"github.com/buonotti/bus-stats-api/middleware"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

func MapRoutes(router *gin.RouterGroup, store *persist.MemoryStore) {
	router.POST("/login", LoginUser)
	router.POST("/register", RegisterUser)
	router.POST("/refresh", RefreshUserToken)

	router.Use(middleware.Auth())
	{
		router.POST("/profile/:id", UploadUserProfile)
		// router.GET("/profile/:id", cache.CacheByRequestPath(store, 1*time.Minute), GetUserProfile)
		router.GET("/profile/:id", GetUserProfile)
		router.DELETE("/profile/:id", DeleteUserProfile)
	}
}
