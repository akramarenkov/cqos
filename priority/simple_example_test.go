package priority_test

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/akramarenkov/cqos/priority"
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

	// Used only in this example for detect that all written data are processed
	measures := make(chan string)
	defer close(measures)

	handle := func(ctx context.Context, item string) {
		// Data processing
		select {
		case <-ctx.Done():
		case measures <- item:
		}
	}

	// For equaling use FairDivider, for prioritization use
	// RateDivider or custom divider
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
			defer close(channel)

			base := strconv.Itoa(int(precedency))

			for id := 0; id < itemsQuantity; id++ {
				item := base + ":" + strconv.Itoa(id)

				channel <- item
			}
		}(priority, input)
	}

	received := 0

	// Wait for process all written data
	for range measures {
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
	measures := make(chan string)
	defer close(measures)

	handle := func(ctx context.Context, item string) {
		// Data processing
		select {
		case <-ctx.Done():
		case measures <- item:
		}
	}

	// For equaling use FairDivider, for prioritization use
	// RateDivider or custom divider
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
	defer wg.Wait()

	// Run writers, that write data to input channels
	for priority, input := range inputs {
		wg.Add(1)

		go func(precedency uint, channel chan string) {
			defer wg.Done()
			defer close(channel)

			base := strconv.Itoa(int(precedency))

			for id := 0; id < itemsQuantity; id++ {
				item := base + ":" + strconv.Itoa(id)

				channel <- item
			}
		}(priority, input)
	}

	obtained := make(chan int)

	go func() {
		defer close(obtained)

		received := 0

		// Wait for process all written data
		for range measures {
			received++

			if received == itemsQuantity*len(inputs) {
				obtained <- received
				return
			}
		}
	}()

	simple.GracefulStop()

	fmt.Println("Processed items quantity:", <-obtained)
	// Output: Processed items quantity: 300
}
