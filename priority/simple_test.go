package priority

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSimpleOptsValidation(t *testing.T) {
	opts := SimpleOpts[string]{
		Handle: func(context.Context, string) {},
		Inputs: map[uint]<-chan string{
			1: make(chan string),
		},
	}

	_, err := NewSimple(opts)
	require.Error(t, err)

	opts = SimpleOpts[string]{
		Divider: RateDivider,
		Inputs: map[uint]<-chan string{
			1: make(chan string),
		},
	}

	_, err = NewSimple(opts)
	require.Error(t, err)

	opts = SimpleOpts[string]{
		Divider: RateDivider,
		Handle:  func(context.Context, string) {},
	}

	_, err = NewSimple(opts)
	require.Error(t, err)
}

func TestSimple(t *testing.T) {
	testSimple(t, false)
}

func TestSimpleBadDivider(t *testing.T) {
	testSimple(t, true)
}

func testSimple(t *testing.T, useBadDivider bool) {
	handlersQuantity := 100
	inputCapacity := 10
	itemsQuantity := 100000

	inputs := map[uint]chan string{
		3: make(chan string, inputCapacity),
		2: make(chan string, inputCapacity),
		1: make(chan string, inputCapacity),
	}

	inputsOpts := make(map[uint]<-chan string)

	for priority, channel := range inputs {
		inputsOpts[priority] = channel
	}

	measures := make(chan string)
	defer close(measures)

	handle := func(ctx context.Context, item string) {
		select {
		case <-ctx.Done():
		case measures <- item:
		}
	}

	badDivider := func(
		priorities []uint,
		dividend uint,
		distribution map[uint]uint,
	) map[uint]uint {
		distribution = FairDivider(priorities, dividend, distribution)

		for priority := range distribution {
			distribution[priority] *= 2
		}

		return distribution
	}

	opts := SimpleOpts[string]{
		Divider:          FairDivider,
		Handle:           handle,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
	}

	if useBadDivider {
		opts.Divider = badDivider
	}

	simple, err := NewSimple(opts)
	require.NoError(t, err)

	defer simple.Stop()

	for priority, input := range inputs {
		go func(precedency uint, channel chan string) {
			defer close(channel)

			base := strconv.Itoa(int(precedency))

			for id := 0; id < itemsQuantity; id++ {
				item := base + ":" + strconv.Itoa(id)

				select {
				case <-simple.Err():
					return
				case channel <- item:
				}
			}
		}(priority, input)
	}

	received := 0

	// located in a function for easy use of the return
	func() {
		for {
			select {
			case <-simple.Err():
				return
			case <-measures:
				received++

				if received == itemsQuantity*len(inputs) {
					return
				}
			}
		}
	}()

	if useBadDivider {
		require.Equal(t, 0, received)
	} else {
		require.Equal(t, itemsQuantity*len(inputs), received)
	}
}

func TestSimpleStop(t *testing.T) {
	handlersQuantity := 100
	inputCapacity := 10
	itemsQuantity := 100000

	inputs := map[uint]chan string{
		3: make(chan string, inputCapacity),
		2: make(chan string, inputCapacity),
		1: make(chan string, inputCapacity),
	}

	inputsOpts := make(map[uint]<-chan string)

	for priority, channel := range inputs {
		inputsOpts[priority] = channel
	}

	measures := make(chan string)
	defer close(measures)

	handle := func(ctx context.Context, item string) {
		select {
		case <-ctx.Done():
		case measures <- item:
			time.Sleep(1 * time.Millisecond)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts := SimpleOpts[string]{
		Ctx:              ctx,
		Divider:          RateDivider,
		Handle:           handle,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
	}

	simple, err := NewSimple(opts)
	require.NoError(t, err)

	defer simple.Stop()
	defer simple.Stop()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for priority, input := range inputs {
		wg.Add(1)

		go func(precedency uint, channel chan string) {
			defer wg.Done()
			defer close(channel)

			base := strconv.Itoa(int(precedency))

			for id := 0; id < itemsQuantity; id++ {
				item := base + ":" + strconv.Itoa(id)

				select {
				case <-ctx.Done():
					return
				case channel <- item:
				}
			}
		}(priority, input)
	}

	received := 0

	func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-measures:
				received++

				if received == itemsQuantity*len(inputs) {
					return
				}
			}
		}
	}()

	require.NotEqual(t, 0, received)
	require.NotEqual(t, itemsQuantity*len(inputs), received)
}

func TestSimpleGracefulStop(t *testing.T) {
	handlersQuantity := 100
	inputCapacity := 10
	itemsQuantity := 100

	inputs := map[uint]chan string{
		3: make(chan string, inputCapacity),
		2: make(chan string, inputCapacity),
		1: make(chan string, inputCapacity),
	}

	inputsOpts := make(map[uint]<-chan string)

	for priority, channel := range inputs {
		inputsOpts[priority] = channel
	}

	measures := make(chan string)
	defer close(measures)

	handle := func(ctx context.Context, item string) {
		select {
		case <-ctx.Done():
		case measures <- item:
		}
	}

	opts := SimpleOpts[string]{
		Divider:          RateDivider,
		Handle:           handle,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
	}

	simple, err := NewSimple(opts)
	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

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

		for range measures {
			received++

			if received == itemsQuantity*len(inputs) {
				obtained <- received
				return
			}
		}
	}()

	simple.GracefulStop()

	require.Equal(t, itemsQuantity*len(inputs), <-obtained)
}
