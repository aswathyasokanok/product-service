# Product Service

A high-performance Go microservice for handling product updates with asynchronous processing, built with MVC architecture, Gin framework, and designed for scalability.

## Architecture

This project follows a clean MVC (Model-View-Controller) architecture with the following structure:

```
product-service/
├── cmd/                    # Application entry point
│   └── main.go
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── controllers/      # HTTP request handlers
│   ├── models/           # Data models
│   ├── repositories/     # Data access layer
│   └── services/         # Business logic layer
├── pkg/                  # Public library code
│   ├── queue/           # Event queue implementation
│   └── store/           # Storage implementations
├── api/                  # API definitions
│   └── v1/              # API version 1

├── Dockerfile           # Container definition
├── docker-compose.yml   # Multi-container setup
├── Makefile            # Build automation
└── README.md           # This file
```

## Features

- **MVC Architecture**: Clean separation of concerns with models, views, and controllers
- **RESTful API**: Clean HTTP endpoints for product management
- **Asynchronous Processing**: Event-driven architecture with worker pools
- **Thread-Safe Storage**: Concurrent access to in-memory product store
- **Graceful Shutdown**: Proper cleanup of resources on termination
- **Comprehensive Testing**: Unit tests, concurrency tests, and benchmarks
- **Production Ready**: Configurable workers, structured logging, health checks
- **Docker Support**: Containerized deployment with multi-stage builds
- **Makefile**: Automated build, test, and deployment tasks

## API Endpoints

### POST /api/v1/events
Accepts JSON payloads representing product updates.

**Request Body:**
```json
{
  "product_id": "abc123",
  "price": 49.99,
  "stock": 100
}
```

**Response:**
- `202 Accepted`: Event successfully enqueued
- `400 Bad Request`: Invalid JSON or missing required fields
- `503 Service Unavailable`: Queue is full

### GET /api/v1/products/{id}
Retrieves the current state of a product.

**Response:**
- `200 OK`: Product data
- `404 Not Found`: Product doesn't exist

**Example Response:**
```json
{
  "id": "abc123",
  "price": 49.99,
  "stock": 100
}
```

### GET /health
Health check endpoint for monitoring.

**Response:**
```json
{
  "status": "healthy"
}
```

## How to Run the Application

### Prerequisites

- **Go 1.19+** installed on your system
- **Docker** and **Docker Compose** (for containerized deployment)
- **Make** (optional, for using Makefile commands)

### Method 1: Using Make (Recommended)

This is the easiest way to run the application with all dependencies managed automatically.

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd product-service
   ```

2. **Install dependencies:**
   ```bash
   make deps
   ```

3. **Run the application:**
   ```bash
   make run
   ```

4. **Verify the service is running:**
   ```bash
   curl http://localhost:8080/health
   ```

### Method 2: Using Docker Compose (Containerized)

Perfect for production-like environments and consistent deployment.

1. **Start the service with Docker Compose:**
   ```bash
   docker-compose up --build
   ```

2. **Run in detached mode:**
   ```bash
   docker-compose up -d --build
   ```

3. **View logs:**
   ```bash
   docker-compose logs -f
   ```

4. **Stop the service:**
   ```bash
   docker-compose down
   ```

## Testing the Application

Once the service is running, you can test it using the following commands:

### 1. Health Check
```bash
curl http://localhost:8080/health
```
Expected response:
```json
{
  "status": "healthy"
}
```

### 2. Send a Product Event
```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{"product_id": "abc123", "price": 49.99, "stock": 100}'
```
Expected response: `202 Accepted`

### 3. Retrieve Product State
```bash
curl http://localhost:8080/api/v1/products/abc123
```
Expected response:
```json
{
  "id": "abc123",
  "price": 49.99,
  "stock": 100
}
```

### 4. Load Testing
```bash
# Run load test (sends 1000 events)
make load-test

