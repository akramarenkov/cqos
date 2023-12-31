package priority

import (
	"context"
	"sync"
	"time"

	"github.com/akramarenkov/cqos/priority/internal/starter"
)

const (
	defaultChannelCapacity      = 100
	defaultWaitDevastationDelay = 1 * time.Nanosecond
)

type measureDiscipline[Type any] interface {
	Stop()
	Err() <-chan error
}

type measureKind int

const (
	measureKindCompleted measureKind = iota + 1
	measureKindProcessed
	measureKindReceived
)

type measure struct {
	Data         uint
	Kind         measureKind
	Priority     uint
	RelativeTime time.Duration
}

const (
	measuresFactor = 3
)

type actionKind int

const (
	actionKindDelay actionKind = iota + 1
	actionKindWaitDevastation
	actionKindWrite
	actionKindWriteWithDelay
)

type action struct {
	delay    time.Duration
	kind     actionKind
	quantity uint
}

type measurerOpts struct {
	DisableMeasures  bool
	HandlersQuantity uint
	InputCapacity    uint
	NoFeedback       bool
	UnbufferedInput  bool
}

func (opts measurerOpts) normalize() measurerOpts {
	if opts.InputCapacity == 0 {
		opts.InputCapacity = defaultChannelCapacity
	}

	if opts.UnbufferedInput {
		opts.InputCapacity = 0
	}

	return opts
}

type measurer struct {
	opts measurerOpts

	feedback chan uint
	inputs   map[uint]chan uint
	output   chan Prioritized[uint]

	actions map[uint][]action
	delays  map[uint]time.Duration
}

func newMeasurer(opts measurerOpts) *measurer {
	msr := &measurer{
		opts: opts.normalize(),

		feedback: make(chan uint, defaultChannelCapacity),
		inputs:   make(map[uint]chan uint),
		output:   make(chan Prioritized[uint], defaultChannelCapacity),

		actions: make(map[uint][]action),
		delays:  make(map[uint]time.Duration),
	}

	return msr
}

func (msr *measurer) updateInput(priority uint) {
	if _, exists := msr.inputs[priority]; !exists {
		msr.inputs[priority] = make(chan uint, msr.opts.InputCapacity)
	}
}

func (msr *measurer) addAction(priority uint, action action) {
	msr.actions[priority] = append(msr.actions[priority], action)
}

func (msr *measurer) AddWrite(priority uint, quantity uint) {
	msr.updateInput(priority)

	action := action{
		kind:     actionKindWrite,
		quantity: quantity,
	}

	msr.addAction(priority, action)
}

func (msr *measurer) AddWriteWithDelay(priority uint, quantity uint, delay time.Duration) {
	msr.updateInput(priority)

	action := action{
		kind:     actionKindWriteWithDelay,
		quantity: quantity,
		delay:    delay,
	}

	msr.addAction(priority, action)
}

func (msr *measurer) AddWaitDevastation(priority uint) {
	action := action{
		kind: actionKindWaitDevastation,
	}

	msr.addAction(priority, action)
}

func (msr *measurer) AddDelay(priority uint, delay time.Duration) {
	action := action{
		kind:  actionKindDelay,
		delay: delay,
	}

	msr.addAction(priority, action)
}

func (msr *measurer) SetProcessDelay(priority uint, delay time.Duration) {
	msr.delays[priority] = delay
}

func (msr *measurer) GetExpectedItemsQuantity() uint {
	quantity := uint(0)

	for _, actions := range msr.actions {
		for _, action := range actions {
			switch action.kind {
			case actionKindWrite, actionKindWriteWithDelay:
				quantity += action.quantity
			}
		}
	}

	return quantity
}

func (msr *measurer) GetExpectedMeasuresQuantity() uint {
	return measuresFactor * msr.GetExpectedItemsQuantity()
}

func (msr *measurer) GetInputs() map[uint]<-chan uint {
	out := make(map[uint]<-chan uint, len(msr.inputs))

	for priority, channel := range msr.inputs {
		out[priority] = channel
	}

	return out
}

func (msr *measurer) GetOutput() chan<- Prioritized[uint] {
	return msr.output
}

func (msr *measurer) GetFeedback() <-chan uint {
	return msr.feedback
}

