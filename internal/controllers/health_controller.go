package controllers

import (
	"net/http"

	"product-service/internal/models"

	"github.com/gin-gonic/gin"
)

// HealthController handles health check requests
type HealthController struct{}

// NewHealthController creates a new health controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Health handles GET /health
func (hc *HealthController) Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{Status: "healthy"})
}
