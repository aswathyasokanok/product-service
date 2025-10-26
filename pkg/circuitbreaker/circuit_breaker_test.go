package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_NewCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(3, 5*time.Second)

	if cb.failureThreshold != 3 {
		t.Errorf("Expected failure threshold 3, got %d", cb.failureThreshold)
	}
	if cb.timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", cb.timeout)
	}
	if cb.state != Closed {
		t.Errorf("Expected initial state Closed, got %v", cb.state)
	}
	if cb.failures != 0 {
		t.Errorf("Expected initial failures 0, got %d", cb.failures)
	}
}

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	cb := NewCircuitBreaker(3, 5*time.Second)

	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cb.state != Closed {
		t.Errorf("Expected state Closed, got %v", cb.state)
	}
	if cb.failures != 0 {
		t.Errorf("Expected failures 0, got %d", cb.failures)
	}
}

func TestCircuitBreaker_Execute_Failure(t *testing.T) {
	cb := NewCircuitBreaker(3, 5*time.Second)

	err := cb.Execute(func() error {
		return errors.New("test error")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if cb.state != Closed {
		t.Errorf("Expected state Closed after 1 failure, got %v", cb.state)
	}
	if cb.failures != 1 {
		t.Errorf("Expected failures 1, got %d", cb.failures)
	}
}

func TestCircuitBreaker_Execute_OpenState(t *testing.T) {
	cb := NewCircuitBreaker(2, 5*time.Second)

	// Cause 2 failures to open the circuit
	cb.Execute(func() error { return errors.New("error 1") })
	cb.Execute(func() error { return errors.New("error 2") })

	if cb.state != Open {
		t.Errorf("Expected state Open, got %v", cb.state)
	}

	// Try to execute when circuit is open
	err := cb.Execute(func() error {
		return nil
	})

	if err == nil {
		t.Error("Expected error when circuit is open, got nil")
	}
	if err.Error() != "circuit breaker is open" {
		t.Errorf("Expected 'circuit breaker is open', got '%s'", err.Error())
	}
}

func TestCircuitBreaker_Execute_HalfOpenToClosed(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)

	// Open the circuit
	cb.Execute(func() error { return errors.New("error 1") })
	cb.Execute(func() error { return errors.New("error 2") })

	if cb.state != Open {
		t.Errorf("Expected state Open, got %v", cb.state)
	}

	// Wait for timeout to pass
	time.Sleep(150 * time.Millisecond)

	// Execute should move to half-open and then close on success
	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cb.state != Closed {
		t.Errorf("Expected state Closed, got %v", cb.state)
	}
	if cb.failures != 0 {
		t.Errorf("Expected failures 0, got %d", cb.failures)
	}
}

func TestCircuitBreaker_Execute_HalfOpenToOpen(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)

	// Open the circuit
	cb.Execute(func() error { return errors.New("error 1") })
	cb.Execute(func() error { return errors.New("error 2") })

	// Wait for timeout to pass
	time.Sleep(150 * time.Millisecond)

	// Execute should move to half-open but fail again
	err := cb.Execute(func() error {
		return errors.New("error 3")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if cb.state != Open {
		t.Errorf("Expected state Open, got %v", cb.state)
	}
}

func TestCircuitBreaker_GetState(t *testing.T) {
	cb := NewCircuitBreaker(2, 5*time.Second)

	if cb.GetState() != Closed {
		t.Errorf("Expected initial state Closed, got %v", cb.GetState())
	}

	// Open the circuit
	cb.Execute(func() error { return errors.New("error 1") })
	cb.Execute(func() error { return errors.New("error 2") })

	if cb.GetState() != Open {
		t.Errorf("Expected state Open, got %v", cb.GetState())
	}
}

func TestCircuitBreaker_GetFailureCount(t *testing.T) {
	cb := NewCircuitBreaker(5, 5*time.Second)

	if cb.GetFailureCount() != 0 {
		t.Errorf("Expected initial failure count 0, got %d", cb.GetFailureCount())
	}

	cb.Execute(func() error { return errors.New("error 1") })
	cb.Execute(func() error { return errors.New("error 2") })

	if cb.GetFailureCount() != 2 {
		t.Errorf("Expected failure count 2, got %d", cb.GetFailureCount())
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	cb := NewCircuitBreaker(2, 5*time.Second)

	// Open the circuit
	cb.Execute(func() error { return errors.New("error 1") })
	cb.Execute(func() error { return errors.New("error 2") })

	if cb.state != Open {
		t.Errorf("Expected state Open, got %v", cb.state)
	}

	// Reset
	cb.Reset()

	if cb.state != Closed {
		t.Errorf("Expected state Closed after reset, got %v", cb.state)
	}
	if cb.failures != 0 {
		t.Errorf("Expected failures 0 after reset, got %d", cb.failures)
	}
}

func TestCircuitBreaker_ConcurrentAccess(t *testing.T) {
	cb := NewCircuitBreaker(10, 5*time.Second)

	// Test concurrent access
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			cb.Execute(func() error {
				return errors.New("concurrent error")
			})
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Circuit should be open after 10 failures
	if cb.state != Open {
		t.Errorf("Expected state Open after 10 failures, got %v", cb.state)
	}
}
