package processor

import "errors"

var (
	// ErrProcessorShutdown is returned when attempting to submit orders to a shut down processor.
	ErrProcessorShutdown = errors.New("processor is shut down")
	// ErrStorageInvalid is returned when a nil storage is passed to NewOrderProcessor.
	ErrStorageInvalid = errors.New("storage must not be nil")
	// ErrWorkerPoolInvalid is returned when a nil worker pool is passed to NewOrderProcessor.
	ErrWorkerPoolInvalid = errors.New("worker pool must not be nil")
)
