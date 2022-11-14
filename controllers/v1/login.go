package v1

import (
	"net/http"

	"github.com/buonotti/bus-stats-api/services"
	serviceV1 "github.com/buonotti/bus-stats-api/services/v1"
	"github.com/gin-gonic/gin"
)

// LoginUser godoc
// @Summary Logs a user in
// @Schemes services.ErrorResponse serviceV1.LoginRequest
// @Description Logs a user in using the provided credentials
// @ID login-user
// @Tags authentication
// @Accept application/json
// @Produces application/json
// @Param data body v1.LoginRequest true "content"
// @Success 200 {object} v1.LoginResponse
// @Failure 400 {object} services.ErrorResponse
// @Failure 401 {object} services.ErrorResponse
// @Failure 500 {object} services.ErrorResponse
// @Router /login [post]
func LoginUser(c *gin.Context) {
	var request serviceV1.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, services.ErrorResponse{Message: err.Error()})
		return
	}

	response, status, err := serviceV1.LoginUser(request)
	if err != nil {
		c.AbortWithStatusJSON(status, services.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(status, response)
}
