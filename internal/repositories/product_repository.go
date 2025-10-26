package repositories

import (
	"sync"

	"product-service/internal/models"
)

// ProductRepository interface defines the contract for product storage
type ProductRepository interface {
	Get(id string) (*models.Product, bool)
	Update(id string, price float64, stock int)
}

// InMemoryProductRepository implements ProductRepository using in-memory storage
type InMemoryProductRepository struct {
	mu   sync.RWMutex
	data map[string]*models.Product
}

// NewInMemoryProductRepository creates a new in-memory product repository
func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		data: make(map[string]*models.Product),
	}
}

// Get retrieves a product by ID
func (r *InMemoryProductRepository) Get(id string) (*models.Product, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	product, exists := r.data[id]
	return product, exists
}

// Update updates a product's state
func (r *InMemoryProductRepository) Update(id string, price float64, stock int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[id] = &models.Product{
		ID:    id,
		Price: price,
		Stock: stock,
	}
}
