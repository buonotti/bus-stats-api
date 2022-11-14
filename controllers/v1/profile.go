package v1

import (
	"net/http"

	"github.com/buonotti/bus-stats-api/jwt"
	"github.com/buonotti/bus-stats-api/models"
	"github.com/buonotti/bus-stats-api/services"
	serviceV1 "github.com/buonotti/bus-stats-api/services/v1"
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
	userId := jwt.ExtractUidFromHeader(c)
	fileForm, err := c.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, services.ErrorResponse{Message: err.Error()})
		return
	}

	result, status, err := serviceV1.SaveUserProfile(models.UserId(userId), fileForm)
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
	userId := jwt.ExtractUidFromHeader(c)

	result, status, err := serviceV1.GetUserProfile(models.UserId(userId))
	if err != nil {
		c.AbortWithStatusJSON(status, services.ErrorResponse{Message: err.Error()})
		return
	}

	c.File(result.FileName)
}
