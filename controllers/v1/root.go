package v1

import (
	"github.com/buonotti/bus-stats-api/middleware"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

func MapRoutes(router *gin.RouterGroup, store *persist.MemoryStore) {
	router.POST("/login", middleware.Auth().LoginHandler)
	router.POST("/logout", middleware.Auth().LogoutHandler)
	router.POST("/register", RegisterUser)
	router.POST("/refresh", middleware.Auth().RefreshHandler)
	router.Use(middleware.Auth().MiddlewareFunc())
	{
		router.POST("/profile", UploadUserProfilePicture)
		router.GET("/profile", GetUserProfile)
	}
}
