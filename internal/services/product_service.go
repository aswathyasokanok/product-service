package services

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"product-service/internal/models"
	"product-service/pkg/circuitbreaker"
	"product-service/pkg/queue"
	"product-service/pkg/retry"
)

// ProductService handles business logic for products
type ProductService struct {
	repository     ProductRepository
	queue          queue.EventQueue
	workerPool     *WorkerPool
	circuitBreaker *circuitbreaker.CircuitBreaker
	retryConfig    *retry.RetryConfig
}

// ProductRepository interface for dependency injection
type ProductRepository interface {
	Get(id string) (*models.Product, bool)
	Update(id string, price float64, stock int)
}

// NewProductService creates a new product service
func NewProductService(repo ProductRepository, eventQueue queue.EventQueue, workers int) *ProductService {
	service := &ProductService{
		repository:     repo,
		queue:          eventQueue,
		circuitBreaker: circuitbreaker.NewCircuitBreaker(5, 60*time.Second),
		retryConfig:    retry.DefaultRetryConfig(),
	}

	service.workerPool = NewWorkerPool(workers, eventQueue, repo, service.circuitBreaker, service.retryConfig)
	return service
}

// Start starts the product service and worker pool
func (s *ProductService) Start() {
	s.workerPool.Start()
}

// Stop gracefully stops the product service
func (s *ProductService) Stop() {
	s.workerPool.Stop()
}

// ProcessEvent enqueues a product event for processing with retry
func (s *ProductService) ProcessEvent(event models.ProductEvent) error {
	return s.retryConfig.ExecuteWithRetry(func() error {
		return s.circuitBreaker.Execute(func() error {
			return s.queue.Enqueue(event)
		})
	})
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(id string) (*models.Product, bool) {
	return s.repository.Get(id)
}

// WorkerPool manages a pool of workers for processing events
type WorkerPool struct {
	workers        int
	queue          queue.EventQueue
	repository     ProductRepository
	circuitBreaker *circuitbreaker.CircuitBreaker
	retryConfig    *retry.RetryConfig
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	logger         *log.Logger
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, eventQueue queue.EventQueue, repo ProductRepository, cb *circuitbreaker.CircuitBreaker, rc *retry.RetryConfig) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:        workers,
		queue:          eventQueue,
		repository:     repo,
		circuitBreaker: cb,
		retryConfig:    rc,
		ctx:            ctx,
		cancel:         cancel,
		logger:         log.New(os.Stdout, "[WORKER] ", log.LstdFlags),
	}
}

// Start starts all workers
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	wp.logger.Printf("Started %d workers", wp.workers)
}

// Stop gracefully stops all workers
func (wp *WorkerPool) Stop() {
	wp.logger.Println("Stopping workers...")
	wp.cancel()
	wp.wg.Wait()
	wp.logger.Println("All workers stopped")
}

// worker processes events from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	wp.logger.Printf("Worker %d started", id)

	for {
		select {
		case <-wp.ctx.Done():
			wp.logger.Printf("Worker %d stopping", id)
			return
		default:
			event, ok := wp.queue.Dequeue()
			if !ok {
				// Channel closed, exit
				wp.logger.Printf("Worker %d: queue closed, exiting", id)
				return
			}

			wp.processEvent(event, id)
		}
	}
}

// processEvent processes a single product event with retry and error handling
func (wp *WorkerPool) processEvent(event models.ProductEvent, workerID int) {
	wp.logger.Printf("Worker %d processing event for product %s", workerID, event.ProductID)

	// Process with retry and circuit breaker
	err := wp.retryConfig.ExecuteWithRetryAndCallback(
		func() error {
			return wp.circuitBreaker.Execute(func() error {
				// Simulate some processing time
				time.Sleep(10 * time.Millisecond)

				// Update the product repository
				wp.repository.Update(event.ProductID, event.Price, event.Stock)

				wp.logger.Printf("Worker %d updated product %s: price=%.2f, stock=%d",
					workerID, event.ProductID, event.Price, event.Stock)

				return nil
			})
		},
		func(attempt int, err error) {
			wp.logger.Printf("Worker %d attempt %d failed for product %s: %v",
				workerID, attempt, event.ProductID, err)
		},
	)

	if err != nil {
		// Log the final failure
		wp.logger.Printf("Worker %d failed to process event for product %s after all retries: %v",
			workerID, event.ProductID, err)

		// In a production system, you would send this to a dead letter queue
		// or persistent storage for later analysis
	}
}
