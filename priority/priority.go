// Discipline that used to distributes data among handlers according to priority.
package priority

import (
	"context"
	"errors"
	"time"

	"github.com/akramarenkov/breaker"
	"github.com/akramarenkov/cqos/internal/general"
	"github.com/akramarenkov/cqos/priority/internal/common"
)

var (
	ErrEmptyDivider         = errors.New("priorities divider was not specified")
	ErrEmptyFeedback        = errors.New("feedback channel was not specified")
	ErrEmptyOutput          = errors.New("output channel was not specified")
	ErrHandlersQuantityZero = errors.New("handlers quantity is zero")
	ErrQuantityExceeded     = errors.New("value of handlers quantity has been exceeded")
)

const (
	defaultFeedbackLimitDivider = 10
	defaultIdleDelay            = 1 * time.Nanosecond
	defaultInterruptTimeout     = 1 * time.Nanosecond
)

type inputAdd[Type any] struct {
	channel  <-chan Type
	priority uint
}

// Describes the data distributed by the prioritization discipline.
type Prioritized[Type any] struct {
	Item     Type
	Priority uint
}

// Options of the created discipline.
type Opts[Type any] struct {
	// Roughly terminates (cancels) work of the discipline
	Ctx context.Context
	// Determines how handlers are distributed among priorities
	Divider Divider
	// Handlers must write priority of processed data to feedback channel after it has been processed
	Feedback <-chan uint
	// Between how many handlers you need to distribute data
	HandlersQuantity uint
	// Channels with input data, should be buffered for performance reasons
	// Map key is a value of priority
	// For graceful termination need close all input channels or remove them
	Inputs map[uint]<-chan Type
	// Handlers should read distributed data from this channel
	Output chan<- Prioritized[Type]
}

// Prioritization discipline.
//
// Preferably input channels should be buffered for performance reasons.
//
// Data from input channels passed to handlers by output channel.
//
// Handlers must write priority of processed data to feedback channel after it has been processed.
//
// For equaling use FairDivider, for prioritization use RateDivider or custom divider.
type Discipline[Type any] struct {
	opts Opts[Type]

	breaker  *breaker.Breaker
	graceful *breaker.Breaker

	inputs map[uint]common.Input[Type]

	priorities []uint

	inputAdds chan inputAdd[Type]
	inputRmvs chan uint

	actual    map[uint]uint
	strategic map[uint]uint
	tactic    map[uint]uint

	uncrowded []uint
	useful    []uint

	feedbackLimit uint

	interrupter *time.Ticker

	err chan error
}

func (opts Opts[Type]) isValid() error {
	if opts.Divider == nil {
		return ErrEmptyDivider
	}

	if opts.HandlersQuantity == 0 {
		return ErrHandlersQuantityZero
	}

	if opts.Feedback == nil {
		return ErrEmptyFeedback
	}

	if opts.Output == nil {
		return ErrEmptyOutput
	}

	return nil
}

func (opts Opts[Type]) normalize() Opts[Type] {
	if opts.Ctx == nil {
		opts.Ctx = context.Background()
	}

	return opts
}

// Creates and runs discipline.
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	feedbackLimit := general.DivideWithMin(
		opts.HandlersQuantity,
		defaultFeedbackLimitDivider,
		1,
	)

	dsc := &Discipline[Type]{
		opts: opts.normalize(),

		breaker:  breaker.New(),
		graceful: breaker.New(),

		inputs: make(map[uint]common.Input[Type]),

		inputAdds: make(chan inputAdd[Type]),
		inputRmvs: make(chan uint),

		actual:    make(map[uint]uint),
		strategic: make(map[uint]uint),
		tactic:    make(map[uint]uint),

		feedbackLimit: feedbackLimit,

		interrupter: time.NewTicker(defaultInterruptTimeout),

		err: make(chan error, 1),
	}

	dsc.updateInputs(opts.Inputs)

	go dsc.main()

	return dsc, nil
}

// Returns a channel with errors. If an error occurs (the value from the channel
// is not equal to nil) the discipline terminates its work. The most likely cause of
// the error is an incorrectly working dividing function in which the sum of
// the distributed quantities is not equal to the original quantity.
//
// The single nil value means that the discipline has terminated in normal mode.
//
// If you are sure that the divider is working correctly, then you don’t have to
// read from this channel and you don’t have to check the received value.
func (dsc *Discipline[Type]) Err() <-chan error {
	return dsc.err
}

