// An internal package with the implementation of a discipline that does not
// manage data according to their priority
package unmanaged

import (
	"sync"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/internal/common"
	"github.com/akramarenkov/cqos/v2/priority/internal/consts"
	"github.com/akramarenkov/cqos/v2/priority/types"
)

type Opts[Type any] struct {
	HandlersQuantity uint
	Inputs           map[uint]<-chan Type
}

type Discipline[Type any] struct {
	opts Opts[Type]

	inputs map[uint]common.Input[Type]
	output chan types.Prioritized[Type]

	interrupter *time.Ticker
}

func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	capacity := common.CalcCapacity(
		int(opts.HandlersQuantity),
		consts.DefaultCapacityFactor,
		len(opts.Inputs),
	)

	dsc := &Discipline[Type]{
		opts: opts,

		inputs: make(map[uint]common.Input[Type], len(opts.Inputs)),
		output: make(chan types.Prioritized[Type], capacity),

		interrupter: time.NewTicker(consts.DefaultInterruptTimeout),
	}

	dsc.updateInputs(opts.Inputs)

	go dsc.main()

	return dsc, nil
}

func (dsc *Discipline[Type]) Output() <-chan types.Prioritized[Type] {
	return dsc.output
}

func (dsc *Discipline[Type]) Release(uint) {
}

func (dsc *Discipline[Type]) updateInputs(inputs map[uint]<-chan Type) {
	for priority, channel := range inputs {
		input := common.Input[Type]{
			Channel: channel,
		}

		dsc.inputs[priority] = input
	}
}

func (dsc *Discipline[Type]) main() {
	defer close(dsc.output)
	defer dsc.interrupter.Stop()

	waiter := &sync.WaitGroup{}
	defer waiter.Wait()

	for priority := range dsc.opts.Inputs {
		waiter.Add(1)

		if cap(dsc.inputs[priority].Channel) != 0 {
			go dsc.io(waiter, priority)
		} else {
			go dsc.iou(waiter, priority)
		}
	}
}

func (dsc *Discipline[Type]) io(waiter *sync.WaitGroup, priority uint) {
	defer waiter.Done()

	for item := range dsc.opts.Inputs[priority] {
		dsc.send(item, priority)
	}
}

func (dsc *Discipline[Type]) iou(waiter *sync.WaitGroup, priority uint) {
	defer waiter.Done()

	sleep := false

	for {
		select {
		case item, opened := <-dsc.inputs[priority].Channel:
			if !opened {
				return
			}

			dsc.send(item, priority)
		case <-dsc.interrupter.C:
			if sleep {
				sleep = false

				time.Sleep(consts.DefaultIdleDelay)

				continue
			}

			sleep = true
		}
	}
}

func (dsc *Discipline[Type]) send(item Type, priority uint) {
	prioritized := types.Prioritized[Type]{
		Priority: priority,
		Item:     item,
	}

	dsc.output <- prioritized
}
