package queue

import (
	"sync"
	"testing"

	"product-service/internal/models"
)

func TestInMemoryEventQueue(t *testing.T) {
	q := NewInMemoryEventQueue(10)

	// Test enqueue
	event := models.ProductEvent{ProductID: "test", Price: 10.0, Stock: 5}
	err := q.Enqueue(event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test dequeue
	dequeuedEvent, ok := q.Dequeue()
	if !ok {
		t.Error("Expected to dequeue event")
	}
	if dequeuedEvent.ProductID != event.ProductID {
		t.Errorf("Expected product ID %s, got %s", event.ProductID, dequeuedEvent.ProductID)
	}

	// Test queue full
	q = NewInMemoryEventQueue(2)
	if err := q.Enqueue(models.ProductEvent{ProductID: "1", Price: 1.0, Stock: 1}); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if err := q.Enqueue(models.ProductEvent{ProductID: "2", Price: 2.0, Stock: 2}); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	err = q.Enqueue(models.ProductEvent{ProductID: "3", Price: 3.0, Stock: 3})
	if err == nil {
		t.Error("Expected error when queue is full")
	}
}

func TestInMemoryEventQueue_ConcurrentAccess(t *testing.T) {
	q := NewInMemoryEventQueue(100)

	// Test concurrent enqueue
	var wg sync.WaitGroup
	numEvents := 50

	for i := 0; i < numEvents; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			event := models.ProductEvent{ProductID: string(rune(id)), Price: float64(id), Stock: id}
			err := q.Enqueue(event)
			if err != nil {
				t.Errorf("Unexpected error enqueuing event %d: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	// Test concurrent dequeue
	dequeuedCount := 0
	for i := 0; i < numEvents; i++ {
		_, ok := q.Dequeue()
		if ok {
			dequeuedCount++
		}
	}

	if dequeuedCount != numEvents {
		t.Errorf("Expected to dequeue %d events, got %d", numEvents, dequeuedCount)
	}
}

func TestInMemoryEventQueue_Close(t *testing.T) {
	q := NewInMemoryEventQueue(10)

	// Add some events
	q.Enqueue(models.ProductEvent{ProductID: "1", Price: 1.0, Stock: 1})
	q.Enqueue(models.ProductEvent{ProductID: "2", Price: 2.0, Stock: 2})

	// Close the queue
	q.Close()

	// Try to enqueue after close (should fail)
	err := q.Enqueue(models.ProductEvent{ProductID: "3", Price: 3.0, Stock: 3})
	if err == nil {
		t.Error("Expected error when enqueuing to closed queue")
	}

	// Should still be able to dequeue existing events
	event1, ok1 := q.Dequeue()
	if !ok1 {
		t.Error("Expected to dequeue first event")
	}
	if event1.ProductID != "1" {
		t.Errorf("Expected product ID '1', got '%s'", event1.ProductID)
	}

	event2, ok2 := q.Dequeue()
	if !ok2 {
		t.Error("Expected to dequeue second event")
	}
	if event2.ProductID != "2" {
		t.Errorf("Expected product ID '2', got '%s'", event2.ProductID)
	}

	// Queue should be empty now
	_, ok3 := q.Dequeue()
	if ok3 {
		t.Error("Expected queue to be empty")
	}
}

func TestInMemoryEventQueue_EmptyDequeue(t *testing.T) {
	q := NewInMemoryEventQueue(10)

	// Try to dequeue from empty queue
	_, ok := q.Dequeue()
	if ok {
		t.Error("Expected no event from empty queue")
	}
}

func TestInMemoryEventQueue_ZeroSize(t *testing.T) {
	q := NewInMemoryEventQueue(0)

	// Should fail to enqueue to zero-size queue
	err := q.Enqueue(models.ProductEvent{ProductID: "1", Price: 1.0, Stock: 1})
	if err == nil {
		t.Error("Expected error when enqueuing to zero-size queue")
	}
}

func TestInMemoryEventQueue_BlockingBehavior(t *testing.T) {
	q := NewInMemoryEventQueue(1)

	// Fill the queue
	err := q.Enqueue(models.ProductEvent{ProductID: "1", Price: 1.0, Stock: 1})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Try to enqueue another event (should fail immediately)
	err = q.Enqueue(models.ProductEvent{ProductID: "2", Price: 2.0, Stock: 2})
	if err == nil {
		t.Error("Expected error when queue is full")
	}
	if err.Error() != "queue is full" {
		t.Errorf("Expected 'queue is full', got '%s'", err.Error())
	}
}

func TestInMemoryEventQueue_EventIntegrity(t *testing.T) {
	q := NewInMemoryEventQueue(10)

	originalEvent := models.ProductEvent{
		ProductID: "test-product",
		Price:     99.99,
		Stock:     50,
	}

	// Enqueue event
	err := q.Enqueue(originalEvent)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Dequeue event
	dequeuedEvent, ok := q.Dequeue()
	if !ok {
		t.Error("Expected to dequeue event")
	}

	// Verify event integrity
	if dequeuedEvent.ProductID != originalEvent.ProductID {
		t.Errorf("Expected ProductID %s, got %s", originalEvent.ProductID, dequeuedEvent.ProductID)
	}
	if dequeuedEvent.Price != originalEvent.Price {
		t.Errorf("Expected Price %.2f, got %.2f", originalEvent.Price, dequeuedEvent.Price)
	}
	if dequeuedEvent.Stock != originalEvent.Stock {
		t.Errorf("Expected Stock %d, got %d", originalEvent.Stock, dequeuedEvent.Stock)
	}
}