// Roughly terminates work of the discipline.
//
// Use for wait completion at terminates via context.
func (dsc *Discipline[Type]) Stop() {
	dsc.breaker.Break()
}

// Graceful terminates work of the discipline.
//
// Waits draining input channels, waits end processing data in handlers and terminates.
//
// You must end write to input channels and close them (or remove), otherwise graceful
// stop not be ended.
func (dsc *Discipline[Type]) GracefulStop() {
	dsc.graceful.Break()
}

func (dsc *Discipline[Type]) addPriority(channel <-chan Type, priority uint) {
	_, exists := dsc.inputs[priority]

	input := common.Input[Type]{
		Channel: channel,
	}

	dsc.inputs[priority] = input

	if exists {
		return
	}

	dsc.priorities = append(dsc.priorities, priority)
}

func (dsc *Discipline[Type]) updateInputs(inputs map[uint]<-chan Type) {
	for priority, channel := range inputs {
		dsc.addPriority(channel, priority)
	}

	common.SortPriorities(dsc.priorities)

	dsc.strategic = dsc.opts.Divider(dsc.priorities, dsc.opts.HandlersQuantity, nil)
}

// Adds or updates (if it added previously) input channel for specified priority.
func (dsc *Discipline[Type]) AddInput(channel <-chan Type, priority uint) {
	in := inputAdd[Type]{
		channel:  channel,
		priority: priority,
	}

	dsc.inputAdds <- in
}

func (dsc *Discipline[Type]) addInput(channel <-chan Type, priority uint) {
	dsc.addPriority(channel, priority)

	common.SortPriorities(dsc.priorities)

	dsc.strategic = dsc.opts.Divider(dsc.priorities, dsc.opts.HandlersQuantity, nil)
}

// Removes input channel for specified priority.
func (dsc *Discipline[Type]) RemoveInput(priority uint) {
	dsc.inputRmvs <- priority
}

func (dsc *Discipline[Type]) removeInput(priority uint) {
	delete(dsc.inputs, priority)
	delete(dsc.tactic, priority)

	dsc.priorities = removePriority(dsc.priorities, priority)
	dsc.strategic = dsc.opts.Divider(dsc.priorities, dsc.opts.HandlersQuantity, nil)
}

func (dsc *Discipline[Type]) main() {
	defer dsc.breaker.Complete()
	defer dsc.graceful.Complete()
	defer close(dsc.err)
	defer close(dsc.inputAdds)
	defer close(dsc.inputRmvs)
	defer dsc.interrupter.Stop()

	if err := dsc.loop(); err != nil {
		dsc.err <- err
	}
}

func (dsc *Discipline[Type]) loop() error {
	defer dsc.waitZeroActual()

	for {
		select {
		case <-dsc.breaker.IsBreaked():
			return nil
		case <-dsc.opts.Ctx.Done():
			return nil
		case add := <-dsc.inputAdds:
			dsc.addInput(add.channel, add.priority)
		case priority := <-dsc.inputRmvs:
			dsc.removeInput(priority)
		case priority := <-dsc.opts.Feedback:
			dsc.decreaseActual(priority)
		default:
		}

		dsc.clearActual()

		processed, err := dsc.base()
		if err != nil {
			return err
		}

		if processed == 0 {
			select {
			case <-dsc.graceful.IsBreaked():
				if dsc.isDrainedInputs() {
					return nil
				}
			default:
			}

			time.Sleep(defaultIdleDelay)
		}

		dsc.getLimitedFeedback()
	}
}

func (dsc *Discipline[Type]) waitZeroActual() {
	for !dsc.isZeroActual() {
		select {
		case <-dsc.breaker.IsBreaked():
			return
		case <-dsc.opts.Ctx.Done():
			return
		case priority := <-dsc.opts.Feedback:
			dsc.decreaseActual(priority)
		}
	}
}

func (dsc *Discipline[Type]) getOneFeedback() {
	select {
	case <-dsc.breaker.IsBreaked():
		return
	case <-dsc.opts.Ctx.Done():
		return
	case priority := <-dsc.opts.Feedback:
		dsc.decreaseActual(priority)
	}
}

