package unmanaged

import (
	"sync"

	"github.com/akramarenkov/cqos/v2/priority/types"
)

type Opts[Type any] struct {
	Inputs map[uint]<-chan Type
}

type Unmanaged[Type any] struct {
	opts Opts[Type]

	output chan types.Prioritized[Type]
}

func New[Type any](opts Opts[Type]) (*Unmanaged[Type], error) {
	nmn := &Unmanaged[Type]{
		opts: opts,

		output: make(chan types.Prioritized[Type], 1),
	}

	go nmn.main()

	return nmn, nil
}

func (nmn *Unmanaged[Type]) Output() <-chan types.Prioritized[Type] {
	return nmn.output
}

func (nmn *Unmanaged[Type]) Release(uint) {
}

func (nmn *Unmanaged[Type]) main() {
	defer close(nmn.output)

	waiter := &sync.WaitGroup{}
	defer waiter.Wait()

	for priority := range nmn.opts.Inputs {
		waiter.Add(1)

		go nmn.io(waiter, priority)
	}
}

func (nmn *Unmanaged[Type]) io(waiter *sync.WaitGroup, priority uint) {
	defer waiter.Done()

	for item := range nmn.opts.Inputs[priority] {
		nmn.send(item, priority)
	}
}

func (nmn *Unmanaged[Type]) send(item Type, priority uint) {
	prioritized := types.Prioritized[Type]{
		Priority: priority,
		Item:     item,
	}

	nmn.output <- prioritized
}
