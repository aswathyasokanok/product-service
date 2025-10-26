package errors

import (
	"fmt"
	"testing"
)

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		errorType ErrorType
		expected  string
	}{
		{RetryableError, "RetryableError"},
		{NonRetryableError, "NonRetryableError"},
		{ValidationError, "ValidationError"},
		{SystemError, "SystemError"},
		{NetworkError, "NetworkError"},
		{TimeoutError, "TimeoutError"},
	}

	for _, test := range tests {
		result := fmt.Sprintf("%s", test.errorType)
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestClassifiedError_ShouldRetry(t *testing.T) {
	tests := []struct {
		errorType ErrorType
		expected  bool
	}{
		{RetryableError, true},
		{NonRetryableError, false},
		{ValidationError, false},
		{SystemError, false},
		{NetworkError, true},
		{TimeoutError, true},
	}

	for _, test := range tests {
		ce := &ClassifiedError{
			Type:    test.errorType,
			Message: "test error",
			Cause:   nil,
		}

		result := ce.ShouldRetry()
		if result != test.expected {
			t.Errorf("Expected ShouldRetry() to return %v for %s, got %v",
				test.expected, fmt.Sprintf("%s", test.errorType), result)
		}
	}
}

func TestClassifiedError_Error(t *testing.T) {
	ce := &ClassifiedError{
		Type:    ValidationError,
		Message: "test error message",
		Cause:   nil,
	}

	result := ce.Error()
	expected := "test error message"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestClassifiedError_Error_WithCause(t *testing.T) {
	cause := &ClassifiedError{
		Type:    SystemError,
		Message: "underlying error",
		Cause:   nil,
	}

	ce := &ClassifiedError{
		Type:    ValidationError,
		Message: "test error message",
		Cause:   cause,
	}

	result := ce.Error()
	expected := "test error message: underlying error"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestClassifiedError_Unwrap(t *testing.T) {
	cause := &ClassifiedError{
		Type:    SystemError,
		Message: "underlying error",
		Cause:   nil,
	}

	ce := &ClassifiedError{
		Type:    ValidationError,
		Message: "test error message",
		Cause:   cause,
	}

	unwrapped := ce.Unwrap()
	if unwrapped != cause {
		t.Errorf("Expected unwrapped error to be the cause")
	}
}

func TestClassifiedError_Unwrap_NoCause(t *testing.T) {
	ce := &ClassifiedError{
		Type:    ValidationError,
		Message: "test error message",
		Cause:   nil,
	}

	unwrapped := ce.Unwrap()
	if unwrapped != nil {
		t.Errorf("Expected unwrapped error to be nil when no cause")
	}
}

func TestNewRetryableError(t *testing.T) {
	err := NewRetryableError("retryable error message", nil)

	if err.Type != RetryableError {
		t.Errorf("Expected Type RetryableError, got %s", fmt.Sprintf("%s", err.Type))
	}
	if err.Message != "retryable error message" {
		t.Errorf("Expected Message 'retryable error message', got '%s'", err.Message)
	}
	if err.Cause != nil {
		t.Errorf("Expected Cause to be nil, got %v", err.Cause)
	}
	if !err.ShouldRetry() {
		t.Error("Expected ShouldRetry() to return true")
	}
}

func TestNewNonRetryableError(t *testing.T) {
	err := NewNonRetryableError("non-retryable error message", nil)

	if err.Type != NonRetryableError {
		t.Errorf("Expected Type NonRetryableError, got %s", fmt.Sprintf("%s", err.Type))
	}
	if err.Message != "non-retryable error message" {
		t.Errorf("Expected Message 'non-retryable error message', got '%s'", err.Message)
	}
	if err.Cause != nil {
		t.Errorf("Expected Cause to be nil, got %v", err.Cause)
	}
	if err.ShouldRetry() {
		t.Error("Expected ShouldRetry() to return false")
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("validation error message", nil)

	if err.Type != ValidationError {
		t.Errorf("Expected Type ValidationError, got %s", fmt.Sprintf("%s", err.Type))
	}
	if err.Message != "validation error message" {
		t.Errorf("Expected Message 'validation error message', got '%s'", err.Message)
	}
	if err.Cause != nil {
		t.Errorf("Expected Cause to be nil, got %v", err.Cause)
	}
	if err.ShouldRetry() {
		t.Error("Expected ShouldRetry() to return false")
	}
}

func TestNewSystemError(t *testing.T) {
	err := NewSystemError("system error message", nil)

	if err.Type != SystemError {
		t.Errorf("Expected Type SystemError, got %s", fmt.Sprintf("%s", err.Type))
	}
	if err.Message != "system error message" {
		t.Errorf("Expected Message 'system error message', got '%s'", err.Message)
	}
	if err.Cause != nil {
		t.Errorf("Expected Cause to be nil, got %v", err.Cause)
	}
	if err.ShouldRetry() {
		t.Error("Expected ShouldRetry() to return false")
	}
}

func TestNewRetryableErrorWithCause(t *testing.T) {
	cause := &ClassifiedError{
		Type:    SystemError,
		Message: "underlying error",
		Cause:   nil,
	}

	err := NewRetryableErrorWithCause("retryable error message", cause)

	if err.Type != RetryableError {
		t.Errorf("Expected Type RetryableError, got %s", err.Type.String())
	}
	if err.Message != "retryable error message" {
		t.Errorf("Expected Message 'retryable error message', got '%s'", err.Message)
	}
	if err.Cause != cause {
		t.Errorf("Expected Cause to be the provided cause")
	}
	if !err.ShouldRetry() {
		t.Error("Expected ShouldRetry() to return true")
	}
}

func TestNewNonRetryableErrorWithCause(t *testing.T) {
	cause := &ClassifiedError{
		Type:    ValidationError,
		Message: "underlying error",
		Cause:   nil,
	}

	err := NewNonRetryableErrorWithCause("non-retryable error message", cause)

	if err.Type != NonRetryableError {
		t.Errorf("Expected Type NonRetryableError, got %s", err.Type.String())
	}
	if err.Message != "non-retryable error message" {
		t.Errorf("Expected Message 'non-retryable error message', got '%s'", err.Message)
	}
	if err.Cause != cause {
		t.Errorf("Expected Cause to be the provided cause")
	}
	if err.ShouldRetry() {
		t.Error("Expected ShouldRetry() to return false")
	}
}
