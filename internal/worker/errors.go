package worker

import "errors"

var (
	// ErrNumWorkersInvalid is returned when the number of workers is less than or equal to 0.
	ErrNumWorkersInvalid = errors.New("number of workers must be greater than 0")
	// ErrBufferInvalid is returned when the buffer size is less than or equal to 0.
	ErrBufferInvalid = errors.New("buffer must be greater than 0")
	// ErrPoolClosed is returned when attempting to add a task to a closed worker pool.
	ErrPoolClosed = errors.New("worker pool is closed")
	// ErrNilTaskChannel is returned when the task channel is nil.
	ErrNilTaskChannel = errors.New("task channel cannot be nil")
)
