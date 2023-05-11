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
	Data     uint
	Duration time.Duration
	Kind     GaugeKind
	Priority uint
}

type GaugerOpts struct {
	HandlersQuantity uint
}

type Gauger struct {
	opts GaugerOpts

	breaker   chan bool
	ready     *sync.WaitGroup
	start     chan bool
	startedAt time.Time

	actions map[uint][]action
	delays  map[uint]time.Duration
	results chan []Gauge

	feedback chan uint
	inputs   map[uint]chan uint
	output   chan types.Prioritized[uint]

	waiter *sync.WaitGroup
}

func NewGauger(opts GaugerOpts) *Gauger {
	gg := &Gauger{
		opts: opts,

		ready: &sync.WaitGroup{},

		actions: make(map[uint][]action),
		delays:  make(map[uint]time.Duration),

		feedback: make(chan uint, defaultChannelCapacity),
		inputs:   make(map[uint]chan uint),
		output:   make(chan types.Prioritized[uint], defaultChannelCapacity),

		waiter: &sync.WaitGroup{},
	}

	return gg
}

func (gg *Gauger) Finalize() {
	close(gg.feedback)

	for _, channel := range gg.inputs {
		close(channel)
	}

	close(gg.output)
}

func (gg *Gauger) AddWrite(priority uint, quantity uint) {
	if _, exists := gg.inputs[priority]; !exists {
		gg.inputs[priority] = make(chan uint, defaultChannelCapacity)
	}

	action := action{
		kind:     actionKindWrite,
		quantity: quantity,
	}

	gg.actions[priority] = append(gg.actions[priority], action)
}

func (gg *Gauger) AddWriteWithDelay(priority uint, quantity uint, delay time.Duration) {
	if _, exists := gg.inputs[priority]; !exists {
		gg.inputs[priority] = make(chan uint, defaultChannelCapacity)
	}

	action := action{
		kind:     actionKindWriteWithDelay,
		quantity: quantity,
		delay:    delay,
	}

	gg.actions[priority] = append(gg.actions[priority], action)
}

func (gg *Gauger) AddWaitDevastation(priority uint) {
	action := action{
		kind: actionKindWaitDevastation,
	}

	gg.actions[priority] = append(gg.actions[priority], action)
}

func (gg *Gauger) AddDelay(priority uint, delay time.Duration) {
	action := action{
		kind:  actionKindDelay,
		delay: delay,
	}

	gg.actions[priority] = append(gg.actions[priority], action)
}

func (gg *Gauger) calcExpectedResultsQuantity() uint {
	quantity := uint(0)

	for _, actions := range gg.actions {
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

func (gg *Gauger) SetProcessDelay(priority uint, delay time.Duration) {
	gg.delays[priority] = delay
}

func (gg *Gauger) GetInputs() map[uint]<-chan uint {
	out := make(map[uint]<-chan uint, len(gg.inputs))

	for priority, channel := range gg.inputs {
		out[priority] = channel
	}

	return out
}

func (gg *Gauger) GetOutput() chan<- types.Prioritized[uint] {
	return gg.output
}

func (gg *Gauger) GetFeedback() <-chan uint {
	return gg.feedback
}

func (gg *Gauger) runWriters() {
	for priority := range gg.inputs {
		gg.waiter.Add(1)

		go gg.writer(priority)
	}
}

func (gg *Gauger) writer(priority uint) {
	defer gg.waiter.Done()

	written := uint(0)

	for _, action := range gg.actions[priority] {
		switch action.kind {
		case actionKindWrite:
			for id := uint(0); id < action.quantity; id++ {
				gg.inputs[priority] <- written

				written++
			}
		case actionKindWriteWithDelay:
			for id := uint(0); id < action.quantity; id++ {
				gg.inputs[priority] <- written

				time.Sleep(action.delay)

				written++
			}
		case actionKindWaitDevastation:
			func() {
				ticker := time.NewTicker(defaultWaitDevastationDelay)
				defer ticker.Stop()

				for range ticker.C {
					if len(gg.inputs[priority]) == 0 {
						break
					}
				}
			}()
		case actionKindDelay:
			time.Sleep(action.delay)
		}
	}
}

func (gg *Gauger) runHandlers() {
	gg.start = make(chan bool)
	defer close(gg.start)

	for counter := uint(0); counter < gg.opts.HandlersQuantity; counter++ {
		gg.ready.Add(1)
		gg.waiter.Add(1)

		go gg.handler()
	}

	gg.ready.Wait()

	gg.startedAt = time.Now()
}

func (gg *Gauger) handler() {
	defer gg.waiter.Done()

	gg.ready.Done()

	<-gg.start

	const batchSize = 3

	for {
		select {
		case <-gg.breaker:
			return
		case prioritized := <-gg.output:
			batch := make([]Gauge, 0, batchSize)

			received := Gauge{
				Duration: time.Since(gg.startedAt),
				Priority: prioritized.Priority,
				Kind:     GaugeKindReceived,
				Data:     prioritized.Item,
			}

			batch = append(batch, received)

			time.Sleep(gg.delays[prioritized.Priority])

			processed := Gauge{
				Duration: time.Since(gg.startedAt),
				Priority: prioritized.Priority,
				Kind:     GaugeKindProcessed,
				Data:     prioritized.Item,
			}

			batch = append(batch, processed)

			gg.feedback <- prioritized.Priority

			completed := Gauge{
				Duration: time.Since(gg.startedAt),
				Priority: prioritized.Priority,
				Kind:     GaugeKindCompleted,
				Data:     prioritized.Item,
			}

			batch = append(batch, completed)

			gg.results <- batch
		}
	}
}

func (gg *Gauger) Play() []Gauge {
	expectedResultsQuantity := gg.calcExpectedResultsQuantity()

	if expectedResultsQuantity == 0 {
		return nil
	}

	gg.breaker = make(chan bool)
	gg.results = make(chan []Gauge, expectedResultsQuantity)

	received := uint(0)
	results := make([]Gauge, 0, expectedResultsQuantity)

	gg.runWriters()
	gg.runHandlers()

	defer close(gg.results)
	defer gg.waiter.Wait()
	defer close(gg.breaker)

	for batch := range gg.results {
		results = append(results, batch...)

		received++

		if received == expectedResultsQuantity {
			return results
		}
	}

	return nil
}
