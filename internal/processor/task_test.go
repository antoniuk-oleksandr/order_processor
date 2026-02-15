package processor_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"

	"github.com/antoniuk-oleksandr/order_processor/internal/order"
	"github.com/antoniuk-oleksandr/order_processor/internal/processor"
	"github.com/antoniuk-oleksandr/order_processor/mock"
)

var _ = Describe("Processor", Label("unit"), func() {
	When("processing an order task", func() {
		It("should call the storage Add method with correct parameters", func() {
			chLength := 2
			o := order.Order{
				UserID: 1,
				Amount: 100,
			}

			ctrl := gomock.NewController(GinkgoT())
			s := mock.NewMockUserStorage(ctrl)

			s.EXPECT().Add(o.UserID, o.Amount).Times(chLength)

			orders := make(chan order.Order, chLength)
			for range chLength {
				orders <- o
			}
			close(orders)

			task := processor.NewOrderTask(orders, s)
			task.Process()
		})
	})
})
