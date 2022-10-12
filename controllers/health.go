package controllers

import (
	"github.com/buonotti/bus-stats-api/services"
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// HealthEndpoint godoc
// @Summary Health check
// @Schemes
// @Description Get the health status of the API
// @Accept json
// @Produce json
// @Success 200 {string} HealthStatus
// @Router /health [get]
func HealthEndpoint(c *gin.Context) {
	data := services.DbInfo()

	c.JSON(200, gin.H{"status": data.Status})
}