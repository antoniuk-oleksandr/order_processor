// Package worker provides a worker pool implementation for concurrent task processing.
//
// The worker package offers a thread-safe worker pool that manages a fixed number
// of worker goroutines. Tasks are submitted to a buffered channel and processed
// concurrently by available workers.
//
// Example usage:
//
//	pool, err := worker.NewWorkerPool(10, 100)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer pool.Shutdown()
//	defer pool.Wait()
//
//	err = pool.AddTask(myTask)
//	if err != nil {
//		log.Println("failed to add task:", err)
//	}
package worker
