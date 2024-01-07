// Discipline that used to accumulates elements from the input channel into a slice and
// writes it to the output channel when the size or timeout is reached.
package join

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/akramarenkov/cqos/breaker"
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
	// channel
	Input <-chan Type
	// Output slice size
	JoinSize uint
	// By default, to the output channel is written a copy of the accumulated slice
	// If the Released channel is set, then to the output channel will be directly
	// written the accumulated slice
	// In this case, after the accumulated slice is used it is necessary to inform
	// the discipline about it by writing to Released channel
	Released <-chan struct{}
	// Send timeout of accumulated slice. A zero or negative value means that no data is
	// written to the output channel after the time has elapsed
	Timeout time.Duration
	// Due to the fact that it is not possible to reliably reset the timer/ticker
	// (without false ticks), a ticker with a duration several times shorter than
	// the timeout is used and to determine the expiration of the timeout,
	// the current time is compared with the time of the last recording to
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
		opts.TimeoutInaccuracy = defaultTimeoutInaccuracy
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
	sendAt            time.Time
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

	dsc.resetSendAt()

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
	defer dsc.send()

	ticker := time.NewTicker(dsc.interruptInterval)
	defer ticker.Stop()

	for {
		select {
		case <-dsc.breaker.Breaked():
			return
		case <-dsc.opts.Ctx.Done():
			return
		case <-ticker.C:
			if dsc.isTimeouted() {
				dsc.send()
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
	defer dsc.send()

	for {
		select {
		case <-dsc.breaker.Breaked():
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
	dsc.join = append(dsc.join, item)

	if len(dsc.join) < int(dsc.opts.JoinSize) {
		return
	}

	dsc.send()
}

func (dsc *Discipline[Type]) send() {
	defer dsc.resetSendAt()

	join := dsc.prepareJoin()

	if len(join) == 0 {
		return
	}

	select {
	case <-dsc.breaker.Breaked():
		return
	case <-dsc.opts.Ctx.Done():
		return
	case dsc.output <- join:
	}

	if dsc.opts.Released != nil {
		select {
		case <-dsc.breaker.Breaked():
			return
		case <-dsc.opts.Ctx.Done():
			return
		case <-dsc.opts.Released:
		}
	}

	dsc.resetJoin()
}

func (dsc *Discipline[Type]) resetSendAt() {
	dsc.sendAt = time.Now()
}

func (dsc *Discipline[Type]) isTimeouted() bool {
	return time.Since(dsc.sendAt) >= dsc.opts.Timeout
}

func (dsc *Discipline[Type]) resetJoin() {
	dsc.join = dsc.join[:0]
}

func (dsc *Discipline[Type]) copyJoin() []Type {
	return slices.Clone(dsc.join)
}

func (dsc *Discipline[Type]) prepareJoin() []Type {
	if dsc.opts.Released != nil {
		return dsc.join
	}

	return dsc.copyJoin()
}
