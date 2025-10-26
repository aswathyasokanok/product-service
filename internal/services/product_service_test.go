package services

import (
	"errors"
	"testing"
	"time"

	"product-service/internal/models"
)

// MockProductRepository for testing
type MockProductRepository struct {
	products map[string]*models.Product
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products: make(map[string]*models.Product),
	}
}

func (m *MockProductRepository) Get(id string) (*models.Product, bool) {
	product, exists := m.products[id]
	return product, exists
}

func (m *MockProductRepository) Update(id string, price float64, stock int) {
	m.products[id] = &models.Product{
		ID:    id,
		Price: price,
		Stock: stock,
	}
}

// MockEventQueue for testing
type MockEventQueue struct {
	events chan models.ProductEvent
	closed bool
}

func NewMockEventQueue(bufferSize int) *MockEventQueue {
	return &MockEventQueue{
		events: make(chan models.ProductEvent, bufferSize),
		closed: false,
	}
}

func (m *MockEventQueue) Enqueue(event models.ProductEvent) error {
	if m.closed {
		return errors.New("queue is closed")
	}
	select {
	case m.events <- event:
		return nil
	default:
		return errors.New("queue is full")
	}
}

func (m *MockEventQueue) Dequeue() (models.ProductEvent, bool) {
	select {
	case event, ok := <-m.events:
		return event, ok
	default:
		return models.ProductEvent{}, false
	}
}

func (m *MockEventQueue) Close() {
	close(m.events)
	m.closed = true
}

func TestProductService_ProcessEvent(t *testing.T) {
	repo := NewMockProductRepository()
	eventQueue := NewMockEventQueue(10)
	service := NewProductService(repo, eventQueue, 1)

	t.Run("ProcessEvent_Success", func(t *testing.T) {
		event := models.ProductEvent{ProductID: "test-1", Price: 10.0, Stock: 5}

		err := service.ProcessEvent(event)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("ProcessEvent_QueueFull", func(t *testing.T) {
		// Create a queue with size 1 and fill it
		smallQueue := NewMockEventQueue(1)
		smallService := NewProductService(repo, smallQueue, 1)

		// Fill the queue
		event1 := models.ProductEvent{ProductID: "test-1", Price: 10.0, Stock: 5}
		err := smallService.ProcessEvent(event1)
		if err != nil {
			t.Errorf("Expected no error for first event, got %v", err)
		}

		// Try to add another event (should fail)
		event2 := models.ProductEvent{ProductID: "test-2", Price: 20.0, Stock: 10}
		err = smallService.ProcessEvent(event2)
		if err == nil {
			t.Error("Expected error for second event when queue is full")
		}
	})
}

func TestProductService_GetProduct(t *testing.T) {
	repo := NewMockProductRepository()
	eventQueue := NewMockEventQueue(10)
	service := NewProductService(repo, eventQueue, 1)

	t.Run("GetProduct_Exists", func(t *testing.T) {
		// Add a product directly to repository
		repo.Update("test-product", 99.99, 50)

		product, exists := service.GetProduct("test-product")
		if !exists {
			t.Error("Expected product to exist")
		}
		if product.Price != 99.99 || product.Stock != 50 {
			t.Errorf("Expected price=99.99, stock=50, got price=%.2f, stock=%d", product.Price, product.Stock)
		}
	})

	t.Run("GetProduct_NotExists", func(t *testing.T) {
		_, exists := service.GetProduct("nonexistent")
		if exists {
			t.Error("Expected product to not exist")
		}
	})
}

func TestWorkerPool_ProcessEvent(t *testing.T) {
	repo := NewMockProductRepository()
	eventQueue := NewMockEventQueue(10)
	service := NewProductService(repo, eventQueue, 1)

	// Start the service
	service.Start()
	defer service.Stop()

	t.Run("WorkerProcessesEvent", func(t *testing.T) {
		event := models.ProductEvent{ProductID: "worker-test", Price: 15.0, Stock: 8}

		// Process the event
		err := service.ProcessEvent(event)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Wait for processing
		time.Sleep(50 * time.Millisecond)

		// Check if product was created
		product, exists := service.GetProduct("worker-test")
		if !exists {
			t.Error("Expected product to exist after processing")
		}
		if product.Price != 15.0 || product.Stock != 8 {
			t.Errorf("Expected price=15.0, stock=8, got price=%.2f, stock=%d", product.Price, product.Stock)
		}
	})
}

func TestWorkerPool_StartStop(t *testing.T) {
	repo := NewMockProductRepository()
	eventQueue := NewMockEventQueue(10)
	service := NewProductService(repo, eventQueue, 2)

	t.Run("StartStop", func(t *testing.T) {
		// Start service
		service.Start()

		// Verify it's running by processing an event
		event := models.ProductEvent{ProductID: "start-stop-test", Price: 25.0, Stock: 12}
		err := service.ProcessEvent(event)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Wait for processing
		time.Sleep(50 * time.Millisecond)

		// Stop service
		service.Stop()

		// Verify product was processed
		product, exists := service.GetProduct("start-stop-test")
		if !exists {
			t.Error("Expected product to exist after processing")
		}
		if product.Price != 25.0 || product.Stock != 12 {
			t.Errorf("Expected price=25.0, stock=12, got price=%.2f, stock=%d", product.Price, product.Stock)
		}
	})
}
