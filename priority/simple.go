package priority

import (
	"context"
	"sync"
)

const (
	defaultCapacityFactor = 0.1
)

// Callback function called in handlers of simplified prioritization discipline when an item is received.
//
// Function should be interrupted when context is canceled
type Handle[Type any] func(ctx context.Context, item Type)

// Options of the created simplified prioritization discipline
type SimpleOpts[Type any] struct {
	// Terminates (cancels) work of the discipline
	Ctx context.Context
	// Determines how handlers are distributed among priorities
	Divider Divider
	// Callback function called in handlers when an item is received
	Handle Handle[Type]
	// Between how many handlers you need to distribute data
	HandlersQuantity uint
	// Channels with input data, should be buffered for performance reasons. Map key is a value of priority
	Inputs map[uint]<-chan Type
}

func (opts SimpleOpts[Type]) normalize() SimpleOpts[Type] {
	if opts.Ctx == nil {
		opts.Ctx = context.Background()
	}

	return opts
}

// Simplified version that runs handlers on its own and hides the output and
// feedback channels
type Simple[Type any] struct {
	opts SimpleOpts[Type]

	discipline *Discipline[Type]

	breaked      bool
	breaker      chan bool
	breakerMutex *sync.Mutex
	completer    chan bool

	output   chan Prioritized[Type]
	feedback chan uint

	wg *sync.WaitGroup
}

// Creates and runs simplified prioritization discipline
func NewSimple[Type any](opts SimpleOpts[Type]) (*Simple[Type], error) {
	opts = opts.normalize()

	capacity := calcCapacity(int(opts.HandlersQuantity), defaultCapacityFactor, 1)

	output := make(chan Prioritized[Type], capacity)
	feedback := make(chan uint, capacity)

	disciplineOpts := Opts[Type]{
		Divider:          opts.Divider,
		Feedback:         feedback,
		HandlersQuantity: opts.HandlersQuantity,
		Inputs:           opts.Inputs,
		Output:           output,
	}

	discipline, err := New(disciplineOpts)
	if err != nil {
		return nil, err
	}

	smpl := &Simple[Type]{
		opts: opts,

		discipline: discipline,

		breaker:      make(chan bool),
		breakerMutex: &sync.Mutex{},
		completer:    make(chan bool),

		output:   output,
		feedback: feedback,

		wg: &sync.WaitGroup{},
	}

	go smpl.handlers()

	return smpl, nil
}

// Terminates work of the discipline.
//
// Use for wait completion at terminates via context
func (smpl *Simple[Type]) Stop() {
	smpl.stop()
	<-smpl.completer
}

func (smpl *Simple[Type]) stop() {
	smpl.breakerMutex.Lock()
	defer smpl.breakerMutex.Unlock()

	if smpl.breaked {
		return
	}

	close(smpl.breaker)

	smpl.breaked = true
}

func (smpl *Simple[Type]) handlers() {
	defer close(smpl.completer)
	defer close(smpl.output)
	defer close(smpl.feedback)
	defer smpl.discipline.Stop()
	defer smpl.wg.Wait()

	ctx, cancel := context.WithCancel(smpl.opts.Ctx)
	defer cancel()

	for id := 0; id < int(smpl.opts.HandlersQuantity); id++ {
		smpl.wg.Add(1)

		go smpl.handler(ctx)
	}

	select {
	case <-smpl.breaker:
		return
	case <-smpl.opts.Ctx.Done():
		return
	}
}

func (smpl *Simple[Type]) handler(ctx context.Context) {
	defer smpl.wg.Done()

	for {
		select {
		case <-smpl.breaker:
			return
		case <-smpl.opts.Ctx.Done():
			return
		case prioritized := <-smpl.output:
			smpl.opts.Handle(ctx, prioritized.Item)
			smpl.feedback <- prioritized.Priority
		}
	}
}