# Run performance test
make perf-test
```

### Method 3: Manual Go Commands

For development and debugging purposes.

1. **Install dependencies:**
   ```bash
   go mod tidy
   go mod download
   ```

2. **Run the service:**
   ```bash
   go run cmd/main.go
   ```

3. **Run with custom configuration:**
   ```bash
   WORKERS=5 QUEUE_SIZE=2000 PORT=8080 go run cmd/main.go
   ```

4. **Build and run binary:**
   ```bash
   go build -o product-service cmd/main.go
   ./product-service
   ```

### Method 4: Using Docker (Single Container)

1. **Build the Docker image:**
   ```bash
   make docker-build
   # or
   docker build -t product-service .
   ```

2. **Run the container:**
   ```bash
   make docker-run
   # or
   docker run -p 8080:8080 -e WORKERS=5 -e QUEUE_SIZE=2000 product-service
   ```

## Quick Start

### Using Make

1. **Clone and setup:**
   ```bash
   git clone <repository-url>
   cd product-service
   make deps
   ```

2. **Run the service:**
   ```bash
   make run
   ```

3. **Run tests:**
   ```bash
   make test
   ```

4. **Build Docker image:**
   ```bash
   make docker-build
   ```

### Using Docker Compose

1. **Start the service:**
   ```bash
   docker-compose up --build
   ```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `WORKERS` | 3 | Number of worker goroutines |
| `QUEUE_SIZE` | 1000 | Size of the event queue buffer |
| `PORT` | 8080 | HTTP server port |

### Example Usage

1. **Send a product update:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/events \
     -H "Content-Type: application/json" \
     -d '{"product_id": "abc123", "price": 49.99, "stock": 100}'
   ```

2. **Retrieve product state:**
   ```bash
   curl http://localhost:8080/api/v1/products/abc123
   ```

3. **Health check:**
   ```bash
   curl http://localhost:8080/health
   ```

## Development

### Available Make Targets

```bash
make help              # Show all available targets
make deps              # Install dependencies
make build             # Build the application
make run               # Run the application
make dev               # Run with hot reload (requires air)
make test              # Run tests
make test-race         # Run tests with race detection
make test-coverage     # Run tests with coverage report
make benchmark         # Run benchmark tests
make clean             # Clean build artifacts
make fmt               # Format code
make lint              # Lint code
make docker-build      # Build Docker image
make docker-run        # Run Docker container
make docker-up         # Start with docker-compose
make docker-down       # Stop docker-compose
make load-test         # Run load test
make perf-test         # Run performance test
```

### Development Tools

Install development tools:
```bash
make install-tools
```

This installs:
- `golangci-lint` for linting

### Hot Reload

For development with hot reload:
```bash
make dev
```


## Testing

### Run All Tests
```bash
make test
```

### Run Tests with Race Detection
```bash
make test-race
```

### Run Tests with Coverage
```bash
make test-coverage
```

### Run Benchmarks
```bash
make benchmark
```

### Load Testing
```bash
# Start the service
make run

# In another terminal, run load test
make load-test
```

## Docker

### Build Image
```bash
make docker-build
```

### Run Container
```bash
make docker-run
```



## Architecture Details

### MVC Components

1. **Models** (`internal/models/`):
   - `Product`: Represents a product with price and stock
   - `ProductEvent`: Represents incoming product updates
   - Response models for API endpoints

2. **Controllers** (`internal/controllers/`):
   - `ProductController`: Handles product-related HTTP requests
   - `HealthController`: Handles health check requests

3. **Services** (`internal/services/`):
   - `ProductService`: Contains business logic and orchestrates operations
   - `WorkerPool`: Manages asynchronous event processing

4. **Repositories** (`internal/repositories/`):
   - `ProductRepository`: Interface for data access
   - `InMemoryProductRepository`: In-memory implementation

### Key Design Patterns

