package queue

import "errors"

// Common queue errors
var (
	ErrQueueFull          = errors.New("queue is full")
	ErrQueueClosed        = errors.New("queue is closed")
	ErrBatchProcessorFull = errors.New("batch processor is full")
	ErrInvalidEvent       = errors.New("invalid event")
	ErrEventTooLarge      = errors.New("event too large")
)
