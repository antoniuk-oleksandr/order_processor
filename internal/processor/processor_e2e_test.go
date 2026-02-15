package processor_test

import (
	"github.com/antoniuk-oleksandr/order_processor/internal/order"
	"github.com/antoniuk-oleksandr/order_processor/internal/processor"
	"github.com/antoniuk-oleksandr/order_processor/internal/storage"
	"github.com/antoniuk-oleksandr/order_processor/internal/worker"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OrderProcessor E2E", Label("e2e"), func() {
	It("should process multiple orders concurrently and update balances correctly", func() {
		numOfWorker := 5
		buffer := 100
		
		s := storage.NewStorage()
		pool, err := worker.NewWorkerPool(numOfWorker, buffer)
		Expect(err).NotTo(HaveOccurred(), "creating worker pool should not return an error")
		processor, err := processor.NewOrderProcessor(s, pool)
		Expect(err).NotTo(HaveOccurred(), "creating processor should not return an error")

		orders := []order.Order{
			{ID: 1, UserID: 1, Amount: 100},
			{ID: 2, UserID: 2, Amount: 200},
			{ID: 3, UserID: 1, Amount: 50},
			{ID: 4, UserID: 3, Amount: 300},
		}

		var wg sync.WaitGroup
		for _, o := range orders {
			wg.Go(func() {
				_ = processor.Submit(o)
			})
		}

		wg.Wait()

		processor.Shutdown()

		amount, ok := s.Get(1)
		Expect(ok).To(BeTrue(), "user 1 should exist in storage")
		Expect(amount).To(Equal(150), "user 1 balance should be 150")

		amount, ok = s.Get(2)
		Expect(ok).To(BeTrue(), "user 2 should exist in storage")
		Expect(amount).To(Equal(200), "user 2 balance should be 200")

		amount, ok = s.Get(3)
		Expect(ok).To(BeTrue(), "user 3 should exist in storage")
		Expect(amount).To(Equal(300), "user 3 balance should be 300")
	})

	It("should return an error if submitting after shutdown", func() {
		numOfWorker := 2
		buffer := 10
		
		s := storage.NewStorage()
		pool, err := worker.NewWorkerPool(numOfWorker, buffer)
		Expect(err).NotTo(HaveOccurred(), "creating worker pool should not return an error")

		proc, err := processor.NewOrderProcessor(s, pool)
		Expect(err).NotTo(HaveOccurred(), "creating processor should not return an error")

		proc.Shutdown()

		err = proc.Submit(order.Order{ID: 1, UserID: 1, Amount: 100})
		Expect(err).To(MatchError(processor.ErrProcessorShutdown))
	})

	It("should handle high concurrency without losing orders", func() {
		numOfWorker := 20
		buffer := 1000
		
		s := storage.NewStorage()
		pool, err := worker.NewWorkerPool(numOfWorker, buffer)
		Expect(err).NotTo(HaveOccurred(), "creating worker pool should not return an error")
		proc, err := processor.NewOrderProcessor(s, pool)
		Expect(err).NotTo(HaveOccurred(), "creating processor should not return an error")

		numUsers := 50
		numOrdersPerUser := 100

		var wg sync.WaitGroup
		for u := 1; u <= numUsers; u++ {
			for i := range numOrdersPerUser {
				userID := u
				amount := i + 1
				wg.Go(func() {
					defer GinkgoRecover()
					localErr := proc.Submit(order.Order{UserID: userID, Amount: amount})
					Expect(localErr).NotTo(HaveOccurred(), "submitting order should not return an error")
				})
			}
		}

		wg.Wait()
		
		proc.Shutdown()
		for u := 1; u <= numUsers; u++ {
			amount, ok := s.Get(u)
			Expect(ok).To(BeTrue(), "user %d should exist in storage", u)
			expected := (numOrdersPerUser * (numOrdersPerUser + 1)) / 2
			Expect(amount).To(Equal(expected), "user %d balance should be correct", u)
		}
	})
})
