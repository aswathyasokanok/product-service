package models

// Product represents a product with its current state
type Product struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// ProductEvent represents an incoming product update event
type ProductEvent struct {
	ProductID string  `json:"product_id"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// EventResponse represents the response after accepting an event
type EventResponse struct {
	Message   string `json:"message"`
	ProductID string `json:"product_id"`
}
