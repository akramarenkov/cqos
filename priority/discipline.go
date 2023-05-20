package priority

import (
	"errors"
	"sync"
	"time"

	"github.com/akramarenkov/cqos/priority/divider"
	"github.com/akramarenkov/cqos/types"
)

var (
	ErrEmptyDivider  = errors.New("priorities divider was not specified")
	ErrEmptyFeedback = errors.New("feedback channel was not specified")
	ErrEmptyOutput   = errors.New("output channel was not specified")
)

const (
	defaultIdleDelay        = 1 * time.Nanosecond
	defaultInterruptTimeout = 1 * time.Nanosecond
)

type input[Type any] struct {
	channel  chan Type
	priority uint
}

type Opts[Type any] struct {
	Divider          divider.Divider
	Feedback         <-chan uint
	HandlersQuantity uint
	Inputs           map[uint]<-chan Type
	Output           chan<- types.Prioritized[Type]
}

type Discipline[Type any] struct {
	opts Opts[Type]

	breaker   chan bool
	completer chan bool
	stopMutex *sync.Mutex
	stopped   bool

	inputs     map[uint]<-chan Type
	priorities []uint
	inputAdds  chan input[Type]
	inputRmvs  chan uint

	strategic map[uint]uint
	tactic    map[uint]uint
	actual    map[uint]uint

	uncrowded []uint
	useful    []uint

	interrupter *time.Ticker
	unbuffered  map[uint]bool
}

func (opts Opts[Type]) isValid() error {
	if opts.Divider == nil {
		return ErrEmptyDivider
	}

	if opts.Feedback == nil {
		return ErrEmptyFeedback
	}

	if opts.Output == nil {
		return ErrEmptyOutput
	}

	return nil
}

func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		breaker:   make(chan bool),
		completer: make(chan bool),
		stopMutex: &sync.Mutex{},

		inputs:    make(map[uint]<-chan Type),
		inputAdds: make(chan input[Type]),
		inputRmvs: make(chan uint),

		strategic: make(map[uint]uint),
		tactic:    make(map[uint]uint),
		actual:    make(map[uint]uint),

		interrupter: time.NewTicker(defaultInterruptTimeout),
		unbuffered:  make(map[uint]bool),
	}

	dsc.updateInputs(opts.Inputs)

	go dsc.loop()

	return dsc, nil
}

func (dsc *Discipline[Type]) Stop() {
	dsc.stopMutex.Lock()
	defer dsc.stopMutex.Unlock()

	if dsc.stopped {
		return
	}

	dsc.stop()

	dsc.stopped = true
}

func (dsc *Discipline[Type]) stop() {
	close(dsc.breaker)
	<-dsc.completer
	close(dsc.inputAdds)
	close(dsc.inputRmvs)
	dsc.interrupter.Stop()
}

func (dsc *Discipline[Type]) addPriority(channel <-chan Type, priority uint) {
	_, exists := dsc.inputs[priority]

	dsc.inputs[priority] = channel

	if cap(channel) == 0 {
		dsc.unbuffered[priority] = true
	}

	if exists {
		return
	}

	dsc.priorities = append(dsc.priorities, priority)
}

func (dsc *Discipline[Type]) updateInputs(inputs map[uint]<-chan Type) {
	for priority, channel := range inputs {
		dsc.addPriority(channel, priority)
	}

	sortPriorities(dsc.priorities)

	dsc.strategic = dsc.opts.Divider(dsc.priorities, dsc.opts.HandlersQuantity, nil)
}

func (dsc *Discipline[Type]) AddInput(channel chan Type, priority uint) {
	in := input[Type]{
		channel:  channel,
		priority: priority,
	}

	dsc.inputAdds <- in
}

func (dsc *Discipline[Type]) addInput(channel chan Type, priority uint) {
	dsc.addPriority(channel, priority)

	sortPriorities(dsc.priorities)

	dsc.strategic = dsc.opts.Divider(dsc.priorities, dsc.opts.HandlersQuantity, nil)
}

func (dsc *Discipline[Type]) RemoveInput(priority uint) {
	dsc.inputRmvs <- priority
}

func (dsc *Discipline[Type]) removeInput(priority uint) {
	delete(dsc.inputs, priority)
	delete(dsc.tactic, priority)
	delete(dsc.unbuffered, priority)

	dsc.priorities = removePriority(dsc.priorities, priority)
	dsc.strategic = dsc.opts.Divider(dsc.priorities, dsc.opts.HandlersQuantity, nil)
}

