package controllers

import (
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// HealthEndpoint godoc
// @Summary Health check
// @Description Get the health status of the API
// @ID health
// @Tags health
// @Success 200
// @Router /health [get]
func HealthEndpoint(c *gin.Context) {
	c.Status(200)
}