package priority

import (
	"fmt"
	"strconv"
	"sync"
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

	defer func() {
		for _, input := range inputs {
			close(input)
		}
	}()

	// Data from input channels passed to handlers by output channel
	output := make(chan Prioritized[string])

	// Handlers must write priority of processed data to feedback channel after it has been processed
	feedback := make(chan uint)
	defer close(feedback)

	// Used only in this example for detect that all written data are processed
	measurements := make(chan bool)
	defer close(measurements)

	// For equaling use FairDivider, for prioritization use RateDivider or custom divider
	disciplineOpts := Opts[string]{
		Divider:          RateDivider,
		Feedback:         feedback,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
		Output:           output,
	}

	discipline, err := New(disciplineOpts)
	if err != nil {
		panic(err)
	}

	defer discipline.Stop()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	// Run handlers, that process data
	for handler := 0; handler < handlersQuantity; handler++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for prioritized := range output {
				// Data processing
				// fmt.Println(prioritized.Item)
				measurements <- true

				feedback <- prioritized.Priority
			}
		}()
	}

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

	// Terminate handlers
	defer close(output)

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
	output := make(chan Prioritized[string])

	// Handlers must write priority of processed data to feedback channel after it has been processed
	feedback := make(chan uint)
	defer close(feedback)

	// Used only in this example for detect that all written data are processed
	measurements := make(chan bool)

	// For equaling use FairDivider, for prioritization use RateDivider or custom divider
	disciplineOpts := Opts[string]{
		Divider:          RateDivider,
		Feedback:         feedback,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
		Output:           output,
	}

	discipline, err := New(disciplineOpts)
	if err != nil {
		panic(err)
	}

	wgh := &sync.WaitGroup{}
	defer wgh.Wait()

	// Run handlers, that process data
	for handler := 0; handler < handlersQuantity; handler++ {
		wgh.Add(1)

		go func() {
			defer wgh.Done()

			for prioritized := range output {
				// Data processing
				// fmt.Println(prioritized.Item)
				measurements <- true

				feedback <- prioritized.Priority
			}
		}()
	}

	wgw := &sync.WaitGroup{}

	// Run writers, that write data to input channels
	for priority, input := range inputs {
		wgw.Add(1)

		go func(precedency uint, channel chan string) {
			defer wgw.Done()

			base := strconv.Itoa(int(precedency))

			for id := 0; id < itemsQuantity; id++ {
				item := base + ":" + strconv.Itoa(id)

				channel <- item
			}
		}(priority, input)
	}

	// Terminate handlers
	defer close(output)

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
	wgw.Wait()

	for _, input := range inputs {
		close(input)
	}

	discipline.GracefulStop()

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
