package priority_test

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/akramarenkov/cqos/priority"
)

func ExampleDiscipline() {
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

	// Data from input channels passed to handlers by output channel
	output := make(chan priority.Prioritized[string])

	// Handlers must write priority of processed data to feedback channel after it has been processed
	feedback := make(chan uint)
	defer close(feedback)

	// Used only in this example for detect that all written data are processed
	measures := make(chan string)
	defer close(measures)

	// For equaling use FairDivider, for prioritization use
	// RateDivider or custom divider
	opts := priority.Opts[string]{
		Divider:          priority.RateDivider,
		Feedback:         feedback,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
		Output:           output,
	}

	discipline, err := priority.New(opts)
	if err != nil {
		panic(err)
	}

	defer discipline.Stop()

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

	// Run handlers, that process data
	for handler := 0; handler < handlersQuantity; handler++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for prioritized := range output {
				// Data processing
				measures <- prioritized.Item

				// Handler must indicate that current data has been processed and
				// handler is ready to receive new data
				feedback <- prioritized.Priority
			}
		}()
	}

	// Terminate handlers
	defer close(output)

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

func ExampleDiscipline_GracefulStop() {
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

	// Data from input channels passed to handlers by output channel
	output := make(chan priority.Prioritized[string])

	// Handlers must write priority of processed data to feedback channel after it has been processed
	feedback := make(chan uint)
	defer close(feedback)

	// Used only in this example for detect that all written data are processed
	measures := make(chan string)
	defer close(measures)

	// For equaling use FairDivider, for prioritization use
	// RateDivider or custom divider
	opts := priority.Opts[string]{
		Divider:          priority.RateDivider,
		Feedback:         feedback,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
		Output:           output,
	}

	discipline, err := priority.New(opts)
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

	// Run handlers, that process data
	for handler := 0; handler < handlersQuantity; handler++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for prioritized := range output {
				// Data processing
				measures <- prioritized.Item

				feedback <- prioritized.Priority
			}
		}()
	}

	// Terminate handlers
	defer close(output)

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

	discipline.GracefulStop()

	fmt.Println("Processed items quantity:", <-obtained)
	// Output: Processed items quantity: 300
}
