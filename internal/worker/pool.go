package worker

import (
	"sync"
	"sync/atomic"
)

// Task represents a unit of work that can be processed by a worker.
type Task interface {
	// Process executes the task's work.
	Process()
}

// WorkerPool manages a pool of workers that process tasks concurrently.
type WorkerPool interface {
	// AddTask adds a task to the pool for processing.
	// Returns ErrPoolClosed if the pool has been shut down.
	AddTask(task Task) error
	// Shutdown gracefully shuts down the pool, preventing new tasks from being added.
	Shutdown()
	// Wait blocks until all workers have finished processing their tasks.
	Wait()
}

type workerPool struct {
	numOfWorkers int
	tasks        chan Task
	wg           sync.WaitGroup
	shutdownOnce sync.Once
	closed       atomic.Bool
}

// NewWorkerPool creates a new WorkerPool with the specified number of workers and buffer size.
// Returns ErrNumWorkersInvalid if numOfWorkers <= 0, or ErrBufferInvalid if buffer <= 0.
func NewWorkerPool(numOfWorkers, buffer int) (WorkerPool, error) {
	if numOfWorkers <= 0 {
		return nil, ErrNumWorkersInvalid
	}

	if buffer <= 0 {
		return nil, ErrBufferInvalid
	}

	w := &workerPool{
		tasks:        make(chan Task, buffer),
		numOfWorkers: numOfWorkers,
		wg:           sync.WaitGroup{},
		shutdownOnce: sync.Once{},
		closed:       atomic.Bool{},
	}

	for range numOfWorkers {
		w.wg.Go(func() {
			worker := NewWorker(w.tasks)
			worker.Run()
		})
	}

	return w, nil
}

func (w *workerPool) AddTask(task Task) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrPoolClosed
		}
	}()

	if w.closed.Load() {
		return ErrPoolClosed
	}

	w.tasks <- task
	return nil
}

func (w *workerPool) Shutdown() {
	w.shutdownOnce.Do(func() {
		w.closed.Store(true)
		close(w.tasks)
	})
}

func (w *workerPool) Wait() {
	w.wg.Wait()
}
