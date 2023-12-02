// Discipline that used to distributes data among handlers according to priority
package priority

import (
	"errors"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/divider"
	"github.com/akramarenkov/cqos/v2/priority/internal/common"
	"github.com/akramarenkov/cqos/v2/priority/internal/consts"
	"github.com/akramarenkov/cqos/v2/priority/types"
)

var (
	ErrEmptyDivider             = errors.New("priorities divider was not specified")
	ErrHandlersQuantityExceeded = errors.New("value of handlers quantity has been exceeded")
)

const (
	defaultGetFeedbackFactor = 0.1
	defaultIdleDelay         = 1 * time.Nanosecond
	defaultInterruptTimeout  = 1 * time.Nanosecond
)

// Options of the created discipline
type Opts[Type any] struct {
	// Determines how handlers are distributed among priorities
	Divider divider.Divider
	// Between how many handlers you need to distribute data
	HandlersQuantity uint
	// Channels with input data, should be buffered for performance reasons
	// Map key is a value of priority
	// For terminate discipline it is necessary and sufficient to close all input channels
	Inputs map[uint]<-chan Type
}

// Prioritization discipline.
//
// Preferably input channels should be buffered for performance reasons.
//
// Data from input channels passed to handlers by output channel.
//
// Handlers must call Release() method after the current data item has been processed.
//
// For equaling use divider.Fair divider, for prioritization use divider.Rate divider or
// custom divider
type Discipline[Type any] struct {
	opts Opts[Type]

	feedback chan uint
	inputs   map[uint]common.Input[Type]
	output   chan types.Prioritized[Type]

	priorities []uint

	actual    map[uint]uint
	strategic map[uint]uint
	tactic    map[uint]uint

	uncrowded []uint
	useful    []uint

	maxGetFeedback int

	interrupter *time.Ticker

	err chan error
}

func (opts Opts[Type]) isValid() error {
	if opts.Divider == nil {
		return ErrEmptyDivider
	}

	return nil
}

// Creates and runs discipline
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	capacity := common.CalcByFactor(
		int(opts.HandlersQuantity),
		consts.DefaultCapacityFactor,
		len(opts.Inputs),
	)

	maxGetFeedback := common.CalcByFactor(
		int(opts.HandlersQuantity),
		defaultGetFeedbackFactor,
		len(opts.Inputs),
	)

	inputs, priorities, strategic, err := prepare(opts)
	if err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		feedback: make(chan uint, capacity),
		inputs:   inputs,
		output:   make(chan types.Prioritized[Type], capacity),

		priorities: priorities,

		actual:    make(map[uint]uint),
		strategic: strategic,
		tactic:    make(map[uint]uint),

		maxGetFeedback: maxGetFeedback,

		interrupter: time.NewTicker(defaultInterruptTimeout),

		err: make(chan error, 1),
	}

	go dsc.loop()

	return dsc, nil
}

