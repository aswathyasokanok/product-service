package repositories

import (
	"sync"
	"testing"
)

func TestInMemoryProductRepository(t *testing.T) {
	repo := NewInMemoryProductRepository()

	// Test initial state
	_, exists := repo.Get("test-product")
	if exists {
		t.Error("Expected product to not exist initially")
	}

	// Test update
	repo.Update("test-product", 99.99, 50)
	product, exists := repo.Get("test-product")
	if !exists {
		t.Error("Expected product to exist after update")
	}
	if product.Price != 99.99 || product.Stock != 50 {
		t.Errorf("Expected price=99.99, stock=50, got price=%.2f, stock=%d", product.Price, product.Stock)
	}

	// Test concurrent access
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			repo.Update("product-"+string(rune(id)), float64(id), id*10)
		}(i)
	}
	wg.Wait()

	// Verify the original product still exists
	originalProduct, exists := repo.Get("test-product")
	if !exists {
		t.Error("Expected original product to still exist")
	}
	if originalProduct.Price != 99.99 || originalProduct.Stock != 50 {
		t.Errorf("Expected original product to have price=99.99, stock=50, got price=%.2f, stock=%d", originalProduct.Price, originalProduct.Stock)
	}
}
