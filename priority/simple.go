package priority

import (
	"context"
	"sync"
)

type Pickup func(priorities []uint, divider Divider, quantity uint) uint
type Handle[Type any] func(item Type) uint

type SimpleOpts[Type any] struct {
	// Terminates (cancels) work of the discipline
	Ctx context.Context
	// Determines how handlers are distributed among priorities
	Divider Divider
	// Between how many handlers you need to distribute data
	HandlersQuantity uint
	// Channels with input data, should be buffered for performance reasons. Map key is a value of priority
	Inputs map[uint]<-chan Type
	Pickup Pickup
	Handle Handle[Type]
}

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

func NewSimple[Type any](opts SimpleOpts[Type]) (*Simple[Type], error) {
	output := make(chan Prioritized[Type])
	feedback := make(chan uint)

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

	smpl.runHandlers()

	return smpl, nil
}

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

func (smpl *Simple[Type]) runHandlers() {
	defer close(smpl.completer)
	defer smpl.wg.Wait()

	for id := 0; id < int(smpl.opts.HandlersQuantity); id++ {
		smpl.wg.Add(1)

		go smpl.handler()
	}
}

func (smpl *Simple[Type]) handler() {
	defer smpl.wg.Done()

	for {
		select {
		case <-smpl.breaker:
			return
		case <-smpl.opts.Ctx.Done():
			return
		case prioritized := <-smpl.output:
			smpl.opts.Handle(prioritized.Item)
		}
	}
}
