package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// State represents the state of the circuit breaker
type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	failureThreshold int
	timeout          time.Duration
	state            State
	failures         int
	lastFailureTime  time.Time
	mutex            sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(failureThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		timeout:          timeout,
		state:            Closed,
		failures:         0,
	}
}

// Execute executes an operation with circuit breaker protection
func (cb *CircuitBreaker) Execute(operation func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Check if circuit breaker is open
	if cb.state == Open {
		if time.Since(cb.lastFailureTime) < cb.timeout {
			return errors.New("circuit breaker is open")
		}
		// Timeout has passed, move to half-open state
		cb.state = HalfOpen
	}

	// Execute the operation
	err := operation()

	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

// recordFailure records a failure and updates the circuit breaker state
func (cb *CircuitBreaker) recordFailure() {
	cb.failures++
	cb.lastFailureTime = time.Now()

	if cb.failures >= cb.failureThreshold {
		cb.state = Open
	}
}

// recordSuccess records a success and resets the circuit breaker if needed
func (cb *CircuitBreaker) recordSuccess() {
	if cb.state == HalfOpen {
		cb.state = Closed
	}
	cb.failures = 0
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// GetFailureCount returns the current failure count
func (cb *CircuitBreaker) GetFailureCount() int {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.failures
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.state = Closed
	cb.failures = 0
}
