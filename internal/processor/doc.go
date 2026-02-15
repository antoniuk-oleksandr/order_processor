// Package processor provides order processing functionality with user-specific queuing.
//
// The processor package implements an order processing system that maintains
// separate queues for each user to ensure orders from the same user are processed
// sequentially while allowing parallel processing across different users.
//
// Example usage:
//
//	storage := storage.NewStorage()
//	pool, _ := worker.NewWorkerPool(10, 100)
//	proc, err := processor.NewOrderProcessor(storage, pool)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer proc.Shutdown()
//
//	err = proc.Submit(order.Order{ID: 1, UserID: 1, Amount: 100})
//	if err != nil {
//		log.Println("failed to submit order:", err)
//	}
package processor
