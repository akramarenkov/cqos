// Discipline that used to accumulates elements from the input channel into a slice and
// writes it to the output channel when the size or timeout is reached
package join

import (
	"errors"
	"time"
)

var (
	ErrInputEmpty   = errors.New("input channel was not specified")
	ErrJoinSizeZero = errors.New("join size is zero")
)

const (
	defaultTimeoutInaccuracy = 25
)

// Options of the created discipline
type Opts[Type any] struct {
	// Input data channel. For terminate discipline it is necessary and sufficient to
	// close the input channel
	Input <-chan Type
	// Output slice size
	JoinSize uint
	// By default, to the output channel is written a copy of the accumulated slice
	// If the NoCopy is set to true, then to the output channel will be directly
	// written the accumulated slice
	// In this case, after the accumulated slice is no longer used it is necessary to
	// inform the discipline about it by calling Release()
	NoCopy bool
	// Send timeout of accumulated slice
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
		return ErrInputEmpty
	}

	if opts.JoinSize == 0 {
		return ErrJoinSizeZero
	}

	return nil
}

func (opts Opts[Type]) normalize() Opts[Type] {
	if opts.TimeoutInaccuracy == 0 {
		opts.TimeoutInaccuracy = defaultTimeoutInaccuracy
	}

	return opts
}

// Join discipline
type Discipline[Type any] struct {
	opts Opts[Type]

	interruptInterval time.Duration
	join              []Type
	output            chan []Type
	release           chan struct{}
	sendAt            time.Time
}

// Creates and runs discipline
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	opts = opts.normalize()

	interval, err := calcInterruptIntervalZeroAllowed(
		opts.Timeout,
		opts.TimeoutInaccuracy,
	)
	if err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		interruptInterval: interval,
		join:              make([]Type, 0, opts.JoinSize),
		output:            make(chan []Type, 1),
		release:           make(chan struct{}),
	}

	dsc.resetSendAt()

	go dsc.main()

	return dsc, nil
}

// Returns output channel.
//
// If this channel is closed, it means that the discipline is terminated
func (dsc *Discipline[Type]) Output() <-chan []Type {
	return dsc.output
}

// Marks accumulated slice as no longer used.
//
// Must be used only if NoCopy option is set to true
func (dsc *Discipline[Type]) Release() {
	dsc.release <- struct{}{}
}

func (dsc *Discipline[Type]) main() {
	defer close(dsc.release)
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

	for item := range dsc.opts.Input {
		dsc.process(item)
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

	dsc.output <- join

	if dsc.opts.NoCopy {
		<-dsc.release
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
	sent := make([]Type, len(dsc.join))

	copy(sent, dsc.join)

	return sent
}

func (dsc *Discipline[Type]) prepareJoin() []Type {
	if dsc.opts.NoCopy {
		return dsc.join
	}

	return dsc.copyJoin()
}
