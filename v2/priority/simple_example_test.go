package priority_test

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/akramarenkov/cqos/v2/priority"
)

func ExampleSimple() {
	handlersQuantity := 100
	// Preferably input channels should be buffered
	inputCapacity := 10
	itemsQuantity := 100

	inputs := map[uint]chan string{
		3: make(chan string, inputCapacity),
		2: make(chan string, inputCapacity),
		1: make(chan string, inputCapacity),
	}

	// Map key is a value of priority
	inputsOpts := map[uint]<-chan string{
		3: inputs[3],
		2: inputs[2],
		1: inputs[1],
	}

	defer func() {
		for _, input := range inputs {
			close(input)
		}
	}()

	// Used only in this example for detect that all written data are processed
	measurements := make(chan bool)
	defer close(measurements)

	handle := func(ctx context.Context, item string) {
		// Data processing
		// fmt.Println(item)
		select {
		case <-ctx.Done():
		case measurements <- true:
		}
	}

	// For equaling use FairDivider, for prioritization use RateDivider or custom divider
	opts := priority.SimpleOpts[string]{
		Divider:          priority.RateDivider,
		Handle:           handle,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
	}

	simple, err := priority.NewSimple(opts)
	if err != nil {
		panic(err)
	}

	defer simple.Stop()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	// Run writers, that write data to input channels
	for priority, input := range inputs {
		wg.Add(1)

		go func(precedency uint, channel chan string) {
			defer wg.Done()

			base := strconv.Itoa(int(precedency))

			for id := 0; id < itemsQuantity; id++ {
				item := base + ":" + strconv.Itoa(id)

				channel <- item
			}
		}(priority, input)
	}

	received := 0

	// Wait for process all written data
	for range measurements {
		received++

		if received == itemsQuantity*len(inputs) {
			break
		}
	}

	fmt.Println("Processed items quantity:", received)
	// Output: Processed items quantity: 300
}

func ExampleSimple_GracefulStop() {
	handlersQuantity := 100
	// Preferably input channels should be buffered
	inputCapacity := 10
	itemsQuantity := 100

	inputs := map[uint]chan string{
		3: make(chan string, inputCapacity),
		2: make(chan string, inputCapacity),
		1: make(chan string, inputCapacity),
	}

	// Map key is a value of priority
	inputsOpts := map[uint]<-chan string{
		3: inputs[3],
		2: inputs[2],
		1: inputs[1],
	}

	// Used only in this example for detect that all written data are processed
	measurements := make(chan bool)

	handle := func(ctx context.Context, item string) {
		// Data processing
		// fmt.Println(item)
		select {
		case <-ctx.Done():
		case measurements <- true:
		}
	}

	// For equaling use FairDivider, for prioritization use RateDivider or custom divider
	opts := priority.SimpleOpts[string]{
		Divider:          priority.RateDivider,
		Handle:           handle,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
	}

	simple, err := priority.NewSimple(opts)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}

	// Run writers, that write data to input channels
	for priority, input := range inputs {
		wg.Add(1)

		go func(precedency uint, channel chan string) {
			defer wg.Done()

			base := strconv.Itoa(int(precedency))

			for id := 0; id < itemsQuantity; id++ {
				item := base + ":" + strconv.Itoa(id)

				channel <- item
			}
		}(priority, input)
	}

	obtained := make(chan int)
	defer close(obtained)

	// Counting the amount of received data
	go func() {
		received := 0

		for range measurements {
			received++
		}

		obtained <- received
	}()

	// You must end write to input channels and close them (or remove),
	// otherwise graceful stop not be ended
	wg.Wait()

	for _, input := range inputs {
		close(input)
	}

	simple.GracefulStop()

	// Terminate measurements
	close(measurements)

	received := <-obtained

	// Verify data received from discipline
	if received != itemsQuantity*len(inputs) {
		panic("graceful stop work not properly")
	}

	fmt.Println("Processed items quantity:", received)
	// Output: Processed items quantity: 300
}
