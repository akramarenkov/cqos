package priority

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	handlersQuantity := 100
	inputCapacity := 10
	itemsQuantity := 100000

	inputs := map[uint]chan string{
		3: make(chan string, inputCapacity),
		2: make(chan string, inputCapacity),
		1: make(chan string, inputCapacity),
	}

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

	measurements := make(chan bool)
	defer close(measurements)

	handle := func(ctx context.Context, item string) {
		measurements <- true
	}

	opts := SimpleOpts[string]{
		Divider:          RateDivider,
		Handle:           handle,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
	}

	simple, err := NewSimple(opts)
	require.NoError(t, err)

	defer simple.Stop()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

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

	for range measurements {
		received++

		if received == itemsQuantity*len(inputs) {
			break
		}
	}

	require.Equal(t, itemsQuantity*len(inputs), received)
}