func (msr *measurer) runWriters(ctx context.Context, wg *sync.WaitGroup) {
	for priority := range msr.inputs {
		wg.Add(1)

		go msr.writer(ctx, wg, priority)
	}
}

func (msr *measurer) writer(ctx context.Context, wg *sync.WaitGroup, priority uint) {
	defer wg.Done()
	defer close(msr.inputs[priority])

	sequence := uint(0)

	for _, action := range msr.actions[priority] {
		switch action.kind {
		case actionKindWrite, actionKindWriteWithDelay:
			increased, proceed := write(ctx, action, msr.inputs[priority], sequence)
			if !proceed {
				return
			}

			sequence = increased
		case actionKindWaitDevastation:
			if proceed := waitDevastation(ctx, msr.inputs[priority]); !proceed {
				return
			}
		case actionKindDelay:
			time.Sleep(action.delay)
		}
	}
}

func write(
	ctx context.Context,
	action action,
	channel chan uint,
	sequence uint,
) (uint, bool) {
	for id := uint(0); id < action.quantity; id++ {
		select {
		case <-ctx.Done():
			return sequence, false
		case channel <- sequence:
		}

		if action.kind == actionKindWriteWithDelay {
			time.Sleep(action.delay)
		}

		sequence++
	}

	return sequence, true
}

func waitDevastation(ctx context.Context, channel chan uint) bool {
	ticker := time.NewTicker(defaultWaitDevastationDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			if len(channel) == 0 {
				return true
			}
		}
	}
}

func (msr *measurer) runHandlers(
	ctx context.Context,
	wg *sync.WaitGroup,
	channel chan measure,
) {
	starter := starter.New()
	defer starter.Go()

	for id := uint(0); id < msr.opts.HandlersQuantity; id++ {
		wg.Add(1)
		starter.Ready(1)

		go msr.handler(ctx, wg, channel, starter)
	}
}

func (msr *measurer) handler(
	ctx context.Context,
	wg *sync.WaitGroup,
	channel chan measure,
	starter *starter.Starter,
) {
	defer wg.Done()

	starter.Set()

	for {
		select {
		case <-ctx.Done():
			return
		case item, opened := <-msr.output:
			if !opened {
				return
			}

			msr.handle(item, channel, starter)
		}
	}
}

func (msr *measurer) handle(
	item Prioritized[uint],
	channel chan measure,
	starter *starter.Starter,
) {
	if msr.opts.DisableMeasures {
		channel <- measure{}
		channel <- measure{}
		channel <- measure{}

		if !msr.opts.NoFeedback {
			msr.feedback <- item.Priority
		}

		return
	}

	received := measure{
		RelativeTime: time.Since(starter.StartedAt),
		Priority:     item.Priority,
		Kind:         measureKindReceived,
		Data:         item.Item,
	}

	channel <- received

	time.Sleep(msr.delays[item.Priority])

	processed := measure{
		RelativeTime: time.Since(starter.StartedAt),
		Priority:     item.Priority,
		Kind:         measureKindProcessed,
		Data:         item.Item,
	}

	channel <- processed

	if !msr.opts.NoFeedback {
		msr.feedback <- item.Priority
	}

	completed := measure{
		RelativeTime: time.Since(starter.StartedAt),
		Priority:     item.Priority,
		Kind:         measureKindCompleted,
		Data:         item.Item,
	}

	channel <- completed
}

func (msr *measurer) prepare() (chan measure, []measure) {
	quantity := msr.GetExpectedMeasuresQuantity()

	if msr.opts.DisableMeasures {
		return make(chan measure, quantity), make([]measure, 0)
	}

	return make(chan measure, quantity), make([]measure, 0, quantity)
}

func (msr *measurer) Play(
	discipline measureDiscipline[uint],
	incomplete bool,
) []measure {
	defer close(msr.feedback)

	channel, measures := msr.prepare()
	defer close(channel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	msr.runWriters(ctx, wg)
	msr.runHandlers(ctx, wg, channel)

	defer close(msr.output)

	received := 0

	for {
		select {
		case err := <-discipline.Err():
			if err != nil || incomplete {
				cancel()
				return measures
			}
		case measure := <-channel:
			if !msr.opts.DisableMeasures {
				measures = append(measures, measure)
			}

			received++

			if received == cap(channel) {
				return measures
			}
		}
	}
}
