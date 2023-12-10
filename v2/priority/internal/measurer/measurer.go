// Internal package with implementation of the measurer which is used for testing
package measurer

import (
	"context"
	"sync"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/internal/common"
	"github.com/akramarenkov/cqos/v2/priority/internal/starter"
)

const (
	defaultChannelCapacity      = 100
	defaultWaitDevastationDelay = 1 * time.Nanosecond
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
	NoFeedback       bool
	NoInputBuffer    bool
}

func (opts Opts) normalize() Opts {
	if opts.InputCapacity == 0 {
		opts.InputCapacity = defaultChannelCapacity
	}

	if opts.NoInputBuffer {
		opts.InputCapacity = 0
	}

	return opts
}

type Measurer struct {
	opts Opts

	inputs map[uint]chan uint

	actions map[uint][]action
	delays  map[uint]time.Duration

	measures chan []Measure
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

func (msr *Measurer) AddWrite(priority uint, quantity uint) {
	msr.updateInput(priority)

	action := action{
		kind:     actionKindWrite,
		quantity: quantity,
	}

	msr.actions[priority] = append(msr.actions[priority], action)
}

func (msr *Measurer) AddWriteWithDelay(priority uint, quantity uint, delay time.Duration) {
	msr.updateInput(priority)

	action := action{
		kind:     actionKindWriteWithDelay,
		quantity: quantity,
		delay:    delay,
	}

	msr.actions[priority] = append(msr.actions[priority], action)
}

func (msr *Measurer) AddWaitDevastation(priority uint) {
	action := action{
		kind: actionKindWaitDevastation,
	}

	msr.actions[priority] = append(msr.actions[priority], action)
}

func (msr *Measurer) AddDelay(priority uint, delay time.Duration) {
	action := action{
		kind:  actionKindDelay,
		delay: delay,
	}

	msr.actions[priority] = append(msr.actions[priority], action)
}

func (msr *Measurer) SetProcessDelay(priority uint, delay time.Duration) {
	msr.delays[priority] = delay
}

func (msr *Measurer) GetExpectedMeasuresQuantity() uint {
	quantity := uint(0)

	for _, actions := range msr.actions {
		for _, action := range actions {
			switch action.kind {
			case actionKindWrite:
				quantity += action.quantity
			case actionKindWriteWithDelay:
				quantity += action.quantity
			}
		}
	}

	return quantity
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

	written := uint(0)

	for _, action := range msr.actions[priority] {
		switch action.kind {
		case actionKindWrite:
			for id := uint(0); id < action.quantity; id++ {
				select {
				case <-ctx.Done():
					return
				case msr.inputs[priority] <- written:
				}

				written++
			}
		case actionKindWriteWithDelay:
			for id := uint(0); id < action.quantity; id++ {
				select {
				case <-ctx.Done():
					return
				case msr.inputs[priority] <- written:
				}

				time.Sleep(action.delay)

				written++
			}
		case actionKindWaitDevastation:
			func() {
				ticker := time.NewTicker(defaultWaitDevastationDelay)
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						if len(msr.inputs[priority]) == 0 {
							return
						}
					}
				}
			}()
		case actionKindDelay:
			time.Sleep(action.delay)
		}
	}
}

func (msr *Measurer) runHandlers(
	ctx context.Context,
	wg *sync.WaitGroup,
	discipline common.Discipline[uint],
) {
	starter := starter.New()

	for counter := uint(0); counter < msr.opts.HandlersQuantity; counter++ {
		wg.Add(1)
		starter.Ready(1)

		go msr.handler(ctx, wg, starter, discipline)
	}

	starter.Go()
}

func (msr *Measurer) handler(
	ctx context.Context,
	wg *sync.WaitGroup,
	starter *starter.Starter,
	discipline common.Discipline[uint],
) {
	defer wg.Done()

	starter.Set()

	const batchSize = 3

	for {
		select {
		case <-ctx.Done():
			return
		case prioritized, opened := <-discipline.Output():
			if !opened {
				return
			}

			if msr.opts.DisableMeasures {
				discipline.Release(prioritized.Priority)
				msr.measures <- nil

				continue
			}

			batch := make([]Measure, 0, batchSize)

			received := Measure{
				RelativeTime: time.Since(starter.StartedAt),
				Priority:     prioritized.Priority,
				Kind:         MeasureKindReceived,
				Data:         prioritized.Item,
			}

			batch = append(batch, received)

			time.Sleep(msr.delays[prioritized.Priority])

			processed := Measure{
				RelativeTime: time.Since(starter.StartedAt),
				Priority:     prioritized.Priority,
				Kind:         MeasureKindProcessed,
				Data:         prioritized.Item,
			}

			batch = append(batch, processed)

			if !msr.opts.NoFeedback {
				discipline.Release(prioritized.Priority)
			}

			completed := Measure{
				RelativeTime: time.Since(starter.StartedAt),
				Priority:     prioritized.Priority,
				Kind:         MeasureKindCompleted,
				Data:         prioritized.Item,
			}

			batch = append(batch, completed)

			msr.measures <- batch
		}
	}
}

func (msr *Measurer) Play(discipline common.Discipline[uint]) []Measure {
	expectedMeasuresQuantity := msr.GetExpectedMeasuresQuantity()

	if expectedMeasuresQuantity == 0 {
		return nil
	}

	measuresCapacity := expectedMeasuresQuantity

	if msr.opts.DisableMeasures {
		measuresCapacity = 0
	}

	msr.measures = make(chan []Measure, expectedMeasuresQuantity)

	received := uint(0)
	measures := make([]Measure, 0, measuresCapacity)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := <-discipline.Err(); err != nil {
			cancel()
		}
	}()

	wg := &sync.WaitGroup{}

	msr.runWriters(ctx, wg)
	msr.runHandlers(ctx, wg, discipline)

	defer close(msr.measures)
	defer wg.Wait()

	for {
		select {
		case <-ctx.Done():
			return measures
		case batch := <-msr.measures:
			if !msr.opts.DisableMeasures {
				measures = append(measures, batch...)
			}

			received++

			if received == expectedMeasuresQuantity {
				return measures
			}
		}
	}
}
