// Internal package with implementation of the gauger which is used for testing
package gauger

import (
	"context"
	"sync"
	"time"

	"github.com/akramarenkov/cqos/v2/priority/types"
)

const (
	defaultChannelCapacity      = 100
	defaultWaitDevastationDelay = 1 * time.Microsecond
)

type disciplineInterface[Type any] interface {
	Output() <-chan types.Prioritized[Type]
	Release(priority uint)
	Err() <-chan error
}

type GaugeKind int

const (
	GaugeKindCompleted GaugeKind = iota + 1
	GaugeKindProcessed
	GaugeKindReceived
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

type Gauge struct {
	Data         uint
	Kind         GaugeKind
	Priority     uint
	RelativeTime time.Duration
}

type Opts struct {
	DisableGauges    bool
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

type Gauger struct {
	opts Opts

	ready     *sync.WaitGroup
	start     chan bool
	startedAt time.Time

	actions map[uint][]action
	delays  map[uint]time.Duration
	gauges  chan []Gauge

	inputs     map[uint]chan uint
	discipline disciplineInterface[uint]

	waiter *sync.WaitGroup
}

func New(opts Opts) *Gauger {
	ggr := &Gauger{
		opts: opts.normalize(),

		ready: &sync.WaitGroup{},

		actions: make(map[uint][]action),
		delays:  make(map[uint]time.Duration),

		inputs: make(map[uint]chan uint),

		waiter: &sync.WaitGroup{},
	}

	return ggr
}

func (ggr *Gauger) AddWrite(priority uint, quantity uint) {
	if _, exists := ggr.inputs[priority]; !exists {
		ggr.inputs[priority] = make(chan uint, ggr.opts.InputCapacity)
	}

	action := action{
		kind:     actionKindWrite,
		quantity: quantity,
	}

	ggr.actions[priority] = append(ggr.actions[priority], action)
}

func (ggr *Gauger) AddWriteWithDelay(priority uint, quantity uint, delay time.Duration) {
	if _, exists := ggr.inputs[priority]; !exists {
		ggr.inputs[priority] = make(chan uint, ggr.opts.InputCapacity)
	}

	action := action{
		kind:     actionKindWriteWithDelay,
		quantity: quantity,
		delay:    delay,
	}

	ggr.actions[priority] = append(ggr.actions[priority], action)
}

func (ggr *Gauger) AddWaitDevastation(priority uint) {
	action := action{
		kind: actionKindWaitDevastation,
	}

	ggr.actions[priority] = append(ggr.actions[priority], action)
}

func (ggr *Gauger) AddDelay(priority uint, delay time.Duration) {
	action := action{
		kind:  actionKindDelay,
		delay: delay,
	}

	ggr.actions[priority] = append(ggr.actions[priority], action)
}

func (ggr *Gauger) CalcExpectedGuagesQuantity() uint {
	quantity := uint(0)

	for _, actions := range ggr.actions {
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

func (ggr *Gauger) SetProcessDelay(priority uint, delay time.Duration) {
	ggr.delays[priority] = delay
}

func (ggr *Gauger) GetInputs() map[uint]<-chan uint {
	out := make(map[uint]<-chan uint, len(ggr.inputs))

	for priority, channel := range ggr.inputs {
		out[priority] = channel
	}

	return out
}

func (ggr *Gauger) SetDiscipline(discipline disciplineInterface[uint]) {
	ggr.discipline = discipline
}

func (ggr *Gauger) runWriters(ctx context.Context) {
	for priority := range ggr.inputs {
		ggr.waiter.Add(1)

		go ggr.writer(ctx, priority)
	}
}

func (ggr *Gauger) writer(ctx context.Context, priority uint) {
	defer ggr.waiter.Done()
	defer close(ggr.inputs[priority])

	written := uint(0)

	for _, action := range ggr.actions[priority] {
		switch action.kind {
		case actionKindWrite:
			for id := uint(0); id < action.quantity; id++ {
				select {
				case <-ctx.Done():
					return
				case ggr.inputs[priority] <- written:
				}

				written++
			}
		case actionKindWriteWithDelay:
			for id := uint(0); id < action.quantity; id++ {
				select {
				case <-ctx.Done():
					return
				case ggr.inputs[priority] <- written:
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
						if len(ggr.inputs[priority]) == 0 {
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

func (ggr *Gauger) runHandlers(ctx context.Context) {
	ggr.start = make(chan bool)
	defer close(ggr.start)

	for counter := uint(0); counter < ggr.opts.HandlersQuantity; counter++ {
		ggr.ready.Add(1)
		ggr.waiter.Add(1)

		go ggr.handler(ctx)
	}

	ggr.ready.Wait()

	ggr.startedAt = time.Now()
}

func (ggr *Gauger) handler(ctx context.Context) {
	defer ggr.waiter.Done()

	ggr.ready.Done()

	<-ggr.start

	const batchSize = 3

	for {
		select {
		case <-ctx.Done():
			return
		case prioritized, opened := <-ggr.discipline.Output():
			if !opened {
				return
			}

			if ggr.opts.DisableGauges {
				ggr.discipline.Release(prioritized.Priority)
				ggr.gauges <- nil

				continue
			}

			batch := make([]Gauge, 0, batchSize)

			received := Gauge{
				RelativeTime: time.Since(ggr.startedAt),
				Priority:     prioritized.Priority,
				Kind:         GaugeKindReceived,
				Data:         prioritized.Item,
			}

			batch = append(batch, received)

			time.Sleep(ggr.delays[prioritized.Priority])

			processed := Gauge{
				RelativeTime: time.Since(ggr.startedAt),
				Priority:     prioritized.Priority,
				Kind:         GaugeKindProcessed,
				Data:         prioritized.Item,
			}

			batch = append(batch, processed)

			if !ggr.opts.NoFeedback {
				ggr.discipline.Release(prioritized.Priority)
			}

			completed := Gauge{
				RelativeTime: time.Since(ggr.startedAt),
				Priority:     prioritized.Priority,
				Kind:         GaugeKindCompleted,
				Data:         prioritized.Item,
			}

			batch = append(batch, completed)

			ggr.gauges <- batch
		}
	}
}

func (ggr *Gauger) Play(ctx context.Context) []Gauge {
	expectedGaugesQuantity := ggr.CalcExpectedGuagesQuantity()

	if expectedGaugesQuantity == 0 {
		return nil
	}

	gaugesCapacity := expectedGaugesQuantity

	if ggr.opts.DisableGauges {
		gaugesCapacity = 0
	}

	ggr.gauges = make(chan []Gauge, expectedGaugesQuantity)

	received := uint(0)
	gauges := make([]Gauge, 0, gaugesCapacity)

	ggr.runWriters(ctx)
	ggr.runHandlers(ctx)

	defer close(ggr.gauges)
	defer ggr.waiter.Wait()

	for {
		select {
		case <-ctx.Done():
			return gauges
		case batch := <-ggr.gauges:
			if !ggr.opts.DisableGauges {
				gauges = append(gauges, batch...)
			}

			received++

			if received == expectedGaugesQuantity {
				return gauges
			}
		}
	}
}
