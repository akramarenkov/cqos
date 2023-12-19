package limit

import (
	"errors"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/general"
)

var (
	ErrEmptyInput = errors.New("input channel was not specified")
)

const (
	defaultCapacityFactor = 0.1
)

// Options of the created discipline
type Opts[Type any] struct {
	// Input data channel. For terminate discipline it is necessary and sufficient to
	// close the input channel
	Input <-chan Type
	// Rate limit
	Limit Rate
}

func (opts Opts[Type]) isValid() error {
	if opts.Input == nil {
		return ErrEmptyInput
	}

	return opts.Limit.IsValid()
}

// Limit discipline
type Discipline[Type any] struct {
	opts Opts[Type]

	breaker chan struct{}

	output chan Type
	passer chan struct{}
}

// Creates and runs discipline
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	capacity := general.CalcByFactor(
		int(opts.Limit.Quantity),
		defaultCapacityFactor,
		1,
	)

	dsc := &Discipline[Type]{
		opts: opts,

		breaker: make(chan struct{}),

		output: make(chan Type, capacity),
		passer: make(chan struct{}, capacity),
	}

	go dsc.main()

	return dsc, nil
}

// Returns output channel.
//
// If this channel is closed, it means that the discipline is terminated
func (dsc *Discipline[Type]) Output() <-chan Type {
	return dsc.output
}

func (dsc *Discipline[Type]) main() {
	defer close(dsc.output)
	defer close(dsc.passer)

	dsc.loop()
}

func (dsc *Discipline[Type]) loop() {
	ticker := time.NewTicker(dsc.opts.Limit.Interval)
	defer ticker.Stop()

	go dsc.process()

	for {
		select {
		case <-dsc.breaker:
			return
		case <-ticker.C:
			if stop := dsc.pass(); stop {
				return
			}
		}
	}
}

func (dsc *Discipline[Type]) pass() bool {
	for quantity := uint64(0); quantity < dsc.opts.Limit.Quantity; quantity++ {
		select {
		case <-dsc.breaker:
			return true
		case dsc.passer <- struct{}{}:
		}
	}

	return false
}

func (dsc *Discipline[Type]) process() {
	defer close(dsc.breaker)

	for item := range dsc.opts.Input {
		for range dsc.passer {
			dsc.send(item)
			break
		}
	}
}

func (dsc *Discipline[Type]) send(item Type) {
	dsc.output <- item
}
