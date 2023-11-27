package priority

import (
	"context"
	"sync"
	"time"
)

const (
	defaultChannelCapacity      = 100
	defaultWaitDevastationDelay = 1 * time.Microsecond
)

type gaugeKind int

const (
	gaugeKindCompleted gaugeKind = iota + 1
	gaugeKindProcessed
	gaugeKindReceived
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

type gauge struct {
	Data         uint
	Kind         gaugeKind
	Priority     uint
	RelativeTime time.Duration
}

type gaugerOpts struct {
	DisableGauges    bool
	HandlersQuantity uint
	InputCapacity    uint
	NoFeedback       bool
	NoInputBuffer    bool
}

func (opts gaugerOpts) normalize() gaugerOpts {
	if opts.InputCapacity == 0 {
		opts.InputCapacity = defaultChannelCapacity
	}

	if opts.NoInputBuffer {
		opts.InputCapacity = 0
	}

	return opts
}

type gauger struct {
	opts gaugerOpts

	ready     *sync.WaitGroup
	start     chan bool
	startedAt time.Time

	actions map[uint][]action
	delays  map[uint]time.Duration
	gauges  chan []gauge

	feedback chan uint
	inputs   map[uint]chan uint
	output   <-chan Prioritized[uint]

	waiter *sync.WaitGroup
}

func newGauger(opts gaugerOpts) *gauger {
	ggr := &gauger{
		opts: opts.normalize(),

		ready: &sync.WaitGroup{},

		actions: make(map[uint][]action),
		delays:  make(map[uint]time.Duration),

		feedback: make(chan uint, defaultChannelCapacity),
		inputs:   make(map[uint]chan uint),

		waiter: &sync.WaitGroup{},
	}

	return ggr
}

func (ggr *gauger) AddWrite(priority uint, quantity uint) {
	if _, exists := ggr.inputs[priority]; !exists {
		ggr.inputs[priority] = make(chan uint, ggr.opts.InputCapacity)
	}

	action := action{
		kind:     actionKindWrite,
		quantity: quantity,
	}

	ggr.actions[priority] = append(ggr.actions[priority], action)
}

func (ggr *gauger) AddWriteWithDelay(priority uint, quantity uint, delay time.Duration) {
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

func (ggr *gauger) AddWaitDevastation(priority uint) {
	action := action{
		kind: actionKindWaitDevastation,
	}

	ggr.actions[priority] = append(ggr.actions[priority], action)
}

func (ggr *gauger) AddDelay(priority uint, delay time.Duration) {
	action := action{
		kind:  actionKindDelay,
		delay: delay,
	}

	ggr.actions[priority] = append(ggr.actions[priority], action)
}

func (ggr *gauger) CalcExpectedGuagesQuantity() uint {
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

func (ggr *gauger) SetProcessDelay(priority uint, delay time.Duration) {
	ggr.delays[priority] = delay
}

func (ggr *gauger) GetInputs() map[uint]<-chan uint {
	out := make(map[uint]<-chan uint, len(ggr.inputs))

	for priority, channel := range ggr.inputs {
		out[priority] = channel
	}

	return out
}

func (ggr *gauger) SetOutput(output <-chan Prioritized[uint]) {
	ggr.output = output
}

func (ggr *gauger) GetFeedback() <-chan uint {
	return ggr.feedback
}

func (ggr *gauger) runWriters(ctx context.Context) {
	for priority := range ggr.inputs {
		ggr.waiter.Add(1)

		go ggr.writer(ctx, priority)
	}
}

func (ggr *gauger) writer(ctx context.Context, priority uint) {
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

func (ggr *gauger) runHandlers(ctx context.Context) {
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

func (ggr *gauger) handler(ctx context.Context) {
	defer ggr.waiter.Done()

	ggr.ready.Done()

	<-ggr.start

	const batchSize = 3

	for {
		select {
		case <-ctx.Done():
			return
		case prioritized, opened := <-ggr.output:
			if !opened {
				return
			}

			if ggr.opts.DisableGauges {
				ggr.feedback <- prioritized.Priority
				ggr.gauges <- nil

				continue
			}

			batch := make([]gauge, 0, batchSize)

			received := gauge{
				RelativeTime: time.Since(ggr.startedAt),
				Priority:     prioritized.Priority,
				Kind:         gaugeKindReceived,
				Data:         prioritized.Item,
			}

			batch = append(batch, received)

			time.Sleep(ggr.delays[prioritized.Priority])

			processed := gauge{
				RelativeTime: time.Since(ggr.startedAt),
				Priority:     prioritized.Priority,
				Kind:         gaugeKindProcessed,
				Data:         prioritized.Item,
			}

			batch = append(batch, processed)

			if !ggr.opts.NoFeedback {
				ggr.feedback <- prioritized.Priority
			}

			completed := gauge{
				RelativeTime: time.Since(ggr.startedAt),
				Priority:     prioritized.Priority,
				Kind:         gaugeKindCompleted,
				Data:         prioritized.Item,
			}

			batch = append(batch, completed)

			ggr.gauges <- batch
		}
	}
}

func (ggr *gauger) Play(ctx context.Context) []gauge {
	defer close(ggr.feedback)

	expectedGaugesQuantity := ggr.CalcExpectedGuagesQuantity()

	if expectedGaugesQuantity == 0 {
		return nil
	}

	gaugesCapacity := expectedGaugesQuantity

	if ggr.opts.DisableGauges {
		gaugesCapacity = 0
	}

	ggr.gauges = make(chan []gauge, expectedGaugesQuantity)

	received := uint(0)
	gauges := make([]gauge, 0, gaugesCapacity)

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
