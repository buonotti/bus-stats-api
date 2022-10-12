package v1

import (
	serviceV1 "github.com/buonotti/bus-stats-api/services/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// RegisterUser godoc
// @Summary Register user
// @Schemes serviceV1.RegisterRequest
// @Description Register a user. The response contains a token used for authentication
// @Accept json
// @Produce json
// @Success 200 {serviceV1.RegisterResponse} RegisterResponse
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	var request serviceV1.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := serviceV1.RegisterUser(request)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"result": "ERROR", "message": err.Error()})
		return
	}

	c.JSON(200, response)
}
