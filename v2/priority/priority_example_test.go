package priority_test

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/akramarenkov/cqos/v2/priority"
	"github.com/akramarenkov/cqos/v2/priority/divider"
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

	// Used only in this example for detect that all written data are processed
	measures := make(chan bool)
	defer close(measures)

	// For equaling use divider.Fair divider, for prioritization use
	// divider.Rate divider or custom divider
	opts := priority.Opts[string]{
		Divider:          divider.Rate,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
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

			for prioritized := range discipline.Output() {
				// Data processing
				// fmt.Println(prioritized.Item)
				measures <- true

				// Handler must indicate that current data has been processed and
				// handler is ready to receive new data
				discipline.Release(prioritized.Priority)
			}
		}()
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
