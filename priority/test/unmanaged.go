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

	nmn := &Unmanaged[Type]{
		opts: opts,

		breaker:   make(chan bool),
		completer: make(chan bool),
		stopMutex: &sync.Mutex{},
	}

	go nmn.main()

	return nmn, nil
}

func (nmn *Unmanaged[Type]) Stop() {
	nmn.stopMutex.Lock()
	defer nmn.stopMutex.Unlock()

	if nmn.stopped {
		return
	}

	nmn.stop()

	nmn.stopped = true
}

func (nmn *Unmanaged[Type]) stop() {
	close(nmn.breaker)
	<-nmn.completer
}

func (nmn *Unmanaged[Type]) main() {
	defer close(nmn.completer)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for priority, channel := range nmn.opts.Inputs {
		wg.Add(1)

		go nmn.io(wg, priority, channel)
	}
}

func (nmn *Unmanaged[Type]) io(wg *sync.WaitGroup, priority uint, channel <-chan Type) {
	defer wg.Done()

	for {
		select {
		case <-nmn.breaker:
			return
		case item, opened := <-nmn.opts.Inputs[priority]:
			if !opened {
				return
			}

			nmn.send(item, priority)
		}
	}
}

func (nmn *Unmanaged[Type]) send(item Type, priority uint) {
	prioritized := types.Prioritized[Type]{
		Priority: priority,
		Item:     item,
	}

	for {
		select {
		case <-nmn.breaker:
			return
		case nmn.opts.Output <- prioritized:
			return
		}
	}
}
