package simple_test

import (
	"fmt"

	"github.com/akramarenkov/cqos/v2/priority/divider"
	"github.com/akramarenkov/cqos/v2/priority/simple"
)

func ExampleDiscipline() {
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

	handle := func(item int) {
		// Data processing
		measurements <- item
	}

	// For equaling use divider.Fair divider, for prioritization use
	// divider.Rate divider or custom divider
	opts := simple.Opts[int]{
		Divider:          divider.Rate,
		Handle:           handle,
		HandlersQuantity: handlersQuantity,
		Inputs:           inputsOpts,
	}

	discipline, err := simple.New(opts)
	if err != nil {
		panic(err)
	}

	// Running writers, that write data to input channels
	for _, input := range inputs {
		go func(channel chan int) {
			defer close(channel)

			for id := range itemsQuantity {
				channel <- id
			}
		}(input)
	}

	// Waiting for the completion of the discipline
	go func() {
		defer close(measurements)

		for err := range discipline.Err() {
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
