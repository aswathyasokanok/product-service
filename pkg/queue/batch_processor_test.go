package queue

import (
	"errors"
	"sync"
	"testing"
	"time"

	"product-service/internal/models"
)

func TestBatchProcessor_NewBatchProcessor(t *testing.T) {
	processedBatches := make([][]models.ProductEvent, 0)
	processor := NewBatchProcessor(5, 100*time.Millisecond, func(events []models.ProductEvent) error {
		processedBatches = append(processedBatches, events)
		return nil
	})

	if processor.batchSize != 5 {
		t.Errorf("Expected batch size 5, got %d", processor.batchSize)
	}
	if processor.flushInterval != 100*time.Millisecond {
		t.Errorf("Expected flush interval 100ms, got %v", processor.flushInterval)
	}
	if len(processor.events) != 0 {
		t.Errorf("Expected empty events slice, got %d events", len(processor.events))
	}
}

func TestBatchProcessor_AddEvent_SingleEvent(t *testing.T) {
	processedBatches := make([][]models.ProductEvent, 0)
	var mu sync.Mutex

	processor := NewBatchProcessor(5, 100*time.Millisecond, func(events []models.ProductEvent) error {
		mu.Lock()
		processedBatches = append(processedBatches, events)
		mu.Unlock()
		return nil
	})

	event := models.ProductEvent{ProductID: "test-1", Price: 10.0, Stock: 5}
	err := processor.AddEvent(event)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Wait for processing
	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	if len(processedBatches) != 1 {
		t.Errorf("Expected 1 processed batch, got %d", len(processedBatches))
	}
	if len(processedBatches[0]) != 1 {
		t.Errorf("Expected 1 event in batch, got %d", len(processedBatches[0]))
	}
	mu.Unlock()
}

func TestBatchProcessor_AddEvent_BatchSizeReached(t *testing.T) {
	processedBatches := make([][]models.ProductEvent, 0)
	var mu sync.Mutex

	processor := NewBatchProcessor(3, 1*time.Second, func(events []models.ProductEvent) error {
		mu.Lock()
		processedBatches = append(processedBatches, events)
		mu.Unlock()
		return nil
	})

	// Add 3 events to trigger batch processing
	for i := 0; i < 3; i++ {
		event := models.ProductEvent{ProductID: string(rune(i)), Price: float64(i), Stock: i}
		err := processor.AddEvent(event)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}

	// Wait a bit for processing
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	if len(processedBatches) != 1 {
		t.Errorf("Expected 1 processed batch, got %d", len(processedBatches))
	}
	if len(processedBatches[0]) != 3 {
		t.Errorf("Expected 3 events in batch, got %d", len(processedBatches[0]))
	}
	mu.Unlock()
}

func TestBatchProcessor_AddEvent_FlushInterval(t *testing.T) {
	processedBatches := make([][]models.ProductEvent, 0)
	var mu sync.Mutex

	processor := NewBatchProcessor(10, 50*time.Millisecond, func(events []models.ProductEvent) error {
		mu.Lock()
		processedBatches = append(processedBatches, events)
		mu.Unlock()
		return nil
	})

	// Add 1 event
	event := models.ProductEvent{ProductID: "test-1", Price: 10.0, Stock: 5}
	err := processor.AddEvent(event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Wait for flush interval
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if len(processedBatches) != 1 {
		t.Errorf("Expected 1 processed batch, got %d", len(processedBatches))
	}
	if len(processedBatches[0]) != 1 {
		t.Errorf("Expected 1 event in batch, got %d", len(processedBatches[0]))
	}
	mu.Unlock()
}

func TestBatchProcessor_AddEvent_ConcurrentAccess(t *testing.T) {
	processedBatches := make([][]models.ProductEvent, 0)
	var mu sync.Mutex

	processor := NewBatchProcessor(10, 100*time.Millisecond, func(events []models.ProductEvent) error {
		mu.Lock()
		processedBatches = append(processedBatches, events)
		mu.Unlock()
		return nil
	})

	// Add events concurrently
	var wg sync.WaitGroup
	numEvents := 20

	for i := 0; i < numEvents; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			event := models.ProductEvent{ProductID: string(rune(id)), Price: float64(id), Stock: id}
			err := processor.AddEvent(event)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		}(i)
	}

	wg.Wait()

	// Wait for processing
	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	if len(processedBatches) == 0 {
		t.Error("Expected at least one processed batch")
	}

	// Count total processed events
	totalProcessed := 0
	for _, batch := range processedBatches {
		totalProcessed += len(batch)
	}
	if totalProcessed != numEvents {
		t.Errorf("Expected %d total processed events, got %d", numEvents, totalProcessed)
	}
	mu.Unlock()
}