func (dsc *Discipline[Type]) getLimitedFeedback() {
	for range dsc.feedbackLimit {
		select {
		case <-dsc.breaker.IsBreaked():
			return
		case <-dsc.opts.Ctx.Done():
			return
		case priority := <-dsc.opts.Feedback:
			dsc.decreaseActual(priority)
		default:
			return
		}
	}
}

func (dsc *Discipline[Type]) clearActual() {
	for priority, quantity := range dsc.actual {
		if quantity == 0 && !dsc.isInputExists(priority) {
			delete(dsc.actual, priority)
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

func (dsc *Discipline[Type]) isInputExists(priority uint) bool {
	_, exists := dsc.inputs[priority]
	return exists
}

func (dsc *Discipline[Type]) base() (uint, error) {
	processed := uint(0)

	if err := dsc.waitCalcTactic(); err != nil {
		return processed, err
	}

	processed += dsc.prioritize()

	proceed, err := dsc.recalcTactic()
	if err != nil {
		return processed, err
	}

	if !proceed {
		return processed, nil
	}

	processed += dsc.prioritize()

	return processed, nil
}

func (dsc *Discipline[Type]) waitCalcTactic() error {
	for {
		proceed, err := dsc.calcTactic()
		if err != nil {
			return err
		}

		if proceed {
			return nil
		}

		dsc.getOneFeedback()
	}
}

func (dsc *Discipline[Type]) prioritize() uint {
	processed := uint(0)

	for _, priority := range dsc.priorities {
		if dsc.inputs[priority].Drained {
			continue
		}

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

	for dsc.tactic[priority] != 0 {
		select {
		case <-dsc.breaker.IsBreaked():
			return processed
		case <-dsc.opts.Ctx.Done():
			return processed
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

	return processed
}

func (dsc *Discipline[Type]) iou(priority uint) uint {
	processed := uint(0)

	interrupt := false

	for dsc.tactic[priority] != 0 {
		select {
		case <-dsc.breaker.IsBreaked():
			return processed
		case <-dsc.opts.Ctx.Done():
			return processed
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

	return processed
}

func (dsc *Discipline[Type]) markInputAsDrained(priority uint) {
	input := dsc.inputs[priority]

	input.Drained = true

	dsc.inputs[priority] = input
}

func (dsc *Discipline[Type]) send(item Type, priority uint) uint {
	prioritized := Prioritized[Type]{
		Priority: priority,
		Item:     item,
	}

	select {
	case <-dsc.breaker.IsBreaked():
		return 0
	case <-dsc.opts.Ctx.Done():
		return 0
	case dsc.opts.Output <- prioritized:
		dsc.decreaseTactic(priority)
		dsc.increaseActual(priority)
	}

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

	if picked := dsc.calcTacticByAddUpToStrategic(vacants); picked {
		return true, nil
	}

	return dsc.calcTacticBase(vacants)
}

func (dsc *Discipline[Type]) calcVacants() (uint, error) {
	busy := calcDistributionQuantity(dsc.actual)

	// we will not get an overflow because the correspondence of the quantities is
	// checked at all stages of distribution
	if dsc.opts.HandlersQuantity < busy {
		return 0, ErrQuantityExceeded
	}

	return dsc.opts.HandlersQuantity - busy, nil
}

func (dsc *Discipline[Type]) calcTacticByAddUpToStrategic(vacants uint) bool {
	dsc.resetTactic()

	picked := uint(0)

	for _, priority := range dsc.priorities {
		if dsc.actual[priority] > dsc.strategic[priority] {
			return false
		}

		dsc.tactic[priority] = dsc.strategic[priority] - dsc.actual[priority]

		picked += dsc.tactic[priority]
	}

	return picked == vacants
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
		if dsc.actual[priority] < dsc.strategic[priority] {
			dsc.uncrowded = append(dsc.uncrowded, priority)
		}
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
	remainder := calcDistributionQuantity(dsc.tactic)

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

func (dsc *Discipline[Type]) updateUseful() {
	dsc.useful = dsc.useful[:0]

	for _, priority := range dsc.priorities {
		if dsc.tactic[priority] == 0 {
			dsc.useful = append(dsc.useful, priority)
		}
	}
}

func (dsc *Discipline[Type]) updateUsefulLikeUncrowded() {
	dsc.useful = dsc.useful[:0]

	for _, priority := range dsc.priorities {
		if dsc.actual[priority] < dsc.tactic[priority] {
			dsc.useful = append(dsc.useful, priority)
		}
	}
}
