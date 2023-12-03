// Simplified version of the prioritization discipline that runs handlers on its own
package simple

import (
	"errors"
	"sync"

	"github.com/akramarenkov/cqos/v2/priority"
	"github.com/akramarenkov/cqos/v2/priority/divider"
)

var (
	ErrEmptyHandle = errors.New("handle function was not specified")
)

// Callback function called in handlers when an item is received
type Handle[Type any] func(item Type)

// Options of the created discipline
type Opts[Type any] struct {
	// Determines how handlers are distributed among priorities
	Divider divider.Divider
	// Callback function called in handlers when an item is received
	Handle Handle[Type]
	// Between how many handlers you need to distribute data
	HandlersQuantity uint
	// Channels with input data, should be buffered for performance reasons
	// Map key is a value of priority
	// For terminate discipline it is necessary and sufficient to close all input channels
	Inputs map[uint]<-chan Type
}

func (opts Opts[Type]) isValid() error {
	if opts.Handle == nil {
		return ErrEmptyHandle
	}

	return nil
}

// Simplified prioritization discipline.
//
// Preferably input channels should be buffered for performance reasons.
//
// For equaling use divider.Fair divider, for prioritization use divider.Rate divider or
// custom divider
type Discipline[Type any] struct {
	opts Opts[Type]

	priority *priority.Discipline[Type]
}

// Creates and runs discipline
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	priorityOpts := priority.Opts[Type]{
		Divider:          opts.Divider,
		HandlersQuantity: opts.HandlersQuantity,
		Inputs:           opts.Inputs,
	}

	priority, err := priority.New(priorityOpts)
	if err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		priority: priority,
	}

	go dsc.main()

	return dsc, nil
}

// Returns a channel with errors. If an error occurs (the value from the channel
// is not equal to nil) the discipline terminates its work. The most likely cause of
// the error is an incorrectly working dividing function in which the sum of
// the distributed quantities is not equal to the original quantity.
//
// The single nil value means that the discipline has terminated in normal mode
func (dsc *Discipline[Type]) Err() <-chan error {
	return dsc.priority.Err()
}

func (dsc *Discipline[Type]) main() {
	wg := &sync.WaitGroup{}

	for id := uint(0); id < dsc.opts.HandlersQuantity; id++ {
		wg.Add(1)

		go dsc.handler(wg)
	}

	wg.Wait()
}

func (dsc *Discipline[Type]) handler(wg *sync.WaitGroup) {
	defer wg.Done()

	for prioritized := range dsc.priority.Output() {
		dsc.opts.Handle(prioritized.Item)
		dsc.priority.Release(prioritized.Priority)
	}
}
