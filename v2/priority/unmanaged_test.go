package priority

import (
	"sync"
)

type unmanagedOpts[Type any] struct {
	Inputs map[uint]<-chan Type
	Output chan<- Prioritized[Type]
}

type unmanaged[Type any] struct {
	opts unmanagedOpts[Type]

	breaker   chan bool
	completer chan bool
	stopMutex *sync.Mutex
	stopped   bool
}

func newUnmanaged[Type any](opts unmanagedOpts[Type]) (*unmanaged[Type], error) {
	if opts.Output == nil {
		return nil, ErrEmptyOutput
	}

	nmn := &unmanaged[Type]{
		opts: opts,

		breaker:   make(chan bool),
		completer: make(chan bool),
		stopMutex: &sync.Mutex{},
	}

	go nmn.main()

	return nmn, nil
}

func (nmn *unmanaged[Type]) Stop() {
	nmn.stopMutex.Lock()
	defer nmn.stopMutex.Unlock()

	if nmn.stopped {
		return
	}

	nmn.stop()

	nmn.stopped = true
}

func (nmn *unmanaged[Type]) stop() {
	close(nmn.breaker)
	<-nmn.completer
}

func (nmn *unmanaged[Type]) main() {
	defer close(nmn.completer)

	waiter := &sync.WaitGroup{}
	defer waiter.Wait()

	for priority := range nmn.opts.Inputs {
		waiter.Add(1)

		go nmn.io(waiter, priority)
	}
}

func (nmn *unmanaged[Type]) io(waiter *sync.WaitGroup, priority uint) {
	defer waiter.Done()

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

func (nmn *unmanaged[Type]) send(item Type, priority uint) {
	prioritized := Prioritized[Type]{
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
