package retry

import (
	"fmt"
	"time"
)

// RetryConfig defines the configuration for retry operations
type RetryConfig struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
	}
}

// ExecuteWithRetry executes an operation with exponential backoff retry
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

// ExecuteWithRetryAndCallback executes an operation with retry and calls a callback on each failure
func (r *RetryConfig) ExecuteWithRetryAndCallback(operation func() error, onFailure func(attempt int, err error)) error {
	delay := r.InitialDelay

	for attempt := 1; attempt <= r.MaxAttempts; attempt++ {
		if err := operation(); err == nil {
			return nil
		}

		if onFailure != nil {
			onFailure(attempt, fmt.Errorf("attempt %d failed", attempt))
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
