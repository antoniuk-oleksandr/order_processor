package main

import (
	"fmt"
	"log"
	"github.com/antoniuk-oleksandr/order_processor/internal/order"
	"github.com/antoniuk-oleksandr/order_processor/internal/processor"
	"github.com/antoniuk-oleksandr/order_processor/internal/storage"
	"github.com/antoniuk-oleksandr/order_processor/internal/worker"
	"time"
)

func main() {
	s := storage.NewStorage()
	p, err := worker.NewWorkerPool(200, 1000)
	if err != nil {
		log.Fatal("err creating a worker pool:", err)
	}
	proc, err := processor.NewOrderProcessor(s, p)
	if err != nil {
		log.Fatal("err creating an order processor", err)
	}

	before := time.Now()
	for i := 1; i < 30; i++ {
		localErr := proc.Submit(order.Order{ID: i, UserID: 1, Amount: 100})
		if localErr != nil {
			// No errors are expected since processor is not shutdown yet
			log.Println("err:", localErr)
		}
	}

	for i := 1; i < 100; i++ {
		localErr := proc.Submit(order.Order{ID: i, UserID: i, Amount: 200})
		if localErr != nil {
			// No errors are expected since processor is not shutdown yet
			log.Println("err:", localErr)
		}
	}

	proc.Shutdown()
	elapsed := time.Since(before)
	fmt.Println("Processing time:", elapsed.Milliseconds(), "ms")

	err = proc.Submit(order.Order{ID: 11, UserID: 1, Amount: 100})
	if err != nil {
		// An error is expected since processor is already shutdown
		log.Println("err:", err)
	}

	balance, ok := proc.GetBalance(1)
	if ok {
		fmt.Println("User 1 balance:", balance)
	}
}
