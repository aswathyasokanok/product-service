package retry

import (
	"errors"
	"testing"
	"time"
)

func TestRetryConfig_DefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxAttempts != 3 {
		t.Errorf("Expected MaxAttempts 3, got %d", config.MaxAttempts)
	}
	if config.InitialDelay != 100*time.Millisecond {
		t.Errorf("Expected InitialDelay 100ms, got %v", config.InitialDelay)
	}
	if config.MaxDelay != 30*time.Second {
		t.Errorf("Expected MaxDelay 30s, got %v", config.MaxDelay)
	}
	if config.Multiplier != 2.0 {
		t.Errorf("Expected Multiplier 2.0, got %f", config.Multiplier)
	}
}

func TestRetryConfig_ExecuteWithRetry_Success(t *testing.T) {
	config := &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}

	attempts := 0
	err := config.ExecuteWithRetry(func() error {
		attempts++
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestRetryConfig_ExecuteWithRetry_Failure(t *testing.T) {
	config := &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}

	attempts := 0
	err := config.ExecuteWithRetry(func() error {
		attempts++
		return errors.New("test error")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
	if err.Error() != "operation failed after 3 attempts" {
		t.Errorf("Expected 'operation failed after 3 attempts', got '%s'", err.Error())
	}
}

func TestRetryConfig_ExecuteWithRetry_SuccessOnRetry(t *testing.T) {
	config := &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}

	attempts := 0
	err := config.ExecuteWithRetry(func() error {
		attempts++
		if attempts == 2 {
			return nil // Success on second attempt
		}
		return errors.New("test error")
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestRetryConfig_ExecuteWithRetryAndCallback(t *testing.T) {
	config := &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}

	attempts := 0
	callbackCalls := 0

	err := config.ExecuteWithRetryAndCallback(
		func() error {
			attempts++
			return errors.New("test error")
		},
		func(attempt int, err error) {
			callbackCalls++
		},
	)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
	if callbackCalls != 3 { // Callback called for all 3 failures
		t.Errorf("Expected 3 callback calls, got %d", callbackCalls)
	}
}

func TestRetryConfig_ExecuteWithRetryAndCallback_Success(t *testing.T) {
	config := &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}

	attempts := 0
	callbackCalls := 0

	err := config.ExecuteWithRetryAndCallback(
		func() error {
			attempts++
			if attempts == 2 {
				return nil // Success on second attempt
			}
			return errors.New("test error")
		},
		func(attempt int, err error) {
			callbackCalls++
		},
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
	if callbackCalls != 1 { // Callback called for first failure only
		t.Errorf("Expected 1 callback call, got %d", callbackCalls)
	}
}

func TestRetryConfig_ExecuteWithRetryAndCallback_NilCallback(t *testing.T) {
	config := &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     100 * time.Millisecond,
		Multiplier:   2.0,
	}

	attempts := 0

	err := config.ExecuteWithRetryAndCallback(
		func() error {
			attempts++
			return errors.New("test error")
		},
		nil, // nil callback
	)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryConfig_DelayCalculation(t *testing.T) {
	config := &RetryConfig{
		MaxAttempts:  5,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     50 * time.Millisecond,
		Multiplier:   2.0,
	}

	start := time.Now()
	attempts := 0

	err := config.ExecuteWithRetry(func() error {
		attempts++
		return errors.New("test error")
	})

	elapsed := time.Since(start)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 5 {
		t.Errorf("Expected 5 attempts, got %d", attempts)
	}

	// Should have delays: 10ms, 20ms, 40ms, 50ms (capped at MaxDelay)
	expectedMinDelay := 10*time.Millisecond + 20*time.Millisecond + 40*time.Millisecond + 50*time.Millisecond
	if elapsed < expectedMinDelay {
		t.Errorf("Expected elapsed time >= %v, got %v", expectedMinDelay, elapsed)
	}
}

func TestRetryConfig_MaxDelayCap(t *testing.T) {
	config := &RetryConfig{
		MaxAttempts:  10,
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     50 * time.Millisecond,
		Multiplier:   2.0,
	}

	start := time.Now()
	attempts := 0

	err := config.ExecuteWithRetry(func() error {
		attempts++
		return errors.New("test error")
	})

	elapsed := time.Since(start)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 10 {
		t.Errorf("Expected 10 attempts, got %d", attempts)
	}

	// Should not exceed reasonable time even with many attempts
	maxReasonableTime := 1 * time.Second
	if elapsed > maxReasonableTime {
		t.Errorf("Expected elapsed time <= %v, got %v", maxReasonableTime, elapsed)
	}
}
