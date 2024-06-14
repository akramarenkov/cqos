package priority_test

import (
	"context"
	"fmt"
	"sync"

	"github.com/akramarenkov/cqos/priority"
)

func ExampleSimple() {
	handlersQuantity := uint(100)
	itemsQuantity := 10000
	// Preferably, input channels should be buffered for performance reasons
	inputCapacity := 10

	inputs := map[uint]chan int{
		70: make(chan int, inputCapacity),
		20: make(chan int, inputCapacity),
		10: make(chan int, inputCapacity),
	}

	// Map key is a value of priority
	inputsOpts := make(map[uint]<-chan int, len(inputs))

	for priority, channel := range inputs {
		inputsOpts[priority] = channel
	}

	// Used only in this example for measuring input data
	measurements := make(chan int)

	handle := func(ctx context.Context, item int) {
		// Data processing
		select {
		case <-ctx.Done():
		case measurements <- item:
		}
	}

	// For equaling use FairDivider, for prioritization use
	// RateDivider or custom divider
	opts := priority.SimpleOpts[int]{
		Divider:          priority.RateDivider,
		Handle:           handle,
		HandlersQuantity: handlersQuantity,
		Inputs:           inputsOpts,
	}

	simple, err := priority.NewSimple(opts)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	// Running writers, that write data to input channels
	for _, input := range inputs {
		wg.Add(1)

		go func(channel chan int) {
			defer wg.Done()
			defer close(channel)

			for id := range itemsQuantity {
				channel <- id
			}
		}(input)
	}

	// For simplicity, the process of graceful termination of the discipline is
	// starts immediately
	go simple.GracefulStop()

	// Waiting for the completion of the discipline
	go func() {
		defer close(measurements)

		for err := range simple.Err() {
			if err != nil {
				fmt.Println("An error was received: ", err)
			}
		}
	}()

	received := 0

	// Receiving the measurements data
	for range measurements {
		received++
	}

	fmt.Println("Processed data items quantity:", received)
	// Output: Processed data items quantity: 30000
}
