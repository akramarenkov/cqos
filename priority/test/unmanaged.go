package test

import (
	"errors"
	"sync"

	"github.com/akramarenkov/cqos/types"
)

var (
	ErrEmptyOutput = errors.New("output channel was not specified")
)

type UnmanagedOpts[Type any] struct {
	Inputs map[uint]<-chan Type
	Output chan<- types.Prioritized[Type]
}

type Unmanaged[Type any] struct {
	opts UnmanagedOpts[Type]

	breaker   chan bool
	completer chan bool
	stopMutex *sync.Mutex
	stopped   bool
}

func NewUnmanaged[Type any](opts UnmanagedOpts[Type]) (*Unmanaged[Type], error) {
	if opts.Output == nil {
		return nil, ErrEmptyOutput
	}

	dsc := &Unmanaged[Type]{
		opts: opts,

		breaker:   make(chan bool),
		completer: make(chan bool),
		stopMutex: &sync.Mutex{},
	}

	go dsc.main()

	return dsc, nil
}

func (dsc *Unmanaged[Type]) Stop() {
	dsc.stopMutex.Lock()
	defer dsc.stopMutex.Unlock()

	if dsc.stopped {
		return
	}

	dsc.stop()

	dsc.stopped = true
}

func (dsc *Unmanaged[Type]) stop() {
	close(dsc.breaker)
	<-dsc.completer
}

func (dsc *Unmanaged[Type]) main() {
	defer close(dsc.completer)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for priority, channel := range dsc.opts.Inputs {
		wg.Add(1)

		go dsc.io(wg, priority, channel)
	}
}

func (dsc *Unmanaged[Type]) io(wg *sync.WaitGroup, priority uint, channel <-chan Type) {
	defer wg.Done()

	for {
		select {
		case <-dsc.breaker:
			return
		case item, opened := <-dsc.opts.Inputs[priority]:
			if !opened {
				return
			}

			dsc.send(item, priority)
		}
	}
}

func (dsc *Unmanaged[Type]) send(item Type, priority uint) {
	prioritized := types.Prioritized[Type]{
		Priority: priority,
		Item:     item,
	}

	for {
		select {
		case <-dsc.breaker:
			return
		case dsc.opts.Output <- prioritized:
			return
		}
	}
}
