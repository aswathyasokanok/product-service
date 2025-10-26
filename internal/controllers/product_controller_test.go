package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"product-service/internal/models"
	"product-service/internal/repositories"
	"product-service/internal/services"
	"product-service/pkg/queue"

	"github.com/gin-gonic/gin"
)

func TestProductController(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create real dependencies for testing
	repo := repositories.NewInMemoryProductRepository()
	eventQueue := queue.NewInMemoryEventQueue(100)
	productService := services.NewProductService(repo, eventQueue, 1)

	controller := NewProductController(productService)

	// Create a test router
	router := gin.New()
	router.POST("/events", controller.HandleEvent)
	router.GET("/products/:id", controller.GetProduct)

	// Test POST /events
	t.Run("HandleEvent", func(t *testing.T) {
		event := models.ProductEvent{ProductID: "test-product", Price: 25.0, Stock: 15}
		eventJSON, _ := json.Marshal(event)

		req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(eventJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusAccepted {
			t.Errorf("Expected status 202, got %d", w.Code)
		}
	})

	// Test invalid JSON
	t.Run("HandleEvent_InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/events", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	// Test missing product_id
	t.Run("HandleEvent_MissingProductID", func(t *testing.T) {
		event := models.ProductEvent{Price: 10.0, Stock: 5} // Missing ProductID
		eventJSON, _ := json.Marshal(event)

		req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(eventJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	// Test GET /products/{id} - product exists
	t.Run("GetProduct_Exists", func(t *testing.T) {
		// First create a product by processing an event
		event := models.ProductEvent{ProductID: "get-test", Price: 50.0, Stock: 25}
		eventJSON, _ := json.Marshal(event)

		req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(eventJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusAccepted {
			t.Errorf("Expected status 202, got %d", w.Code)
		}

		// Wait for async processing
		time.Sleep(100 * time.Millisecond)

		// Now get the product
		req, _ = http.NewRequest("GET", "/products/get-test", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var product models.Product
		if err := json.Unmarshal(w.Body.Bytes(), &product); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		if product.ID != "get-test" || product.Price != 50.0 || product.Stock != 25 {
			t.Errorf("Expected product{ID: get-test, Price: 50.0, Stock: 25}, got %+v", product)
		}
	})

	// Test GET /products/{id} - product not found
	t.Run("GetProduct_NotFound", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products/nonexistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}

		var errorResp models.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &errorResp); err != nil {
			t.Errorf("Failed to unmarshal error response: %v", err)
		}

		if errorResp.Error != "Product not found" {
			t.Errorf("Expected error 'Product not found', got '%s'", errorResp.Error)
		}
	})

	// Test queue full scenario
	t.Run("HandleEvent_QueueFull", func(t *testing.T) {
		// Create a small queue to test queue full scenario
		smallQueue := queue.NewInMemoryEventQueue(1)
		smallService := services.NewProductService(repo, smallQueue, 1)
		smallController := NewProductController(smallService)

		router := gin.New()
		router.POST("/events", smallController.HandleEvent)

		// Fill the queue
		event1 := models.ProductEvent{ProductID: "queue1", Price: 1.0, Stock: 1}
		event1JSON, _ := json.Marshal(event1)
		req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(event1JSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusAccepted {
			t.Errorf("Expected status 202 for first event, got %d", w.Code)
		}

		// Try to add another event (should fail due to queue full)
		event2 := models.ProductEvent{ProductID: "queue2", Price: 2.0, Stock: 2}
		event2JSON, _ := json.Marshal(event2)
		req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(event2JSON))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status 503 for queue full, got %d", w.Code)
		}
	})
}
