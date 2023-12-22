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

	// the value was chosen based on studies of the results of graphical tests
	// an attempt to perform a lower delay leads to an increase in it
	defaultMinimumDelay = 1 * time.Millisecond
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

	output chan Type
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

		output: make(chan Type, capacity),
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

	dsc.loop()
}

func (dsc *Discipline[Type]) loop() {
	delay := time.Duration(0)

	for {
		duration, stop := dsc.process()
		if stop {
			return
		}

		remainder := dsc.opts.Limit.Interval - duration

		delay += remainder

		if delay < defaultMinimumDelay {
			continue
		}

		time.Sleep(delay)

		delay = 0
	}
}

func (dsc *Discipline[Type]) process() (time.Duration, bool) {
	startedAt := time.Now()

	for quantity := uint64(0); quantity < dsc.opts.Limit.Quantity; quantity++ {
		item, opened := <-dsc.opts.Input
		if !opened {
			return 0, true
		}

		dsc.send(item)
	}

	return time.Since(startedAt), false
}

func (dsc *Discipline[Type]) send(item Type) {
	dsc.output <- item
}
