package worker_test

import (
	"sync/atomic"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/antoniuk-oleksandr/order_processor/internal/worker"
)

type atomicTestTask struct {
	processed *atomic.Bool
}

func (t *atomicTestTask) Process() {
	t.processed.Store(true)
}

var _ = Describe("Worker", Label("unit"), func() {
	When("creating a new worker", func() {
		It("should create successfully with a valid task channel", func() {
			tasks := make(chan worker.Task, 2)

			w := worker.NewWorker(tasks)

			Expect(w).ToNot(BeNil(), "expected worker to be created")
		})
	})

	When("running a worker directly", func() {
		It("should process all tasks from the channel", func() {
			tasks := make(chan worker.Task, 2)

			processedFlags := make([]*atomic.Bool, 2)
			for i := range 2 {
				flag := &atomic.Bool{}
				processedFlags[i] = flag
				tasks <- &atomicTestTask{processed: flag}
			}

			close(tasks)

			w := worker.NewWorker(tasks)
			w.Run()

			for _, flag := range processedFlags {
				Expect(flag.Load()).To(BeTrue(), "expected task to be processed")
			}
		})
	})
})
