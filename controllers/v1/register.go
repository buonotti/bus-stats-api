package v1

import (
	"net/http"

	"github.com/buonotti/bus-stats-api/services"
	serviceV1 "github.com/buonotti/bus-stats-api/services/v1"
	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary Register user
// @Description Register a user. The response contains a token used for authentication
// @ID register-user
// @Tags authentication
// @Schemes serviceV1.RegisterRequest serviceV1.RegisterResponse
// @Accept application/json
// @Produce application/json
// @Param data body v1.RegisterRequest true "content"
// @Success 200 {object} v1.RegisterResponse
// @Failure 400 {object} services.ErrorResponse
// @Failure 401 {object} services.ErrorResponse
// @Failure 500 {object} services.ErrorResponse
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	var request serviceV1.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, services.ErrorResponse{Message: err.Error()})
		return
	}

	_, err, status := serviceV1.RegisterUser(request)

	if err != nil {
		c.AbortWithStatusJSON(status, services.ErrorResponse{Message: err.Error()})
		return
	}

	c.Status(status)
}
