package simple

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/divider"

	"github.com/stretchr/testify/require"
)

func TestOptsValidation(t *testing.T) {
	opts := Opts[string]{
		Handle: func(ctx context.Context, item string) {},
	}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[string]{
		Divider: divider.Rate,
	}

	_, err = New(opts)
	require.Error(t, err)
}

func TestDiscipline(t *testing.T) {
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

	measurements := make(chan bool)
	defer close(measurements)

	handle := func(ctx context.Context, item string) {
		select {
		case <-ctx.Done():
		case measurements <- true:
		}
	}

	opts := Opts[string]{
		Divider:          divider.Rate,
		Handle:           handle,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
	}

	_, err := New(opts)
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

	received := 0

	for range measurements {
		received++

		if received == itemsQuantity*len(inputs) {
			break
		}
	}

	require.Equal(t, itemsQuantity*len(inputs), received)
}

func TestBadDivider(t *testing.T) {
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

	measurements := make(chan bool)
	defer close(measurements)

	handle := func(ctx context.Context, item string) {
		select {
		case <-ctx.Done():
			return
		case measurements <- true:
			time.Sleep(1 * time.Millisecond)
		}
	}

	divider := func(priorities []uint, dividend uint, distribution map[uint]uint) map[uint]uint {
		out := divider.Fair(priorities, dividend, distribution)

		for priority := range out {
			out[priority] *= 2
		}

		return out
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := Opts[string]{
		Divider:          divider,
		Handle:           handle,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
	}

	simple, err := New(opts)
	require.NoError(t, err)

	go func() {
		if <-simple.Err() != nil {
			cancel()
		}
	}()

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

	defer func() {
		require.NotEqual(t, 0, received)
		require.NotEqual(t, itemsQuantity*len(inputs), received)
	}()

	defer func() {
		<-simple.Err()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-measurements:
			received++

			if received == itemsQuantity*len(inputs) {
				return
			}
		}
	}
}