1. **Dependency Injection**: Services are injected into controllers
2. **Interface Segregation**: Small, focused interfaces for each component
3. **Repository Pattern**: Abstracts data access layer
4. **Worker Pool Pattern**: Configurable workers for async processing
5. **Context Pattern**: Graceful shutdown using context cancellation

## Production Considerations

### Large-Scale Data & High Throughput Strategies

#### 1. **Batch Processing for High Throughput**
```go
// pkg/queue/batch_processor.go
type BatchProcessor struct {
    batchSize    int
    flushInterval time.Duration
    events       []models.ProductEvent
    mutex        sync.Mutex
}

func (bp *BatchProcessor) ProcessBatch(events []models.ProductEvent) error {
    // Process multiple events in a single transaction
    return bp.repository.BatchUpdate(events)
}
```

#### 2. **Connection Pooling & Resource Management**
```go
// pkg/store/connection_pool.go
type ConnectionPool struct {
    maxConnections int
    idleTimeout    time.Duration
    pool           chan *sql.DB
}

func (cp *ConnectionPool) GetConnection() (*sql.DB, error) {
    select {
    case conn := <-cp.pool:
        return conn, nil
    case <-time.After(5 * time.Second):
        return nil, errors.New("connection timeout")
    }
}
```

#### 3. **Memory Management for Large Datasets**
```go
// pkg/store/memory_manager.go
type MemoryManager struct {
    maxMemoryUsage int64
    cleanupThreshold float64
    gcInterval time.Duration
}

func (mm *MemoryManager) MonitorMemory() {
    // Monitor memory usage and trigger cleanup
    if mm.getMemoryUsage() > mm.maxMemoryUsage * mm.cleanupThreshold {
        mm.performCleanup()
    }
}
```

#### 4. **Horizontal Scaling Strategies**
- **Load Balancing**: Multiple service instances behind a load balancer
- **Data Partitioning**: Shard products by ID ranges or hash
- **Read Replicas**: Separate read/write operations
- **Caching Layers**: Multi-level caching (L1: in-memory, L2: Redis)

### Error Handling & Retry Mechanisms

#### 1. **Exponential Backoff Retry**
```go
// pkg/retry/retry.go
type RetryConfig struct {
    MaxAttempts int
    InitialDelay time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
}

func (r *RetryConfig) ExecuteWithRetry(operation func() error) error {
    delay := r.InitialDelay
    for attempt := 1; attempt <= r.MaxAttempts; attempt++ {
        if err := operation(); err == nil {
            return nil
        }
        
        if attempt == r.MaxAttempts {
            return fmt.Errorf("operation failed after %d attempts", r.MaxAttempts)
        }
        
        time.Sleep(delay)
        delay = time.Duration(float64(delay) * r.Multiplier)
        if delay > r.MaxDelay {
            delay = r.MaxDelay
        }
    }
    return nil
}
```

#### 2. **Circuit Breaker Pattern**
```go
// pkg/circuitbreaker/circuit_breaker.go
type CircuitBreaker struct {
    failureThreshold int
    timeout          time.Duration
    state           State
    failures        int
    lastFailureTime time.Time
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
    if cb.state == Open && time.Since(cb.lastFailureTime) < cb.timeout {
        return errors.New("circuit breaker is open")
    }
    
    err := operation()
    if err != nil {
        cb.recordFailure()
        return err
    }
    
    cb.recordSuccess()
    return nil
}
```

#### 3. **Dead Letter Queue for Failed Events**
```go
// pkg/queue/dead_letter_queue.go
type DeadLetterQueue struct {
    failedEvents chan models.FailedEvent
    maxRetries   int
}

type FailedEvent struct {
    Event     models.ProductEvent
    Error     error
    Timestamp time.Time
    RetryCount int
}

func (dlq *DeadLetterQueue) HandleFailedEvent(event models.ProductEvent, err error) {
    failedEvent := models.FailedEvent{
        Event:     event,
        Error:     err,
        Timestamp: time.Now(),
        RetryCount: 0,
    }
    
    select {
    case dlq.failedEvents <- failedEvent:
    default:
        // Log to persistent storage if queue is full
        dlq.logToPersistentStorage(failedEvent)
    }
}
```

