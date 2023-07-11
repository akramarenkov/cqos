package join

import (
	"context"
	"errors"
	"time"

	"github.com/akramarenkov/cqos/breaker"
)

var (
	ErrEmptyInput      = errors.New("input channel was not specified")
	ErrInvalidJoinSize = errors.New("invalid join size")
	ErrTimeoutTooSmall = errors.New("timeout value is too small")
)

const (
	defaultTimeout        = 1 * time.Millisecond
	defaultTimeoutDivider = 4
)

// Options of the created discipline
type Opts[Type any] struct {
	// Roughly terminates (cancels) work of the discipline
	Ctx context.Context
	// Input data channel. For graceful termination it is enough to
	// close the input channel
	Input <-chan Type
	// Output slice size
	JoinSize uint
	// By default, to the output channel is written a copy of the accumulated slice
	// If the Released channel is set, then to the output channel will be directly
	// written the accumulated slice
	// In this case, after the accumulated slice is used it is necessary to inform
	// the discipline about it by writing to Released channel
	Released <-chan struct{}
	// Send timeout of accumulated slice
	Timeout time.Duration
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

	if opts.Timeout == 0 {
		opts.Timeout = defaultTimeout
	}

	return opts
}

// Main discipline
type Discipline[Type any] struct {
	opts Opts[Type]

	breaker *breaker.Breaker

	output    chan []Type
	sendAt    time.Time
	join      []Type
	timeouter *time.Ticker
}

// Creates and runs main discipline
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	opts = opts.normalize()

	duration, err := calcTimeouterDuration(opts.Timeout)
	if err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		breaker: breaker.New(),

		output:    make(chan []Type, 1),
		join:      make([]Type, 0, opts.JoinSize),
		timeouter: time.NewTicker(duration),
	}

	go dsc.loop()

	return dsc, nil
}

// Maximum timeout error is calculated as timeout + timeout/divider.
//
// Relative timeout error in percent is calculated as 100/divider.
//
// Remember that a ticker is 'divider' times more likely to be triggered
func calcTimeouterDuration(timeout time.Duration) (time.Duration, error) {
	timeout /= defaultTimeoutDivider

	if timeout == 0 {
		return 0, ErrTimeoutTooSmall
	}

	return timeout, nil
}

// Returns output channel
func (dsc *Discipline[Type]) Output() <-chan []Type {
	return dsc.output
}

// Roughly terminates work of the discipline.
//
// Use for wait completion at terminates via context
func (dsc *Discipline[Type]) Stop() {
	dsc.breaker.Break()
}

func (dsc *Discipline[Type]) loop() {
	defer dsc.breaker.Complete()
	defer dsc.timeouter.Stop()
	defer close(dsc.output)

	dsc.resetSendAt()

	for {
		select {
		case <-dsc.breaker.Breaked():
			return
		case <-dsc.opts.Ctx.Done():
			return
		case <-dsc.timeouter.C:
			if dsc.isTimeouted() {
				dsc.send()
			}
		case item, opened := <-dsc.opts.Input:
			if !opened {
				dsc.send()
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
	dsc.resetSendAt()
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
	sent := make([]Type, len(dsc.join))

	copy(sent, dsc.join)

	return sent
}

func (dsc *Discipline[Type]) prepareJoin() []Type {
	if dsc.opts.Released != nil {
		return dsc.join
	}

	return dsc.copyJoin()
}
