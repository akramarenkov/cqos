package priority

import (
	"context"
	"errors"
	"sync"

	"github.com/akramarenkov/cqos/v2/priority/divider"
)

var (
	ErrEmptyHandle = errors.New("handle function was not specified")
)

// Callback function called in handlers of simplified prioritization
// discipline when an item is received.
//
// Function should be interrupted when context is canceled
type Handle[Type any] func(ctx context.Context, item Type)

// Options of the created simplified prioritization discipline
type SimpleOpts[Type any] struct {
	// Determines how handlers are distributed among priorities
	Divider divider.Divider
	// Callback function called in handlers when an item is received
	Handle Handle[Type]
	// Between how many handlers you need to distribute data
	HandlersQuantity uint
	// Channels with input data, should be buffered for performance reasons
	// Map key is a value of priority
	// For graceful termination need close all input channels
	Inputs map[uint]<-chan Type
}

func (opts SimpleOpts[Type]) isValid() error {
	if opts.Handle == nil {
		return ErrEmptyHandle
	}

	return nil
}

// Simplified version of the discipline that runs handlers on its own and
// hides the output and feedback channels
type Simple[Type any] struct {
	opts SimpleOpts[Type]

	discipline *Discipline[Type]

	wg *sync.WaitGroup

	err chan error
}

// Creates and runs simplified prioritization discipline
func NewSimple[Type any](opts SimpleOpts[Type]) (*Simple[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	disciplineOpts := Opts[Type]{
		Divider:          opts.Divider,
		HandlersQuantity: opts.HandlersQuantity,
		Inputs:           opts.Inputs,
	}

	discipline, err := New(disciplineOpts)
	if err != nil {
		return nil, err
	}

	smpl := &Simple[Type]{
		opts: opts,

		discipline: discipline,

		wg: &sync.WaitGroup{},

		err: make(chan error, 1),
	}

	go smpl.handlers()

	return smpl, nil
}

// Returns a channel with errors. If an error occurs (the value from the channel
// is not equal to nil) the discipline terminates its work. The most likely cause of
// the error is an incorrectly working dividing function in which the sum of
// the distributed quantities is not equal to the original quantity.
//
// The single nil value means that the discipline has terminated in normal mode
func (smpl *Simple[Type]) Err() <-chan error {
	return smpl.err
}

func (smpl *Simple[Type]) handlers() {
	defer close(smpl.err)
	defer smpl.wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for id := uint(0); id < smpl.opts.HandlersQuantity; id++ {
		smpl.wg.Add(1)

		go smpl.handler(ctx)
	}

	err := <-smpl.discipline.Err()
	smpl.err <- err
}

func (smpl *Simple[Type]) handler(ctx context.Context) {
	defer smpl.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case prioritized, opened := <-smpl.discipline.Output():
			if !opened {
				return
			}

			smpl.opts.Handle(ctx, prioritized.Item)
			smpl.discipline.Release(prioritized.Priority)
		}
	}
}
