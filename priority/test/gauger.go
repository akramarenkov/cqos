// Mostly internaly package used to test and research the discipline
package test

import (
	"sync"
	"time"

	"github.com/akramarenkov/cqos/types"
)

const (
	defaultChannelCapacity      = 100
	defaultWaitDevastationDelay = 1 * time.Microsecond
)

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

type GaugerOpts struct {
	DisableGauges    bool
	HandlersQuantity uint
	InputCapacity    uint
	NoFeedback       bool
	NoInputBuffer    bool
}

func (opts GaugerOpts) normalize() GaugerOpts {
	if opts.InputCapacity == 0 {
		opts.InputCapacity = defaultChannelCapacity
	}

	if opts.NoInputBuffer {
		opts.InputCapacity = 0
	}

	return opts
}

type Gauger struct {
	opts GaugerOpts

	breaker   chan bool
	ready     *sync.WaitGroup
	start     chan bool
	startedAt time.Time

	actions map[uint][]action
	delays  map[uint]time.Duration
	gauges  chan []Gauge

	feedback chan uint
	inputs   map[uint]chan uint
	output   chan types.Prioritized[uint]

	waiter *sync.WaitGroup
}

func NewGauger(opts GaugerOpts) *Gauger {
	ggr := &Gauger{
		opts: opts.normalize(),

		ready: &sync.WaitGroup{},

		actions: make(map[uint][]action),
		delays:  make(map[uint]time.Duration),

		feedback: make(chan uint, defaultChannelCapacity),
		inputs:   make(map[uint]chan uint),
		output:   make(chan types.Prioritized[uint], defaultChannelCapacity),

		waiter: &sync.WaitGroup{},
	}

	return ggr
}

func (ggr *Gauger) Finalize() {
	close(ggr.feedback)

	for _, channel := range ggr.inputs {
		close(channel)
	}

	close(ggr.output)
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

func (ggr *Gauger) calcExpectedGuagesQuantity() uint {
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

func (ggr *Gauger) GetOutput() chan<- types.Prioritized[uint] {
	return ggr.output
}

func (ggr *Gauger) GetFeedback() <-chan uint {
	return ggr.feedback
}

func (ggr *Gauger) runWriters() {
	for priority := range ggr.inputs {
		ggr.waiter.Add(1)

		go ggr.writer(priority)
	}
}

func (ggr *Gauger) writer(priority uint) {
	defer ggr.waiter.Done()

	written := uint(0)

	for _, action := range ggr.actions[priority] {
		switch action.kind {
		case actionKindWrite:
			for id := uint(0); id < action.quantity; id++ {
				ggr.inputs[priority] <- written

				written++
			}
		case actionKindWriteWithDelay:
			for id := uint(0); id < action.quantity; id++ {
				ggr.inputs[priority] <- written

				time.Sleep(action.delay)

				written++
			}
		case actionKindWaitDevastation:
			func() {
				ticker := time.NewTicker(defaultWaitDevastationDelay)
				defer ticker.Stop()

				for range ticker.C {
					if len(ggr.inputs[priority]) == 0 {
						break
					}
				}
			}()
		case actionKindDelay:
			time.Sleep(action.delay)
		}
	}
}

func (ggr *Gauger) runHandlers() {
	ggr.start = make(chan bool)
	defer close(ggr.start)

	for counter := uint(0); counter < ggr.opts.HandlersQuantity; counter++ {
		ggr.ready.Add(1)
		ggr.waiter.Add(1)

		go ggr.handler()
	}

	ggr.ready.Wait()

	ggr.startedAt = time.Now()
}

func (ggr *Gauger) handler() {
	defer ggr.waiter.Done()

	ggr.ready.Done()

	<-ggr.start

	const batchSize = 3

	for {
		select {
		case <-ggr.breaker:
			return
		case prioritized := <-ggr.output:
			if ggr.opts.DisableGauges {
				ggr.feedback <- prioritized.Priority
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
				ggr.feedback <- prioritized.Priority
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

func (ggr *Gauger) Play() []Gauge {
	expectedGaugesQuantity := ggr.calcExpectedGuagesQuantity()

	if expectedGaugesQuantity == 0 {
		return nil
	}

	gaugesCapacity := expectedGaugesQuantity

	if ggr.opts.DisableGauges {
		gaugesCapacity = 0
	}

	ggr.breaker = make(chan bool)
	ggr.gauges = make(chan []Gauge, expectedGaugesQuantity)

	received := uint(0)
	gauges := make([]Gauge, 0, gaugesCapacity)

	ggr.runWriters()
	ggr.runHandlers()

	defer close(ggr.gauges)
	defer ggr.waiter.Wait()
	defer close(ggr.breaker)

	for batch := range ggr.gauges {
		if !ggr.opts.DisableGauges {
			gauges = append(gauges, batch...)
		}

		received++

		if received == expectedGaugesQuantity {
			return gauges
		}
	}

	return nil
}
