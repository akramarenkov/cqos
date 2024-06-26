// Discipline used to accumulate elements from an input channel into a slice and
// write that slice to an output channel when the maximum slice size or timeout
// for its accumulation is reached.
package join

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/akramarenkov/cqos/join/internal/common"

	"github.com/akramarenkov/breaker"
)

var (
	ErrEmptyInput      = errors.New("input channel was not specified")
	ErrInvalidJoinSize = errors.New("invalid join size")
)

// Options of the created discipline.
type Opts[Type any] struct {
	// Roughly terminates (cancels) work of the discipline
	Ctx context.Context
	// Input data channel. For graceful termination it is enough to close the input
	// channel. Preferably input channel should be buffered for performance reasons.
	// Optimal capacity is in the range of one to two JoinSize
	Input <-chan Type
	// Maximum size of the output slice. Actual size of the output slice may be
	// smaller due to the timeout or closure of the input channel
	JoinSize uint
	// By default, to the output channel is written a copy of the accumulated slice
	// If the Released channel is set, then to the output channel will be directly
	// written the accumulated slice. In this case, after the accumulated slice is
	// used it is necessary to inform the discipline about it by writing
	// to Released channel
	Released <-chan struct{}
	// Timeout for slice accumulation. If the slice has not been filled completely
	// in the allotted time, the data accumulated during this time is written to
	// the output channel. A zero or negative value means that discipline will wait
	// for the missing data until they appear or the channel is closed (in this case,
	// the accumulated data will be written to the output channel)
	Timeout time.Duration
	// Due to the fact that it is not possible to reliably reset the timer/ticker
	// (without false ticks), a ticker with a duration several times shorter than
	// the timeout is used and to determine the expiration of the timeout,
	// the current time is compared with the time of the last writing to
	// the output channel. This method has an inaccuracy that can be set by
	// this parameter in percents
	TimeoutInaccuracy uint
}

func (opts Opts[Type]) isValid() error {
	if opts.Input == nil {
		return ErrEmptyInput
	}

	if opts.JoinSize == 0 {
		return ErrInvalidJoinSize
	}

	return nil
}

func (opts Opts[Type]) normalize() Opts[Type] {
	if opts.Ctx == nil {
		opts.Ctx = context.Background()
	}

	if opts.TimeoutInaccuracy == 0 {
		opts.TimeoutInaccuracy = common.DefaultTimeoutInaccuracy
	}

	return opts
}

// Join discipline.
type Discipline[Type any] struct {
	opts Opts[Type]

	breaker *breaker.Breaker

	interruptInterval time.Duration
	join              []Type
	output            chan []Type
	passAt            time.Time

	unreleased bool
}

// Creates and runs discipline.
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	opts = opts.normalize()

	interval, err := calcInterruptIntervalNonPositiveAllowed(
		opts.Timeout,
		opts.TimeoutInaccuracy,
	)
	if err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		breaker: breaker.New(),

		interruptInterval: interval,
		join:              make([]Type, 0, opts.JoinSize),
		output:            make(chan []Type, 1),
	}

	dsc.resetPassAt()

	go dsc.main()

	return dsc, nil
}

// Returns output channel.
//
// If this channel is closed, it means that the discipline is terminated.
func (dsc *Discipline[Type]) Output() <-chan []Type {
	return dsc.output
}

// Roughly terminates work of the discipline.
//
// Use for wait completion at terminates via context.
func (dsc *Discipline[Type]) Stop() {
	dsc.breaker.Break()
}

func (dsc *Discipline[Type]) main() {
	defer dsc.breaker.Complete()
	defer close(dsc.output)

	if dsc.interruptInterval == 0 {
		dsc.loopUntimeouted()
		return
	}

	dsc.loop()
}

func (dsc *Discipline[Type]) loop() {
	defer dsc.pass()

	ticker := time.NewTicker(dsc.interruptInterval)
	defer ticker.Stop()

	for {
		select {
		case <-dsc.breaker.IsBreaked():
			return
		case <-dsc.opts.Ctx.Done():
			return
		case <-ticker.C:
			if dsc.isTimeouted() {
				dsc.pass()
			}
		case item, opened := <-dsc.opts.Input:
			if !opened {
				return
			}

			dsc.process(item)
		}
	}
}

func (dsc *Discipline[Type]) loopUntimeouted() {
	defer dsc.pass()

	for {
		select {
		case <-dsc.breaker.IsBreaked():
			return
		case <-dsc.opts.Ctx.Done():
			return
		case item, opened := <-dsc.opts.Input:
			if !opened {
				return
			}

			dsc.process(item)
		}
	}
}

func (dsc *Discipline[Type]) process(item Type) {
	if dsc.unreleased {
		return
	}

	dsc.join = append(dsc.join, item)

	if len(dsc.join) < int(dsc.opts.JoinSize) {
		return
	}

	dsc.pass()
}

func (dsc *Discipline[Type]) pass() {
	if dsc.unreleased {
		return
	}

	defer dsc.resetPassAt()

	if len(dsc.join) == 0 {
		return
	}

	dsc.send(dsc.join)
	dsc.resetJoin()
}

func (dsc *Discipline[Type]) send(item []Type) {
	item = dsc.prepareItem(item)

	select {
	case <-dsc.breaker.IsBreaked():
		return
	case <-dsc.opts.Ctx.Done():
		return
	case dsc.output <- item:
	}

	if dsc.opts.Released != nil {
		select {
		case <-dsc.breaker.IsBreaked():
			dsc.unreleased = true
			return
		case <-dsc.opts.Ctx.Done():
			dsc.unreleased = true
			return
		case <-dsc.opts.Released:
		}
	}
}

func (dsc *Discipline[Type]) prepareItem(item []Type) []Type {
	if dsc.opts.Released != nil {
		return item
	}

	return slices.Clone(item)
}

func (dsc *Discipline[Type]) resetJoin() {
	if dsc.unreleased {
		return
	}

	dsc.join = dsc.join[:0]
}

func (dsc *Discipline[Type]) resetPassAt() {
	dsc.passAt = time.Now()
}

func (dsc *Discipline[Type]) isTimeouted() bool {
	return time.Since(dsc.passAt) >= dsc.opts.Timeout
}