#### 4. **Comprehensive Error Classification**
```go
// pkg/errors/error_types.go
type ErrorType int

const (
    RetryableError ErrorType = iota
    NonRetryableError
    ValidationError
    SystemError
)

type ClassifiedError struct {
    Type    ErrorType
    Message string
    Cause   error
}

func (ce *ClassifiedError) ShouldRetry() bool {
    return ce.Type == RetryableError
}
```

### Production Configuration

#### Environment Variables
```bash
# High Throughput Configuration
BATCH_SIZE=1000
BATCH_FLUSH_INTERVAL=1s
MAX_CONNECTIONS=100
CONNECTION_TIMEOUT=30s

# Error Handling Configuration
MAX_RETRY_ATTEMPTS=3
INITIAL_RETRY_DELAY=100ms
MAX_RETRY_DELAY=30s
CIRCUIT_BREAKER_THRESHOLD=5
CIRCUIT_BREAKER_TIMEOUT=60s

# Memory Management
MAX_MEMORY_USAGE=1GB
CLEANUP_THRESHOLD=0.8
GC_INTERVAL=30s
```

#### Docker Compose for Production
```yaml
version: '3.8'
services:
  product-service:
    build: .
    environment:
      - BATCH_SIZE=1000
      - MAX_RETRY_ATTEMPTS=3
      - CIRCUIT_BREAKER_THRESHOLD=5
    deploy:
      replicas: 3
      resources:
        limits:
          memory: 1G
          cpus: '0.5'
        reservations:
          memory: 512M
          cpus: '0.25'
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

### Monitoring and Observability

1. **Metrics**: Prometheus metrics for monitoring
2. **Tracing**: OpenTelemetry for distributed tracing
3. **Logging**: Structured logging with zap
4. **Health Checks**: Comprehensive health endpoints

### Security

1. **Input Validation**: Comprehensive request validation
2. **Rate Limiting**: Prevent abuse with rate limiting
3. **Authentication**: JWT-based authentication
4. **Authorization**: Role-based access control

## Troubleshooting

### Common Issues

1. **Port Already in Use**:
   ```bash
   # Change port
   PORT=8081 make run
   # or
   docker run -p 8081:8080 product-service
   ```

2. **Docker Build Fails**:
   ```bash
   # Clean Docker cache
   docker system prune -a
   make docker-build
   ```

3. **Tests Fail**:
   ```bash
   # Run with verbose output
   make test-race
   ```

4. **Service Won't Start**:
   ```bash
   # Check if Go is installed
   go version
   
   # Check if dependencies are installed
   make deps
   
   # Run with debug logging
   LOG_LEVEL=debug make run
   ```

5. **Docker Container Exits Immediately**:
   ```bash
   # Check container logs
   docker logs <container-id>
   
   # Run container interactively
   docker run -it product-service /bin/sh
   ```

6. **Permission Denied on Linux**:
   ```bash
   # Make binary executable
   chmod +x product-service
   
   
   ```

7. **Memory Issues**:
   ```bash
   # Reduce queue size and workers
   WORKERS=1 QUEUE_SIZE=100 make run
   ```

8. **Network Connection Issues**:
   ```bash
   # Check if port is accessible
   netstat -tulpn | grep 8080
   
   # Test with different port
   PORT=8081 make run
   ```

### Debugging

1. **Enable Debug Logging**:
   ```bash
   LOG_LEVEL=debug make run
   ```

2. **Profile the Application**:
   ```bash
   go tool pprof http://localhost:8080/debug/pprof/profile
   ```

3. **Check Container Logs**:
   ```bash
   make docker-logs
   ```

