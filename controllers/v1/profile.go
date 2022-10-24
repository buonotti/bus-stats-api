package v1

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/buonotti/bus-stats-api/middleware"
	"github.com/buonotti/bus-stats-api/services"
	serviceV1 "github.com/buonotti/bus-stats-api/services/v1"
	"github.com/buonotti/bus-stats-api/util"
	"github.com/gin-gonic/gin"
)

// UploadUserProfilePicture godoc
// @Summary Upload user profile picture
// @Description Upload a user profile picture in a form for the currently logged-in user
// @ID upload-profile
// @Tags user-account
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Schemes v1.SaveUserProfileResponse services.ErrorResponse
// @Accept multipart/form-data
// @Produce application/json
// @Param id path string true "user id"
// @Param image formData file true "picture form data"
// @Success 200 {object} v1.SaveUserProfileResponse
// @Failure 400 {object} services.ErrorResponse
// @Failure 401 {object} services.ErrorResponse
// @Failure 500 {object} services.ErrorResponse
// @Router /profile [post]
func UploadUserProfilePicture(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userId := claims[middleware.IdentityKey].(string)
	fileForm, err := c.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, services.ErrorResponse{Message: err.Error()})
		return
	}

	result, err, status := serviceV1.SaveUserProfile(util.ExtractToken(c), userId, fileForm)
	if err != nil {
		c.AbortWithStatusJSON(status, services.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(status, result)
}

// GetUserProfile godoc
// @Summary Get the profile picture a user
// @Description Get the profile picture file for the currently authenticated user
// @ID get-profile
// @Tags user-account
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Schemes services.ErrorResponse
// @Produce multipart/form-data
// @Success 200 binary formData
// @Failure 400 {object} services.ErrorResponse
// @Failure 401 {object} services.ErrorResponse
// @Failure 500 {object} services.ErrorResponse
// @Router /profile [get]
func GetUserProfile(c *gin.Context) {
	userId := jwt.ExtractClaims(c)[middleware.IdentityKey].(string)
	result, err, status := serviceV1.GetUserProfile(userId)
	if err != nil {
		c.AbortWithStatusJSON(status, services.ErrorResponse{Message: err.Error()})
	}

	c.Status(status)
	c.File(result.FileName)
}