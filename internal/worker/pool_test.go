package worker_test

import (
	"github.com/antoniuk-oleksandr/order_processor/internal/worker"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type testTask struct {
}

func (t *testTask) Process() {
}

var _ = Describe("WorkerPool", Label("unit"), func() {
	When("creating a worker pool", func() {
		It("should create a worker pool", func() {
			numOfWorkers := 5
			buffer := 10

			pool, err := worker.NewWorkerPool(numOfWorkers, buffer)

			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a worker pool with valid parameters")
			Expect(pool).ToNot(BeNil(), "expected worker pool to be created with valid parameters")

			pool.Shutdown()
			pool.Wait()
		})
	})

	When("creating a worker pool with invalid number of workers", func() {
		It("should return an error", func() {
			numOfWorkers := -1
			buffer := 10

			pool, err := worker.NewWorkerPool(numOfWorkers, buffer)

			Expect(err).To(Equal(worker.ErrNumWorkersInvalid), "expected error when creating a worker pool with invalid number of workers")
			Expect(pool).To(BeNil(), "expected no worker pool to be created with invalid number of workers")
		})
	})

	When("creating a worker pool with invalid buffer", func() {
		It("should return an error", func() {
			numOfWorkers := 5
			buffer := -1

			pool, err := worker.NewWorkerPool(numOfWorkers, buffer)

			Expect(err).To(MatchError(worker.ErrBufferInvalid), "expected error when creating a worker pool with invalid buffer")
			Expect(pool).To(BeNil(), "expected no worker pool to be created with invalid buffer")
		})
	})

	When("adding a task to the pool", func() {
		It("should accept a task without error", func() {
			numOfWorkers := 5
			buffer := 1000
			ts := &testTask{}

			pool, err := worker.NewWorkerPool(numOfWorkers, buffer)
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a worker pool with valid parameters")
			Expect(pool).ToNot(BeNil(), "expected worker pool to be created with valid parameters")

			err = pool.AddTask(ts)
			Expect(err).ToNot(HaveOccurred(), "expected no error when adding a task to the pool")

			pool.Shutdown()
			pool.Wait()
		})
	})

	When("adding a task to a closed pool", func() {
		It("should return an error", func() {
			numOfWorkers := 5
			buffer := 1000
			ts := &testTask{}

			pool, err := worker.NewWorkerPool(numOfWorkers, buffer)
			Expect(err).ToNot(HaveOccurred(), "expected no error when creating a worker pool with valid parameters")
			Expect(pool).ToNot(BeNil(), "expected worker pool to be created with valid parameters")

			pool.Shutdown()
			pool.Wait()
			err = pool.AddTask(ts)

			Expect(err).To(MatchError(worker.ErrPoolClosed), "expected error when adding a task to a closed pool")
		})
	})
})
