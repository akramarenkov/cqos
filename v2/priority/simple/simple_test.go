package simple

import (
	"strconv"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/divider"

	"github.com/stretchr/testify/require"
)

func TestOptsValidation(t *testing.T) {
	opts := Opts[string]{
		Handle: func(item string) {},
	}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[string]{
		Divider: divider.Rate,
	}

	_, err = New(opts)
	require.Error(t, err)
}

func testDiscipline(t *testing.T, useBadDivider bool) {
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

	measurements := make(chan bool)
	defer close(measurements)

	handle := func(item string) {
		time.Sleep(1 * time.Millisecond)

		measurements <- true
	}

	badDivider := func(
		priorities []uint,
		dividend uint,
		distribution map[uint]uint,
	) map[uint]uint {
		out := divider.Fair(priorities, dividend, distribution)

		for priority := range out {
			out[priority] *= 2
		}

		return out
	}

	opts := Opts[string]{
		Divider:          divider.Rate,
		Handle:           handle,
		HandlersQuantity: uint(handlersQuantity),
		Inputs:           inputsOpts,
	}

	if useBadDivider {
		opts.Divider = badDivider
	}

	simple, err := New(opts)
	require.NoError(t, err)

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
			case <-measurements:
				received++
			}
		}
	}()

	if useBadDivider {
		require.NotEqual(t, 0, received)
		require.NotEqual(t, itemsQuantity*len(inputs), received)
	} else {
		require.Equal(t, itemsQuantity*len(inputs), received)
	}
}

func TestDiscipline(t *testing.T) {
	testDiscipline(t, false)
}

func TestBadDivider(t *testing.T) {
	testDiscipline(t, true)
}