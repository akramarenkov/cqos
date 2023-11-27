package priority

import (
	"sync"
)

type unmanagedOpts[Type any] struct {
	Inputs map[uint]<-chan Type
}

type unmanaged[Type any] struct {
	opts unmanagedOpts[Type]

	output chan Prioritized[Type]
}

func newUnmanaged[Type any](opts unmanagedOpts[Type]) (*unmanaged[Type], error) {
	nmn := &unmanaged[Type]{
		opts: opts,

		output: make(chan Prioritized[Type], 1),
	}

	go nmn.main()

	return nmn, nil
}

func (nmn *unmanaged[Type]) Output() <-chan Prioritized[Type] {
	return nmn.output
}

func (nmn *unmanaged[Type]) Release(uint) {
}

func (nmn *unmanaged[Type]) main() {
	defer close(nmn.output)

	waiter := &sync.WaitGroup{}
	defer waiter.Wait()

	for priority := range nmn.opts.Inputs {
		waiter.Add(1)

		go nmn.io(waiter, priority)
	}
}

func (nmn *unmanaged[Type]) io(waiter *sync.WaitGroup, priority uint) {
	defer waiter.Done()

	for item := range nmn.opts.Inputs[priority] {
		nmn.send(item, priority)
	}
}

func (nmn *unmanaged[Type]) send(item Type, priority uint) {
	prioritized := Prioritized[Type]{
		Priority: priority,
		Item:     item,
	}

	nmn.output <- prioritized
}
