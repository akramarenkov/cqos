// Internal package with implementation of the measurer which is used for testing.
package measurer

import (
	"context"
	"sync"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/internal/starter"
	"github.com/akramarenkov/cqos/v2/priority/types"
)

const (
	defaultChannelCapacity      = 100
	defaultWaitDevastationDelay = 1 * time.Nanosecond
)

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

type Opts struct {
	DisableMeasures  bool
	HandlersQuantity uint
	InputCapacity    uint
	UnbufferedInput  bool
}

func (opts Opts) normalize() Opts {
	if opts.InputCapacity == 0 {
		opts.InputCapacity = defaultChannelCapacity
	}

	if opts.UnbufferedInput {
		opts.InputCapacity = 0
	}

	return opts
}

type Measurer struct {
	opts Opts

	inputs map[uint]chan uint

	actions map[uint][]action
	delays  map[uint]time.Duration
}

func New(opts Opts) *Measurer {
	msr := &Measurer{
		opts: opts.normalize(),

		inputs: make(map[uint]chan uint),

		actions: make(map[uint][]action),
		delays:  make(map[uint]time.Duration),
	}

	return msr
}

func (msr *Measurer) updateInput(priority uint) {
	if _, exists := msr.inputs[priority]; !exists {
		msr.inputs[priority] = make(chan uint, msr.opts.InputCapacity)
	}
}

func (msr *Measurer) addAction(priority uint, action action) {
	msr.actions[priority] = append(msr.actions[priority], action)
}

func (msr *Measurer) AddWrite(priority uint, quantity uint) {
	msr.updateInput(priority)

	action := action{
		kind:     actionKindWrite,
		quantity: quantity,
	}

	msr.addAction(priority, action)
}

func (msr *Measurer) AddWriteWithDelay(priority uint, quantity uint, delay time.Duration) {
	msr.updateInput(priority)

	action := action{
		kind:     actionKindWriteWithDelay,
		quantity: quantity,
		delay:    delay,
	}

	msr.addAction(priority, action)
}

func (msr *Measurer) AddWaitDevastation(priority uint) {
	action := action{
		kind: actionKindWaitDevastation,
	}

	msr.addAction(priority, action)
}

func (msr *Measurer) AddDelay(priority uint, delay time.Duration) {
	action := action{
		kind:  actionKindDelay,
		delay: delay,
	}

	msr.addAction(priority, action)
}

func (msr *Measurer) SetProcessDelay(priority uint, delay time.Duration) {
	msr.delays[priority] = delay
}

// Number of data elements that must be sent to input channels.
func (msr *Measurer) GetExpectedItemsQuantity() uint {
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

// Number of measure elements to be obtained.
func (msr *Measurer) GetExpectedMeasuresQuantity() uint {
	if msr.opts.DisableMeasures {
		return 0
	}

	return measuresFactor * msr.GetExpectedItemsQuantity()
}

func (msr *Measurer) GetInputs() map[uint]<-chan uint {
	out := make(map[uint]<-chan uint, len(msr.inputs))

	for priority, channel := range msr.inputs {
		out[priority] = channel
	}

	return out
}

func (msr *Measurer) runWriters(ctx context.Context, wg *sync.WaitGroup) {
	for priority := range msr.inputs {
		wg.Add(1)

		go msr.writer(ctx, wg, priority)
	}
}

func (msr *Measurer) writer(ctx context.Context, wg *sync.WaitGroup, priority uint) {
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

func (msr *Measurer) runHandlers(
	ctx context.Context,
	wg *sync.WaitGroup,
	channel chan Measure,
	discipline Discipline[uint],
) {
	starter := starter.New()
	defer starter.Go()

	for id := uint(0); id < msr.opts.HandlersQuantity; id++ {
		wg.Add(1)
		starter.Ready(1)

		go msr.handler(ctx, wg, channel, starter, discipline)
	}
}

func (msr *Measurer) handler(
	ctx context.Context,
	wg *sync.WaitGroup,
	channel chan Measure,
	starter *starter.Starter,
	discipline Discipline[uint],
) {
	defer wg.Done()

	starter.Set()

	for {
		select {
		case <-ctx.Done():
			return
		case item, opened := <-discipline.Output():
			if !opened {
				return
			}

			msr.handle(item, channel, starter, discipline)
		}
	}
}

func (msr *Measurer) handle(
	item types.Prioritized[uint],
	channel chan Measure,
	starter *starter.Starter,
	discipline Discipline[uint],
) {
	if msr.opts.DisableMeasures {
		discipline.Release(item.Priority)
		return
	}

	received := Measure{
		RelativeTime: time.Since(starter.StartedAt),
		Priority:     item.Priority,
		Kind:         MeasureKindReceived,
		Data:         item.Item,
	}

	channel <- received

	time.Sleep(msr.delays[item.Priority])

	processed := Measure{
		RelativeTime: time.Since(starter.StartedAt),
		Priority:     item.Priority,
		Kind:         MeasureKindProcessed,
		Data:         item.Item,
	}

	channel <- processed

	discipline.Release(item.Priority)

	completed := Measure{
		RelativeTime: time.Since(starter.StartedAt),
		Priority:     item.Priority,
		Kind:         MeasureKindCompleted,
		Data:         item.Item,
	}

	channel <- completed
}

func (msr *Measurer) prepare() (chan Measure, []Measure) {
	if msr.opts.DisableMeasures {
		return make(chan Measure), make([]Measure, 0)
	}

	quantity := msr.GetExpectedMeasuresQuantity()

	return make(chan Measure, quantity), make([]Measure, 0, quantity)
}

func (msr *Measurer) Play(discipline Discipline[uint]) []Measure {
	channel, measures := msr.prepare()
	defer close(channel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := <-discipline.Err(); err != nil {
			cancel()
		}
	}()

	msr.runWriters(ctx, wg)
	msr.runHandlers(ctx, wg, channel, discipline)

	if msr.opts.DisableMeasures {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return measures
		case measure := <-channel:
			measures = append(measures, measure)

			if len(measures) == cap(channel) {
				return measures
			}
		}
	}
}
