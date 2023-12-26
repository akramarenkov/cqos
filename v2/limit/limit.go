// Discipline that used to limits the speed of passing data elements from the
// input channel to the output channel
package limit

import (
	"errors"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/internal/general"
)

var (
	ErrInputEmpty = errors.New("input channel was not specified")
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
		return ErrInputEmpty
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

		delay = dsc.delay(delay, duration)
	}
}

func (dsc *Discipline[Type]) delay(
	delay time.Duration,
	duration time.Duration,
) time.Duration {
	delay = increaseDelay(delay, dsc.opts.Limit.Interval-duration)

	if delay < consts.ReliablyMeasurableDuration {
		return delay
	}

	time.Sleep(delay)

	return 0
}

func increaseDelay(delay time.Duration, delta time.Duration) time.Duration {
	increased := delay + delta

	switch {
	case delay > 0 && delta > 0:
		if increased < delay {
			return 0
		}
	case delay < 0 && delta < 0:
		if increased > delay {
			return 0
		}
	}

	return increased
}

func (dsc *Discipline[Type]) process() (time.Duration, bool) {
	startedAt := time.Now()

	if stop := dsc.pass(); stop {
		return 0, true
	}

	return time.Since(startedAt), false
}

func (dsc *Discipline[Type]) pass() bool {
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
