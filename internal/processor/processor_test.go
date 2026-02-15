package processor_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/antoniuk-oleksandr/order_processor/internal/order"
	"github.com/antoniuk-oleksandr/order_processor/internal/processor"
	"github.com/antoniuk-oleksandr/order_processor/internal/storage"
	"github.com/antoniuk-oleksandr/order_processor/internal/worker"
	"github.com/antoniuk-oleksandr/order_processor/mock"
)

var _ = Describe("OrderProcessor", Label("unit"), func() {
	When("creating a new order processor", func() {
		It("should create successfully with valid storage and worker pool", func() {
			numOfWorker := 10
			buffer := 1000
			pool, err := worker.NewWorkerPool(numOfWorker, buffer)
			Expect(err).NotTo(HaveOccurred())
			s := storage.NewStorage()

			proc, err := processor.NewOrderProcessor(s, pool)

			Expect(proc).NotTo(BeNil(), "processor should not be nil")
			Expect(err).NotTo(HaveOccurred(), "creating processor should not return an error")
		})

		It("should return an error when creating with nil storage", func() {
			numOfWorker := 10
			buffer := 1000
			pool, err := worker.NewWorkerPool(numOfWorker, buffer)
			Expect(err).NotTo(HaveOccurred())

			proc, err := processor.NewOrderProcessor(nil, pool)

			Expect(proc).To(BeNil(), "processor should be nil when storage is nil")
			Expect(err).To(HaveOccurred(), "creating processor with nil storage should return an error")
		})

		It("should return an error when creating with nil worker pool", func() {
			s := storage.NewStorage()

			proc, err := processor.NewOrderProcessor(s, nil)
			Expect(proc).To(BeNil(), "processor should be nil when worker pool is nil")
			Expect(err).To(HaveOccurred(), "creating processor with nil worker pool should return an error")
		})
	})

	When("submitting an order", func() {
		It("should submit successfully and add a task to the worker pool", func() {
			order := order.Order{
				UserID: 1,
				Amount: 100,
			}

			ctrl := gomock.NewController(GinkgoT())
			s := mock.NewMockUserStorage(ctrl)
			pool := mock.NewMockWorkerPool(ctrl)
			pool.EXPECT().AddTask(gomock.Any()).Return(nil).Times(1)

			proc, err := processor.NewOrderProcessor(s, pool)
			Expect(proc).NotTo(BeNil(), "processor should not be nil")
			Expect(err).NotTo(HaveOccurred(), "creating processor should not return an error")

			err = proc.Submit(order)
			Expect(err).NotTo(HaveOccurred(), "submitting order should not return an error")
		})

		It("should return an error when worker pool fails to add a task", func() {
			order := order.Order{
				UserID: 1,
				Amount: 100,
			}

			ctrl := gomock.NewController(GinkgoT())
			s := mock.NewMockUserStorage(ctrl)

			pool := mock.NewMockWorkerPool(ctrl)
			pool.EXPECT().AddTask(gomock.Any()).Return(worker.ErrPoolClosed).Times(1)

			proc, err := processor.NewOrderProcessor(s, pool)
			Expect(proc).NotTo(BeNil(), "processor should not be nil")
			Expect(err).NotTo(HaveOccurred(), "creating processor should not return an error")

			err = proc.Submit(order)
			Expect(err).To(MatchError(processor.ErrProcessorShutdown), "submitting order after shutdown should return ErrProcessorShutdown")
		})
	})

	When("getting user balance", func() {
		It("should return the correct balance from storage", func() {
			amout := 100
			userId := 1
			success := true
			order := order.Order{
				UserID: userId,
				Amount: amout,
			}

			ctrl := gomock.NewController(GinkgoT())
			pool := mock.NewMockWorkerPool(ctrl)

			s := mock.NewMockUserStorage(ctrl)
			s.EXPECT().Get(order.UserID).Return(amout, success).Times(1)

			proc, err := processor.NewOrderProcessor(s, pool)
			Expect(proc).NotTo(BeNil(), "processor should not be nil")
			Expect(err).NotTo(HaveOccurred(), "creating processor should not return an error")

			result, ok := proc.GetBalance(userId)
			Expect(result).To(Equal(amout), "GetBalance should return the correct balance from storage")
			Expect(ok).To(BeTrue(), "GetBalance should return true when balance is found")
		})

		It("should return zero and false when user is not found in storage", func() {
			userId := 1

			ctrl := gomock.NewController(GinkgoT())
			pool := mock.NewMockWorkerPool(ctrl)

			s := mock.NewMockUserStorage(ctrl)
			s.EXPECT().Get(userId).Return(0, false).Times(1)

			proc, err := processor.NewOrderProcessor(s, pool)
			Expect(proc).NotTo(BeNil(), "processor should not be nil")
			Expect(err).NotTo(HaveOccurred(), "creating processor should not return an error")

			result, ok := proc.GetBalance(userId)
			Expect(result).To(Equal(0), "GetBalance should return zero when user is not found in storage")
			Expect(ok).To(BeFalse(), "GetBalance should return false when user is not found in storage")
		})
	})

	When("shutting down the processor", func() {
		It("should shut down the worker pool successfully", func() {
			ctrl := gomock.NewController(GinkgoT())
			s := mock.NewMockUserStorage(ctrl)

			pool := mock.NewMockWorkerPool(ctrl)
			pool.EXPECT().Shutdown().Times(1)
			pool.EXPECT().Wait().Times(1)

			proc, err := processor.NewOrderProcessor(s, pool)
			Expect(proc).NotTo(BeNil(), "processor should not be nil")
			Expect(err).NotTo(HaveOccurred(), "creating processor should not return an error")

			proc.Shutdown()
		})
	})
})
