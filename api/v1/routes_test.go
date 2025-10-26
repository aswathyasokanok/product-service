package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup router with nil controllers to test route registration
	router := gin.New()
	SetupRoutes(router, nil, nil)

	t.Run("HealthRoute", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 500 because controllers are nil
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500 for nil controller, got %d", w.Code)
		}
	})

	t.Run("EventsRoute", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/events", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 500 because controllers are nil
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500 for nil controller, got %d", w.Code)
		}
	})

	t.Run("ProductsRoute", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/products/test-id", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 500 because controllers are nil
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500 for nil controller, got %d", w.Code)
		}
	})

	t.Run("InvalidRoute", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 404 for invalid route
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 for invalid route, got %d", w.Code)
		}
	})
}

func TestSetupRoutes_WithRealControllers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// This test would require proper setup of real controllers
	// For now, we just test that the function doesn't panic
	router := gin.New()

	// Test with nil controllers (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetupRoutes panicked with nil controllers: %v", r)
		}
	}()

	SetupRoutes(router, nil, nil)
}