func TestBatchProcessor_AddEvent_ProcessorError(t *testing.T) {
	processor := NewBatchProcessor(1, 100*time.Millisecond, func(events []models.ProductEvent) error {
		return errors.New("processing error")
	})

	event := models.ProductEvent{ProductID: "test-1", Price: 10.0, Stock: 5}
	err := processor.AddEvent(event)

	if err != nil {
		t.Errorf("Expected no error on AddEvent, got %v", err)
	}

	// Wait for processing
	time.Sleep(150 * time.Millisecond)

	// The error should be handled internally and not returned to AddEvent
	// This test mainly ensures the processor doesn't crash on errors
}

func TestBatchProcessor_Stop(t *testing.T) {
	processedBatches := make([][]models.ProductEvent, 0)
	var mu sync.Mutex

	processor := NewBatchProcessor(10, 100*time.Millisecond, func(events []models.ProductEvent) error {
		mu.Lock()
		processedBatches = append(processedBatches, events)
		mu.Unlock()
		return nil
	})

	// Add some events
	event1 := models.ProductEvent{ProductID: "test-1", Price: 10.0, Stock: 5}
	event2 := models.ProductEvent{ProductID: "test-2", Price: 20.0, Stock: 10}

	processor.AddEvent(event1)
	processor.AddEvent(event2)

	// Stop the processor
	processor.Stop()

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Should have processed the remaining events
	mu.Lock()
	if len(processedBatches) == 0 {
		t.Error("Expected at least 1 processed batch")
	} else {
		if len(processedBatches[0]) != 2 {
			t.Errorf("Expected 2 events in batch, got %d", len(processedBatches[0]))
		}
	}
	mu.Unlock()
}

func TestBatchProcessor_FlushBatch(t *testing.T) {
	processedBatches := make([][]models.ProductEvent, 0)
	var mu sync.Mutex

	processor := NewBatchProcessor(5, 100*time.Millisecond, func(events []models.ProductEvent) error {
		mu.Lock()
		processedBatches = append(processedBatches, events)
		mu.Unlock()
		return nil
	})

	// Add 2 events
	event1 := models.ProductEvent{ProductID: "test-1", Price: 10.0, Stock: 5}
	event2 := models.ProductEvent{ProductID: "test-2", Price: 20.0, Stock: 10}

	processor.AddEvent(event1)
	processor.AddEvent(event2)

	// Manually flush
	err := processor.flushBatch()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	if len(processedBatches) != 1 {
		t.Errorf("Expected 1 processed batch, got %d", len(processedBatches))
	}
	if len(processedBatches[0]) != 2 {
		t.Errorf("Expected 2 events in batch, got %d", len(processedBatches[0]))
	}
	mu.Unlock()
}

func TestBatchProcessor_EmptyFlush(t *testing.T) {
	processedBatches := make([][]models.ProductEvent, 0)
	var mu sync.Mutex

	processor := NewBatchProcessor(5, 100*time.Millisecond, func(events []models.ProductEvent) error {
		mu.Lock()
		processedBatches = append(processedBatches, events)
		mu.Unlock()
		return nil
	})

	// Flush empty batch
	err := processor.flushBatch()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	if len(processedBatches) != 0 {
		t.Errorf("Expected 0 processed batches, got %d", len(processedBatches))
	}
	mu.Unlock()
}
