// Internal package with the implementation of a discipline that does not
// manage data according to their priority.
package unmanaged

import (
	"sync"

	"github.com/akramarenkov/cqos/v2/internal/general"
	"github.com/akramarenkov/cqos/v2/priority/internal/common"
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

	err chan error
}

func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	capacity, err := general.CalcByFactor(
		int(opts.HandlersQuantity),
		common.DefaultCapacityFactor,
		len(opts.Inputs),
	)
	if err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		inputs: make(map[uint]common.Input[Type], len(opts.Inputs)),
		output: make(chan types.Prioritized[Type], capacity),

		err: make(chan error, 1),
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

func (dsc *Discipline[Type]) Err() <-chan error {
	return dsc.err
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
	defer close(dsc.err)
	defer close(dsc.output)

	dsc.loop()
}

func (dsc *Discipline[Type]) loop() {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for priority := range dsc.inputs {
		wg.Add(1)

		go dsc.io(wg, priority)
	}
}

func (dsc *Discipline[Type]) io(wg *sync.WaitGroup, priority uint) {
	defer wg.Done()

	for item := range dsc.inputs[priority].Channel {
		dsc.send(item, priority)
	}
}

func (dsc *Discipline[Type]) send(item Type, priority uint) {
	prioritized := types.Prioritized[Type]{
		Priority: priority,
		Item:     item,
	}

	dsc.output <- prioritized
}
