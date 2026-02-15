package processor

import (
	"github.com/antoniuk-oleksandr/order_processor/internal/order"
	"github.com/antoniuk-oleksandr/order_processor/internal/storage"
	"github.com/antoniuk-oleksandr/order_processor/internal/worker"
	"sync"
)

// OrderProcessor handles order submission and processing with user-specific queuing.
// Orders from the same user are processed sequentially to maintain consistency,
// while orders from different users are processed concurrently.
type OrderProcessor interface {
	// Submit adds an order to the processing queue.
	// Returns ErrProcessorShutdown if the processor has been shut down.
	Submit(order order.Order) error
	// Shutdown gracefully shuts down the processor and waits for all orders to be processed.
	Shutdown()
	// GetBalance retrieves the current balance for a user.
	// Returns the balance and true if the user exists, or 0 and false if not found.
	GetBalance(userID int) (int, bool)
}

type orderProcessor struct {
	storage      storage.Storage
	workerPool   worker.WorkerPool
	shutdownOnce sync.Once
	userQueues   map[int]chan order.Order
	userQueuesMu sync.Mutex
	shutdownChan chan struct{}
	processorWg  sync.WaitGroup
}

// NewOrderProcessor creates a new OrderProcessor with the given storage and worker pool.
// Returns ErrStorageInvalid if storage is nil, or ErrWorkerPoolInvalid if workerPool is nil.
func NewOrderProcessor(storage storage.Storage, workerPool worker.WorkerPool) (OrderProcessor, error) {
	if storage == nil {
		return nil, ErrStorageInvalid
	}

	if workerPool == nil {
		return nil, ErrWorkerPoolInvalid
	}

	return &orderProcessor{
		storage:      storage,
		workerPool:   workerPool,
		shutdownOnce: sync.Once{},
		userQueues:   make(map[int]chan order.Order),
		shutdownChan: make(chan struct{}),
	}, nil
}

func (o *orderProcessor) GetBalance(userID int) (int, bool) {
	return o.storage.Get(userID)
}

func (o *orderProcessor) Shutdown() {
	o.shutdownOnce.Do(func() {
		close(o.shutdownChan)

		o.userQueuesMu.Lock()
		for _, queue := range o.userQueues {
			close(queue)
		}
		o.userQueuesMu.Unlock()

		o.processorWg.Wait()

		o.workerPool.Shutdown()
		o.workerPool.Wait()
	})
}

func (o *orderProcessor) Submit(ord order.Order) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrProcessorShutdown
		}
	}()

	o.userQueuesMu.Lock()
	queue, exists := o.userQueues[ord.UserID]
	if !exists {
		queue = make(chan order.Order, 100)
		o.userQueues[ord.UserID] = queue

		queueTask := NewOrderTask(queue, o.storage)
		if err := o.workerPool.AddTask(queueTask); err != nil {
			o.userQueuesMu.Unlock()
			return ErrProcessorShutdown
		}
	}
	o.userQueuesMu.Unlock()

	select {
	case queue <- ord:
		return nil
	case <-o.shutdownChan:
		return ErrProcessorShutdown
	}
}
