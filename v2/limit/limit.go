package limit

import (
	"errors"
	"time"
)

var (
	ErrEmptyInput = errors.New("input channel was not specified")
)

const (
	defaultOutputCap = 1
)

// Options of the created discipline
type Opts[Type any] struct {
	// Input data channel. For graceful termination it is enough to
	// close the input channel
	Input <-chan Type
	// Rate limit
	Limit     Rate
	OutputCap uint
	// Do not waits for the first ticker tick and transfer first data batch immediately
	ZeroTick bool
}

func (opts Opts[Type]) isValid() error {
	if opts.Input == nil {
		return ErrEmptyInput
	}

	return opts.Limit.IsValid()
}

func (opts Opts[Type]) normalize() Opts[Type] {
	if opts.OutputCap == 0 {
		opts.OutputCap = defaultOutputCap
	}

	return opts
}

// Main discipline
type Discipline[Type any] struct {
	opts Opts[Type]

	output chan Type
}

// Creates and runs main discipline
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	opts = opts.normalize()

	dsc := &Discipline[Type]{
		opts: opts,

		output: make(chan Type, opts.OutputCap),
	}

	go dsc.loop()

	return dsc, nil
}

// Returns output channel
func (dsc *Discipline[Type]) Output() <-chan Type {
	return dsc.output
}

func (dsc *Discipline[Type]) loop() {
	defer close(dsc.output)

	ticker := time.NewTicker(dsc.opts.Limit.Interval)
	defer ticker.Stop()

	if dsc.opts.ZeroTick {
		if stop := dsc.process(); stop {
			return
		}
	}

	for range ticker.C {
		if stop := dsc.process(); stop {
			return
		}
	}
}

func (dsc *Discipline[Type]) process() bool {
	for quantity := uint64(0); quantity < dsc.opts.Limit.Quantity; quantity++ {
		item, opened := <-dsc.opts.Input
		if !opened {
			return true
		}

		dsc.send(item)
	}

	return false
}

func (dsc *Discipline[Type]) send(item Type) {
	dsc.output <- item
}
