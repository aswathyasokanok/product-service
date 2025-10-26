package v1

import (
	"product-service/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures the API routes
func SetupRoutes(router *gin.Engine, productController *controllers.ProductController, healthController *controllers.HealthController) {
	// Health check
	router.GET("/health", healthController.Health)

	// API v1 routes
	api := router.Group("/api/v1")
	{
		api.POST("/events", productController.HandleEvent)
		api.GET("/products/:id", productController.GetProduct)
	}
}
