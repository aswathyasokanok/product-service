package queue

import (
	"sync"
	"time"

	"product-service/internal/models"
)

// BatchProcessor handles batch processing of events for high throughput
type BatchProcessor struct {
	batchSize     int
	flushInterval time.Duration
	events        []models.ProductEvent
	mutex         sync.Mutex
	flushChan     chan []models.ProductEvent
	stopChan      chan struct{}
	processor     BatchProcessorFunc
}

// BatchProcessorFunc defines the function signature for processing batches
type BatchProcessorFunc func(events []models.ProductEvent) error

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(batchSize int, flushInterval time.Duration, processor BatchProcessorFunc) *BatchProcessor {
	bp := &BatchProcessor{
		batchSize:     batchSize,
		flushInterval: flushInterval,
		events:        make([]models.ProductEvent, 0, batchSize),
		flushChan:     make(chan []models.ProductEvent, 10),
		stopChan:      make(chan struct{}),
		processor:     processor,
	}

	// Start the batch processing goroutine
	go bp.processBatches()

	return bp
}

// AddEvent adds an event to the batch
func (bp *BatchProcessor) AddEvent(event models.ProductEvent) error {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()

	bp.events = append(bp.events, event)

	// Flush if batch is full
	if len(bp.events) >= bp.batchSize {
		return bp.flushBatch()
	}

	return nil
}

// flushBatch flushes the current batch
func (bp *BatchProcessor) flushBatch() error {
	if len(bp.events) == 0 {
		return nil
	}

	// Create a copy of the events to send
	eventsToProcess := make([]models.ProductEvent, len(bp.events))
	copy(eventsToProcess, bp.events)

	// Clear the current batch
	bp.events = bp.events[:0]

	// Send to processing channel
	select {
	case bp.flushChan <- eventsToProcess:
		return nil
	default:
		return ErrBatchProcessorFull
	}
}

// processBatches processes batches from the flush channel
func (bp *BatchProcessor) processBatches() {
	ticker := time.NewTicker(bp.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case events := <-bp.flushChan:
			if err := bp.processor(events); err != nil {
				// Log error or send to dead letter queue
				// In production, you would have proper error handling here
			}
		case <-ticker.C:
			// Periodic flush
			bp.mutex.Lock()
			if len(bp.events) > 0 {
				bp.flushBatch()
			}
			bp.mutex.Unlock()
		case <-bp.stopChan:
			// Flush remaining events before stopping
			bp.mutex.Lock()
			if len(bp.events) > 0 {
				bp.flushBatch()
			}
			bp.mutex.Unlock()
			return
		}
	}
}

// Stop stops the batch processor
func (bp *BatchProcessor) Stop() {
	close(bp.stopChan)
}

// GetBatchSize returns the current batch size
func (bp *BatchProcessor) GetBatchSize() int {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()
	return len(bp.events)
}

// GetPendingEvents returns the number of pending events
func (bp *BatchProcessor) GetPendingEvents() int {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()
	return len(bp.events)
}
