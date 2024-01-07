package priority

import (
	"context"
	"errors"
	"sync"

	"github.com/akramarenkov/cqos/breaker"
	"github.com/akramarenkov/cqos/internal/general"
	"github.com/akramarenkov/cqos/priority/internal/common"
)

var (
	ErrEmptyHandle = errors.New("handle function was not specified")
	ErrEmptyInput  = errors.New("input channels was not specified")
)

// Callback function called in handlers when an item is received.
//
// Function should be interrupted when context is canceled.
type Handle[Type any] func(ctx context.Context, item Type)

// Options of the created discipline.
type SimpleOpts[Type any] struct {
	// Roughly terminates (cancels) work of the discipline
	Ctx context.Context
	// Determines how handlers are distributed among priorities
	Divider Divider
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

	if len(opts.Inputs) == 0 {
		return ErrEmptyInput
	}

	return nil
}

func (opts SimpleOpts[Type]) normalize() SimpleOpts[Type] {
	if opts.Ctx == nil {
		opts.Ctx = context.Background()
	}

	return opts
}

// Simplified version of the discipline that runs handlers on its own and
// hides the output and feedback channels.
//
// Preferably input channels should be buffered for performance reasons.
//
// For equaling use divider.Fair divider, for prioritization use divider.Rate divider or
// custom divider.
type Simple[Type any] struct {
	opts SimpleOpts[Type]

	priority *Discipline[Type]

	breaker  *breaker.Breaker
	graceful *breaker.Breaker

	output   chan Prioritized[Type]
	feedback chan uint

	wg *sync.WaitGroup

	err chan error
}

// Creates and runs discipline.
func NewSimple[Type any](opts SimpleOpts[Type]) (*Simple[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	opts = opts.normalize()

	capacity, err := general.CalcByFactor(
		int(opts.HandlersQuantity),
		common.DefaultCapacityFactor,
		len(opts.Inputs),
	)
	if err != nil {
		return nil, err
	}

	output := make(chan Prioritized[Type], capacity)
	feedback := make(chan uint, capacity)

	priorityOpts := Opts[Type]{
		Divider:          opts.Divider,
		Feedback:         feedback,
		HandlersQuantity: opts.HandlersQuantity,
		Inputs:           opts.Inputs,
		Output:           output,
	}

	priority, err := New(priorityOpts)
	if err != nil {
		return nil, err
	}

	smpl := &Simple[Type]{
		opts: opts,

		priority: priority,

		breaker:  breaker.New(),
		graceful: breaker.New(),

		output:   output,
		feedback: feedback,

		wg: &sync.WaitGroup{},

		err: make(chan error, 1),
	}

	go smpl.main()

	return smpl, nil
}

// Returns a channel with errors. If an error occurs (the value from the channel
// is not equal to nil) the discipline terminates its work. The most likely cause of
// the error is an incorrectly working dividing function in which the sum of
// the distributed quantities is not equal to the original quantity.
//
// The single nil value means that the discipline has terminated in normal mode.
func (smpl *Simple[Type]) Err() <-chan error {
	return smpl.err
}

// Roughly terminates work of the discipline.
//
// Use for wait completion at terminates via context.
func (smpl *Simple[Type]) Stop() {
	smpl.breaker.Break()
}

// Graceful terminates work of the discipline.
//
// Waits draining input channels, waits end processing data in handlers and terminates.
//
// You must end write to input channels and close them, otherwise graceful stop not be
// ended.
func (smpl *Simple[Type]) GracefulStop() {
	smpl.graceful.Break()
}

func (smpl *Simple[Type]) main() {
	defer smpl.breaker.Complete()
	defer smpl.graceful.Complete()
	defer close(smpl.err)
	defer close(smpl.output)
	defer close(smpl.feedback)
	defer smpl.wg.Wait()

	ctx, cancel := context.WithCancel(smpl.opts.Ctx)
	defer cancel()

	defer smpl.priority.Stop()

	for id := uint(0); id < smpl.opts.HandlersQuantity; id++ {
		smpl.wg.Add(1)

		go smpl.handler(ctx)
	}

	select {
	case <-smpl.breaker.Breaked():
	case <-smpl.opts.Ctx.Done():
	case <-smpl.graceful.Breaked():
		smpl.priority.GracefulStop()
	case err := <-smpl.priority.Err():
		smpl.err <- err
	}
}

func (smpl *Simple[Type]) handler(ctx context.Context) {
	defer smpl.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case prioritized := <-smpl.output:
			smpl.opts.Handle(ctx, prioritized.Item)

			select {
			case <-ctx.Done():
				return
			case smpl.feedback <- prioritized.Priority:
			}
		}
	}
}
