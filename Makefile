# Product Service Makefile

# Variables
APP_NAME=product-service
DOCKER_IMAGE=product-service
DOCKER_TAG=latest
PORT=8080
WORKERS=3
QUEUE_SIZE=1000

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=main
BINARY_UNIX=$(BINARY_NAME)_unix

# Default target
.PHONY: all
all: clean deps test build

# Install dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go

# Build for Linux
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/main.go

# Run the application
.PHONY: run
run:
	WORKERS=$(WORKERS) QUEUE_SIZE=$(QUEUE_SIZE) PORT=$(PORT) $(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go && ./$(BINARY_NAME)

# Run with hot reload (requires air)
.PHONY: dev
dev:
	air

# Test the application
.PHONY: test
test:
	$(GOTEST) -v ./...

# Test with race detection
.PHONY: test-race
test-race:
	$(GOTEST) -race -v ./...

# Test with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Benchmark tests
.PHONY: benchmark
benchmark:
	$(GOTEST) -bench=. -benchmem ./...

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out
	rm -f coverage.html

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Docker build
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker run
.PHONY: docker-run
docker-run:
	docker run -p $(PORT):8080 -e WORKERS=$(WORKERS) -e QUEUE_SIZE=$(QUEUE_SIZE) $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker compose up
.PHONY: docker-up
docker-up:
	docker-compose up --build

# Docker compose down
.PHONY: docker-down
docker-down:
	docker-compose down

# Docker compose logs
.PHONY: docker-logs
docker-logs:
	docker-compose logs -f

# Install development tools
.PHONY: install-tools
install-tools:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Generate mocks (if using mockgen)
.PHONY: generate-mocks
generate-mocks:
	mockgen -source=internal/repositories/product_repository.go -destination=internal/mocks/product_repository_mock.go
	mockgen -source=pkg/queue/event_queue.go -destination=pkg/mocks/event_queue_mock.go

# Security scan
.PHONY: security-scan
security-scan:
	gosec ./...

# Performance test
.PHONY: perf-test
perf-test:
	ab -n 10000 -c 100 http://localhost:$(PORT)/health

# Load test with events
.PHONY: load-test
load-test:
	@echo "Starting load test..."
	@for i in {1..1000}; do \
		curl -s -X POST http://localhost:$(PORT)/api/v1/events \
		-H "Content-Type: application/json" \
		-d "{\"product_id\":\"product-$$i\",\"price\":$$(($$RANDOM % 1000)),\"stock\":$$(($$RANDOM % 100))}" > /dev/null; \
	done
	@echo "Load test completed"


