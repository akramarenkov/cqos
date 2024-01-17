package priority

import (
	"sync"

	"github.com/akramarenkov/breaker/breaker"
)

type unmanagedOpts[Type any] struct {
	Inputs map[uint]<-chan Type
	Output chan<- Prioritized[Type]
}

type unmanaged[Type any] struct {
	opts unmanagedOpts[Type]

	breaker *breaker.Breaker

	err chan error
}

func newUnmanaged[Type any](opts unmanagedOpts[Type]) (*unmanaged[Type], error) {
	if opts.Output == nil {
		return nil, ErrEmptyOutput
	}

	nmn := &unmanaged[Type]{
		opts: opts,

		breaker: breaker.New(),

		err: make(chan error, 1),
	}

	go nmn.main()

	return nmn, nil
}

func (nmn *unmanaged[Type]) Stop() {
	nmn.breaker.Break()
}

func (nmn *unmanaged[Type]) Err() <-chan error {
	return nmn.err
}

func (nmn *unmanaged[Type]) main() {
	defer close(nmn.err)
	defer nmn.breaker.Complete()

	nmn.loop()
}

func (nmn *unmanaged[Type]) loop() {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for priority := range nmn.opts.Inputs {
		wg.Add(1)

		go nmn.io(wg, priority)
	}
}

func (nmn *unmanaged[Type]) io(wg *sync.WaitGroup, priority uint) {
	defer wg.Done()

	for {
		select {
		case <-nmn.breaker.IsBreaked():
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
		case <-nmn.breaker.IsBreaked():
			return
		case nmn.opts.Output <- prioritized:
			return
		}
	}
}
