package priority

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

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
		select {
		case <-ctx.Done():
		case measurements <- true:
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

func TestSimpleStop(t *testing.T) {
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
		select {
		case <-ctx.Done():
		case measurements <- true:
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

func TestSimpleOptsValidation(t *testing.T) {
	opts := SimpleOpts[string]{
		Handle: func(ctx context.Context, item string) {},
	}

	_, err := NewSimple(opts)
	require.Error(t, err)

	opts = SimpleOpts[string]{
		Divider: RateDivider,
	}

	_, err = NewSimple(opts)
	require.Error(t, err)
}

func TestSimpleBadDivider(t *testing.T) {
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
		select {
		case <-ctx.Done():
		case measurements <- true:
			time.Sleep(1 * time.Millisecond)
		}
	}

	divider := func(priorities []uint, dividend uint, distribution map[uint]uint) map[uint]uint {
		out := FairDivider(priorities, dividend, distribution)

		for priority, quantity := range out {
			out[priority] = 2 * quantity
		}

		return out
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := SimpleOpts[string]{
		Divider:          divider,
		Handle:           handle,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
	}

	simple, err := NewSimple(opts)
	require.NoError(t, err)

	defer simple.Stop()

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

func TestSimpleGracefulStop(t *testing.T) {
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
		select {
		case <-ctx.Done():
		case measurements <- true:
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

			base := strconv.Itoa(int(precedency))

			for id := 0; id < itemsQuantity; id++ {
				item := base + ":" + strconv.Itoa(id)

				channel <- item
			}
		}(priority, input)
	}

	go func() {
		simple.GracefulStop()
	}()

	received := 0

	for range measurements {
		received++

		if received == itemsQuantity*len(inputs) {
			break
		}
	}

	require.Equal(t, itemsQuantity*len(inputs), received)
}