func prepare[Type any](opts Opts[Type]) (
	map[uint]common.Input[Type],
	[]uint,
	map[uint]uint,
	error,
) {
	inputs := make(map[uint]common.Input[Type])
	priorities := make([]uint, 0)
	strategic := make(map[uint]uint)

	for priority, channel := range opts.Inputs {
		input := common.Input[Type]{
			Channel: channel,
		}

		inputs[priority] = input

		priorities = append(priorities, priority)
	}

	common.SortPriorities(priorities)

	err := safeDivide(
		opts.Divider,
		priorities,
		opts.HandlersQuantity,
		strategic,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	return inputs, priorities, strategic, nil
}

// Returns output channel.
//
// If this channel is closed, it means that the discipline is terminated
func (dsc *Discipline[Type]) Output() <-chan types.Prioritized[Type] {
	return dsc.output
}

// Marks that current data has been processed and handler is ready to receive new data
func (dsc *Discipline[Type]) Release(priority uint) {
	dsc.feedback <- priority
}

// Returns a channel with errors. If an error occurs (the value from the channel
// is not equal to nil) the discipline terminates its work. The most likely cause of
// the error is an incorrectly working dividing function in which the sum of
// the distributed quantities is not equal to the original quantity.
//
// The single nil value means that the discipline has terminated in normal mode
func (dsc *Discipline[Type]) Err() <-chan error {
	return dsc.err
}

func (dsc *Discipline[Type]) loop() {
	defer close(dsc.err)
	defer close(dsc.output)
	defer close(dsc.feedback)
	defer dsc.interrupter.Stop()

	for {
		dsc.getFeedback()

		processed, err := dsc.main()
		if err != nil {
			for !dsc.isZeroActual() {
				dsc.decreaseActual(<-dsc.feedback)
			}

			dsc.err <- err

			return
		}

		if processed == 0 {
			if dsc.isDrainedInputs() {
				for !dsc.isZeroActual() {
					dsc.decreaseActual(<-dsc.feedback)
				}

				return
			}

			time.Sleep(defaultIdleDelay)
		}
	}
}

func (dsc *Discipline[Type]) getFeedback() {
	for collected := 0; collected < dsc.maxGetFeedback; collected++ {
		select {
		case priority := <-dsc.feedback:
			dsc.decreaseActual(priority)
		default:
			return
		}
	}
}

func (dsc *Discipline[Type]) isZeroActual() bool {
	for _, quantity := range dsc.actual {
		if quantity != 0 {
			return false
		}
	}

	return true
}

func (dsc *Discipline[Type]) isDrainedInputs() bool {
	for _, input := range dsc.inputs {
		if !input.Drained {
			return false
		}
	}

	return true
}

func (dsc *Discipline[Type]) main() (uint, error) {
	processed := uint(0)

	proceed, err := dsc.calcTactic()
	if err != nil {
		return processed, err
	}

	if !proceed {
		return processed, nil
	}

	processed += dsc.prioritize()

	proceed, err = dsc.recalcTactic()
	if err != nil {
		return processed, err
	}

	if !proceed {
		return processed, nil
	}

	processed += dsc.prioritize()

	return processed, nil
}

func (dsc *Discipline[Type]) prioritize() uint {
	processed := uint(0)

	for _, priority := range dsc.priorities {
		if cap(dsc.inputs[priority].Channel) != 0 {
			processed += dsc.io(priority)
		} else {
			processed += dsc.iou(priority)
		}
	}

	return processed
}

func (dsc *Discipline[Type]) io(priority uint) uint {
	processed := uint(0)

	for {
		if dsc.tactic[priority] == 0 {
			return processed
		}

		select {
		case item, opened := <-dsc.inputs[priority].Channel:
			if !opened {
				dsc.markInputAsDrained(priority)
				return processed
			}

			processed += dsc.send(item, priority)
		default:
			return processed
		}
	}
}

func (dsc *Discipline[Type]) iou(priority uint) uint {
	processed := uint(0)

	interrupt := false

	for {
		if dsc.tactic[priority] == 0 {
			return processed
		}

		select {
		case item, opened := <-dsc.inputs[priority].Channel:
			if !opened {
				dsc.markInputAsDrained(priority)
				return processed
			}

			interrupt = false

			processed += dsc.send(item, priority)
		case <-dsc.interrupter.C:
			if interrupt {
				return processed
			}

			interrupt = true
		}
	}
}

func (dsc *Discipline[Type]) markInputAsDrained(priority uint) {
	input := dsc.inputs[priority]
	input.Drained = true
	dsc.inputs[priority] = input
}

func (dsc *Discipline[Type]) send(item Type, priority uint) uint {
	prioritized := types.Prioritized[Type]{
		Priority: priority,
		Item:     item,
	}

	dsc.output <- prioritized

	dsc.decreaseTactic(priority)
	dsc.increaseActual(priority)

	return 1
}

func (dsc *Discipline[Type]) increaseActual(priority uint) {
	dsc.actual[priority]++
}

func (dsc *Discipline[Type]) decreaseActual(priority uint) {
	dsc.actual[priority]--
}

func (dsc *Discipline[Type]) decreaseTactic(priority uint) {
	dsc.tactic[priority]--
}

func (dsc *Discipline[Type]) calcTactic() (bool, error) {
	vacants, err := dsc.calcVacants()
	if err != nil {
		return false, err
	}

	if vacants == 0 {
		return false, nil
	}

	if picked := dsc.calcTacticSimpleAddition(vacants); picked {
		return true, nil
	}

	return dsc.calcTacticBase(vacants)
}

func (dsc *Discipline[Type]) calcVacants() (uint, error) {
	busy := calcDistributionQuantity(dsc.actual)

	if dsc.opts.HandlersQuantity < busy {
		return 0, ErrHandlersQuantityExceeded
	}

	return dsc.opts.HandlersQuantity - busy, nil
}

func (dsc *Discipline[Type]) calcTacticSimpleAddition(vacants uint) bool {
	dsc.resetTactic()

	picked := uint(0)

	for _, priority := range dsc.priorities {
		if dsc.actual[priority] > dsc.strategic[priority] {
			return false
		}

		dsc.tactic[priority] = dsc.strategic[priority] - dsc.actual[priority]

		picked += dsc.tactic[priority]
	}

	return picked != 0 && picked <= vacants
}

func (dsc *Discipline[Type]) calcTacticBase(vacants uint) (bool, error) {
	dsc.resetTactic()
	dsc.updateUncrowded()

	err := safeDivide(
		dsc.opts.Divider,
		dsc.uncrowded,
		vacants,
		dsc.tactic,
	)
	if err != nil {
		return false, err
	}

	return dsc.isTacticFilled(dsc.uncrowded), nil
}

func (dsc *Discipline[Type]) updateUncrowded() {
	dsc.uncrowded = dsc.uncrowded[:0]

	for _, priority := range dsc.priorities {
		if dsc.actual[priority] >= dsc.strategic[priority] {
			continue
		}

		dsc.uncrowded = append(dsc.uncrowded, priority)
	}
}

func (dsc *Discipline[Type]) resetTactic() {
	for priority := range dsc.tactic {
		dsc.tactic[priority] = 0
	}
}

func (dsc *Discipline[Type]) isTacticFilled(priorities []uint) bool {
	for _, priority := range priorities {
		if dsc.tactic[priority] == 0 {
			return false
		}
	}

	return true
}

func (dsc *Discipline[Type]) recalcTactic() (bool, error) {
	remainder := dsc.calcVacantsRemainder()

	dsc.updateUseful()
	dsc.resetTactic()

	err := safeDivide(
		dsc.opts.Divider,
		dsc.useful,
		dsc.opts.HandlersQuantity,
		dsc.tactic,
	)
	if err != nil {
		return false, err
	}

	dsc.updateUsefulLikeUncrowded()
	dsc.resetTactic()

	err = safeDivide(
		dsc.opts.Divider,
		dsc.useful,
		remainder,
		dsc.tactic,
	)
	if err != nil {
		return false, err
	}

	return dsc.isTacticFilled(dsc.useful), nil
}

func (dsc *Discipline[Type]) calcVacantsRemainder() uint {
	remainder := uint(0)

	for _, priority := range dsc.priorities {
		if dsc.tactic[priority] != 0 {
			remainder += dsc.tactic[priority]
			continue
		}
	}

	return remainder
}

func (dsc *Discipline[Type]) updateUseful() {
	dsc.useful = dsc.useful[:0]

	for _, priority := range dsc.priorities {
		if dsc.tactic[priority] != 0 {
			continue
		}

		dsc.useful = append(dsc.useful, priority)
	}
}

func (dsc *Discipline[Type]) updateUsefulLikeUncrowded() {
	dsc.useful = dsc.useful[:0]

	for _, priority := range dsc.priorities {
		if dsc.actual[priority] >= dsc.tactic[priority] {
			continue
		}

		dsc.useful = append(dsc.useful, priority)
	}
}
