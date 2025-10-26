package errors

import "fmt"

// ErrorType represents the type of error for classification
type ErrorType int

const (
	RetryableError ErrorType = iota
	NonRetryableError
	ValidationError
	SystemError
	NetworkError
	TimeoutError
)

// String returns the string representation of the ErrorType
func (et ErrorType) String() string {
	switch et {
	case RetryableError:
		return "RetryableError"
	case NonRetryableError:
		return "NonRetryableError"
	case ValidationError:
		return "ValidationError"
	case SystemError:
		return "SystemError"
	case NetworkError:
		return "NetworkError"
	case TimeoutError:
		return "TimeoutError"
	default:
		return "UnknownError"
	}
}

// ClassifiedError represents an error with classification
type ClassifiedError struct {
	Type    ErrorType
	Message string
	Cause   error
}

// Error implements the error interface
func (ce *ClassifiedError) Error() string {
	if ce.Cause != nil {
		return fmt.Sprintf("%s: %v", ce.Message, ce.Cause)
	}
	return ce.Message
}

// Unwrap returns the underlying error
func (ce *ClassifiedError) Unwrap() error {
	return ce.Cause
}

// ShouldRetry returns true if the error is retryable
func (ce *ClassifiedError) ShouldRetry() bool {
	return ce.Type == RetryableError || ce.Type == NetworkError || ce.Type == TimeoutError
}

// IsValidationError returns true if the error is a validation error
func (ce *ClassifiedError) IsValidationError() bool {
	return ce.Type == ValidationError
}

// IsSystemError returns true if the error is a system error
func (ce *ClassifiedError) IsSystemError() bool {
	return ce.Type == SystemError
}

// NewClassifiedError creates a new classified error
func NewClassifiedError(errorType ErrorType, message string, cause error) *ClassifiedError {
	return &ClassifiedError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// NewRetryableError creates a new retryable error
func NewRetryableError(message string, cause error) *ClassifiedError {
	return NewClassifiedError(RetryableError, message, cause)
}

// NewNonRetryableError creates a new non-retryable error
func NewNonRetryableError(message string, cause error) *ClassifiedError {
	return NewClassifiedError(NonRetryableError, message, cause)
}

// NewValidationError creates a new validation error
func NewValidationError(message string, cause error) *ClassifiedError {
	return NewClassifiedError(ValidationError, message, cause)
}

// NewSystemError creates a new system error
func NewSystemError(message string, cause error) *ClassifiedError {
	return NewClassifiedError(SystemError, message, cause)
}

// NewNetworkError creates a new network error
func NewNetworkError(message string, cause error) *ClassifiedError {
	return NewClassifiedError(NetworkError, message, cause)
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(message string, cause error) *ClassifiedError {
	return NewClassifiedError(TimeoutError, message, cause)
}

// NewRetryableErrorWithCause creates a new retryable error with a specific cause
func NewRetryableErrorWithCause(message string, cause error) *ClassifiedError {
	return NewClassifiedError(RetryableError, message, cause)
}

// NewNonRetryableErrorWithCause creates a new non-retryable error with a specific cause
func NewNonRetryableErrorWithCause(message string, cause error) *ClassifiedError {
	return NewClassifiedError(NonRetryableError, message, cause)
}