func (dsc *Discipline[Type]) loop() {
	defer close(dsc.completer)

	for {
		select {
		case <-dsc.breaker:
			return
		case add := <-dsc.inputAdds:
			dsc.addInput(add.channel, add.priority)
		case priority := <-dsc.inputRmvs:
			dsc.removeInput(priority)
		case priority := <-dsc.opts.Feedback:
			dsc.decreaseActual(priority)
		default:
		}

		dsc.clearActual()

		if processed := dsc.main(); processed == 0 {
			time.Sleep(defaultIdleDelay)
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

func (dsc *Discipline[Type]) isInputExists(priority uint) bool {
	_, exists := dsc.inputs[priority]
	return exists
}

func (dsc *Discipline[Type]) main() uint {
	processed := uint(0)

	dsc.calcTactic()

	processed += dsc.prioritize()

	if proceed := dsc.recalcTactic(); !proceed {
		return processed
	}

	processed += dsc.prioritize()

	return processed
}

func (dsc *Discipline[Type]) prioritize() uint {
	processed := uint(0)

	for _, priority := range dsc.priorities {
		if !dsc.unbuffered[priority] {
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
		case <-dsc.breaker:
			return processed
		case item, opened := <-dsc.inputs[priority]:
			if !opened {
				return processed
			}

			processed += dsc.send(item, priority)
		case precedency := <-dsc.opts.Feedback:
			dsc.decreaseActual(precedency)
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
		case <-dsc.breaker:
			return processed
		case item, opened := <-dsc.inputs[priority]:
			if !opened {
				return processed
			}

			processed += dsc.send(item, priority)
		case precedency := <-dsc.opts.Feedback:
			dsc.decreaseActual(precedency)
		case <-dsc.interrupter.C:
			if interrupt {
				return processed
			}

			interrupt = true
		}
	}
}

func (dsc *Discipline[Type]) send(item Type, priority uint) uint {
	prioritized := types.Prioritized[Type]{
		Priority: priority,
		Item:     item,
	}

	for {
		select {
		case <-dsc.breaker:
			return 0
		case dsc.opts.Output <- prioritized:
			dsc.decreaseTactic(priority)
			dsc.increaseActual(priority)

			return 1
		case precedency := <-dsc.opts.Feedback:
			dsc.decreaseActual(precedency)
		}
	}
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

func (dsc *Discipline[Type]) calcTactic() {
	for {
		if !dsc.pickUpTactic() {
			select {
			case <-dsc.breaker:
				return
			case priority := <-dsc.opts.Feedback:
				dsc.decreaseActual(priority)
			}

			continue
		}

		return
	}
}

func (dsc *Discipline[Type]) pickUpTactic() bool {
	vacants := dsc.calcVacants()

	if vacants == 0 {
		return false
	}

	if picked := dsc.pickUpTacticSimpleAddition(vacants); picked {
		return true
	}

	return dsc.pickUpTacticBase(vacants)
}

func (dsc *Discipline[Type]) pickUpTacticBase(vacants uint) bool {
	dsc.resetTactic()
	dsc.updateUncrowded()
	dsc.opts.Divider(dsc.uncrowded, vacants, dsc.tactic)

	return dsc.isTacticFilled(dsc.uncrowded)
}

func (dsc *Discipline[Type]) pickUpTacticSimpleAddition(vacants uint) bool {
	dsc.resetTactic()

	picked := uint(0)

	for priority := range dsc.actual {
		if dsc.actual[priority] > dsc.strategic[priority] {
			return false
		}

		dsc.tactic[priority] = dsc.strategic[priority] - dsc.actual[priority]

		picked += dsc.tactic[priority]
	}

	return picked != 0 && picked <= vacants
}

func (dsc *Discipline[Type]) isTacticFilled(priorities []uint) bool {
	for _, priority := range priorities {
		if dsc.tactic[priority] == 0 {
			return false
		}
	}

	return true
}

func (dsc *Discipline[Type]) resetTactic() {
	for priority := range dsc.tactic {
		dsc.tactic[priority] = 0
	}
}

func (dsc *Discipline[Type]) calcVacants() uint {
	busy := uint(0)

	for _, quantity := range dsc.actual {
		busy += quantity
	}

	if dsc.opts.HandlersQuantity < busy {
		panic("value of handlers quantity has been exceeded")
	}

	return dsc.opts.HandlersQuantity - busy
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

func (dsc *Discipline[Type]) recalcTactic() bool {
	remainder := dsc.updateUseful()

	dsc.resetTactic()
	dsc.opts.Divider(dsc.useful, dsc.opts.HandlersQuantity, dsc.tactic)

	dsc.reUpdateUseful()
	dsc.resetTactic()
	dsc.opts.Divider(dsc.useful, remainder, dsc.tactic)

	return dsc.isTacticFilled(dsc.useful)
}

func (dsc *Discipline[Type]) updateUseful() uint {
	remainder := uint(0)

	dsc.useful = dsc.useful[:0]

	for _, priority := range dsc.priorities {
		if dsc.tactic[priority] != 0 {
			remainder += dsc.tactic[priority]
			continue
		}

		dsc.useful = append(dsc.useful, priority)
	}

	return remainder
}

func (dsc *Discipline[Type]) reUpdateUseful() {
	dsc.useful = dsc.useful[:0]

	for _, priority := range dsc.priorities {
		if dsc.actual[priority] >= dsc.tactic[priority] {
			continue
		}

		dsc.useful = append(dsc.useful, priority)
	}
}
