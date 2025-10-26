package controllers

import (
	"net/http"

	"product-service/internal/models"
	"product-service/internal/services"

	"github.com/gin-gonic/gin"
)

// ProductController handles HTTP requests for products
type ProductController struct {
	productService *services.ProductService
}

// NewProductController creates a new product controller
func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// HandleEvent handles POST /events
func (pc *ProductController) HandleEvent(c *gin.Context) {
	var event models.ProductEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON payload"})
		return
	}

	// Validate required fields
	if event.ProductID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "product_id is required"})
		return
	}

	// Process the event
	if err := pc.productService.ProcessEvent(event); err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{Error: "Queue is full"})
		return
	}

	c.JSON(http.StatusAccepted, models.EventResponse{
		Message:   "Event accepted for processing",
		ProductID: event.ProductID,
	})
}

// GetProduct handles GET /products/{id}
func (pc *ProductController) GetProduct(c *gin.Context) {
	productID := c.Param("id")

	product, exists := pc.productService.GetProduct(productID)
	if !exists {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
