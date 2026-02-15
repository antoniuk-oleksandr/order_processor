package processor

import (
	"github.com/antoniuk-oleksandr/order_processor/internal/order"
	"github.com/antoniuk-oleksandr/order_processor/internal/storage"
	"time"
)

type orderTask interface {
	Process()
}

type orderTaskStr struct {
	queue   chan order.Order
	storage storage.Storage
}

func NewOrderTask(queue chan order.Order, storage storage.Storage) orderTask {
	return &orderTaskStr{
		queue:   queue,
		storage: storage,
	}
}

func (o orderTaskStr) Process() {
	for ord := range o.queue {
		time.Sleep(time.Millisecond * 200)

		o.storage.Add(ord.UserID, ord.Amount)
	}
}
