package models

import (
	"encoding/json"
	"testing"
)

func TestProduct_JSONSerialization(t *testing.T) {
	product := Product{
		ID:    "test-product",
		Price: 99.99,
		Stock: 50,
	}

	// Test marshaling
	jsonData, err := json.Marshal(product)
	if err != nil {
		t.Errorf("Failed to marshal product: %v", err)
	}

	// Test unmarshaling
	var unmarshaledProduct Product
	err = json.Unmarshal(jsonData, &unmarshaledProduct)
	if err != nil {
		t.Errorf("Failed to unmarshal product: %v", err)
	}

	if unmarshaledProduct.ID != product.ID {
		t.Errorf("Expected ID %s, got %s", product.ID, unmarshaledProduct.ID)
	}
	if unmarshaledProduct.Price != product.Price {
		t.Errorf("Expected Price %.2f, got %.2f", product.Price, unmarshaledProduct.Price)
	}
	if unmarshaledProduct.Stock != product.Stock {
		t.Errorf("Expected Stock %d, got %d", product.Stock, unmarshaledProduct.Stock)
	}
}

func TestProductEvent_JSONSerialization(t *testing.T) {
	event := ProductEvent{
		ProductID: "test-product",
		Price:     99.99,
		Stock:     50,
	}

	// Test marshaling
	jsonData, err := json.Marshal(event)
	if err != nil {
		t.Errorf("Failed to marshal product event: %v", err)
	}

	// Test unmarshaling
	var unmarshaledEvent ProductEvent
	err = json.Unmarshal(jsonData, &unmarshaledEvent)
	if err != nil {
		t.Errorf("Failed to unmarshal product event: %v", err)
	}

	if unmarshaledEvent.ProductID != event.ProductID {
		t.Errorf("Expected ProductID %s, got %s", event.ProductID, unmarshaledEvent.ProductID)
	}
	if unmarshaledEvent.Price != event.Price {
		t.Errorf("Expected Price %.2f, got %.2f", event.Price, unmarshaledEvent.Price)
	}
	if unmarshaledEvent.Stock != event.Stock {
		t.Errorf("Expected Stock %d, got %d", event.Stock, unmarshaledEvent.Stock)
	}
}

func TestHealthResponse_JSONSerialization(t *testing.T) {
	response := HealthResponse{
		Status: "healthy",
	}

	// Test marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal health response: %v", err)
	}

	// Test unmarshaling
	var unmarshaledResponse HealthResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal health response: %v", err)
	}

	if unmarshaledResponse.Status != response.Status {
		t.Errorf("Expected Status %s, got %s", response.Status, unmarshaledResponse.Status)
	}
}

func TestErrorResponse_JSONSerialization(t *testing.T) {
	response := ErrorResponse{
		Error: "test error message",
	}

	// Test marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal error response: %v", err)
	}

	// Test unmarshaling
	var unmarshaledResponse ErrorResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal error response: %v", err)
	}

	if unmarshaledResponse.Error != response.Error {
		t.Errorf("Expected Error %s, got %s", response.Error, unmarshaledResponse.Error)
	}
}

func TestEventResponse_JSONSerialization(t *testing.T) {
	response := EventResponse{
		Message:   "Event accepted for processing",
		ProductID: "test-product",
	}

	// Test marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal event response: %v", err)
	}

	// Test unmarshaling
	var unmarshaledResponse EventResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal event response: %v", err)
	}

	if unmarshaledResponse.Message != response.Message {
		t.Errorf("Expected Message %s, got %s", response.Message, unmarshaledResponse.Message)
	}
	if unmarshaledResponse.ProductID != response.ProductID {
		t.Errorf("Expected ProductID %s, got %s", response.ProductID, unmarshaledResponse.ProductID)
	}
}

func TestProduct_ZeroValues(t *testing.T) {
	var product Product

	if product.ID != "" {
		t.Errorf("Expected empty ID, got %s", product.ID)
	}
	if product.Price != 0.0 {
		t.Errorf("Expected Price 0.0, got %.2f", product.Price)
	}
	if product.Stock != 0 {
		t.Errorf("Expected Stock 0, got %d", product.Stock)
	}
}

func TestProductEvent_ZeroValues(t *testing.T) {
	var event ProductEvent

	if event.ProductID != "" {
		t.Errorf("Expected empty ProductID, got %s", event.ProductID)
	}
	if event.Price != 0.0 {
		t.Errorf("Expected Price 0.0, got %.2f", event.Price)
	}
	if event.Stock != 0 {
		t.Errorf("Expected Stock 0, got %d", event.Stock)
	}
}

func TestHealthResponse_ZeroValues(t *testing.T) {
	var response HealthResponse

	if response.Status != "" {
		t.Errorf("Expected empty Status, got %s", response.Status)
	}
}

func TestErrorResponse_ZeroValues(t *testing.T) {
	var response ErrorResponse

	if response.Error != "" {
		t.Errorf("Expected empty Error, got %s", response.Error)
	}
}

func TestEventResponse_ZeroValues(t *testing.T) {
	var response EventResponse

	if response.Message != "" {
		t.Errorf("Expected empty Message, got %s", response.Message)
	}
	if response.ProductID != "" {
		t.Errorf("Expected empty ProductID, got %s", response.ProductID)
	}
}
