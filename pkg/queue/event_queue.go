package queue

import (
	"fmt"

	"product-service/internal/models"
)

// EventQueue interface defines the contract for event queuing
type EventQueue interface {
	Enqueue(event models.ProductEvent) error
	Dequeue() (models.ProductEvent, bool)
	Close()
}

// InMemoryEventQueue implements EventQueue using buffered channels
type InMemoryEventQueue struct {
	events chan models.ProductEvent
}

// NewInMemoryEventQueue creates a new in-memory event queue with specified buffer size
func NewInMemoryEventQueue(bufferSize int) EventQueue {
	return &InMemoryEventQueue{
		events: make(chan models.ProductEvent, bufferSize),
	}
}

// Enqueue adds an event to the queue
func (q *InMemoryEventQueue) Enqueue(event models.ProductEvent) error {
	select {
	case q.events <- event:
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

// Dequeue retrieves an event from the queue
func (q *InMemoryEventQueue) Dequeue() (models.ProductEvent, bool) {
	event, ok := <-q.events
	return event, ok
}

// Close closes the event queue
func (q *InMemoryEventQueue) Close() {
	close(q.events)
}
