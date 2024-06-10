// Discipline that is used to accumulates elements from the input channel into a
// slice and writes it to the output channel when the maximum size or timeout is
// reached.
package join

import (
	"errors"
	"slices"
	"time"

	"github.com/akramarenkov/cqos/v2/join/internal/common"
	"github.com/akramarenkov/cqos/v2/join/internal/spinner"
)

var (
	ErrInputEmpty   = errors.New("input channel was not specified")
	ErrJoinSizeZero = errors.New("join size is zero")
)

// Options of the created discipline.
type Opts[Type any] struct {
	// Input data channel. For terminate discipline it is necessary and sufficient to
	// close the input channel
	Input <-chan Type
	// Maximum size of the output slice. Actual size of the output slice may be
	// smaller due to the timeout or closure of the input channel
	JoinSize uint
	// By default, to the output channel is written a copy of the accumulated slice
	// If the NoCopy is set to true, then to the output channel will be directly
	// written the accumulated slice. In this case, after the accumulated slice is
	// no longer used it is necessary to inform the discipline about it by calling
	// Release() method
	NoCopy bool
	// Disables double buffering. Double buffering is implemented through two
	// accumulation buffers, an additional goroutine of sending the accumulated
	// buffer and an intermediate channel between the goroutines of accumulation and
	// sending. Double buffering allows you to reduce processing time (increase
	// performance) by 15-80% if the data accumulation and sending times are
	// comparable, i.e. they correspond in the range from about 1:4 to 4:1 and these
	// times are longer than the time spent on transferring the buffer between the
	// accumulation and sending goroutines. Disabling double buffering can reduce
	// processing time (increase performance) in other cases by 5-15% at low
	// system load, and up to 70% at high system load
	NoDoubleBuffering bool
	// Timeout for slice accumulation. If the slice has not been filled completely
	// in the allotted time, the data accumulated during this time is written to
	// the output channel. A zero or negative value means that discipline will wait
	// for the missing data until they appear or the channel is closed (in this case,
	// the data will be accumulated data will be written to the output channel)
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
		opts.TimeoutInaccuracy = common.DefaultTimeoutInaccuracy
	}

	return opts
}

// Join discipline.
type Discipline[Type any] struct {
	opts Opts[Type]

	id                *spinner.Spinner
	interim           chan []Type
	interruptInterval time.Duration
	joins             [][]Type
	output            chan []Type
	passAt            time.Time
	release           chan struct{}
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

		id:                prepareSpinner(opts),
		interim:           make(chan []Type, common.InterimCapacity),
		interruptInterval: interval,
		joins:             make([][]Type, common.BuffersQuantity),
		output:            make(chan []Type, 1),
		release:           make(chan struct{}),
	}

	dsc.initJoins()
	dsc.resetPassAt()

	go dsc.main()

	return dsc, nil
}

func prepareSpinner[Type any](opts Opts[Type]) *spinner.Spinner {
	if opts.NoDoubleBuffering {
		return spinner.New(0, 0)
	}

	return spinner.New(0, common.BuffersQuantity-1)
}

func (dsc *Discipline[Type]) initJoins() {
	for id := range dsc.joins {
		dsc.joins[id] = make([]Type, 0, dsc.opts.JoinSize)
	}
}

// Returns output channel.
//
// If this channel is closed, it means that the discipline is terminated.
func (dsc *Discipline[Type]) Output() <-chan []Type {
	return dsc.output
}

// Marks accumulated slice as no longer used.
//
// Must be used only if NoCopy option is set to true.
func (dsc *Discipline[Type]) Release() {
	dsc.release <- struct{}{}
}

func (dsc *Discipline[Type]) main() {
	defer close(dsc.output)
	defer close(dsc.release)

	closing := dsc.runSender()
	defer closing()

	if dsc.interruptInterval == 0 {
		dsc.accumulatorUntimeouted()
		return
	}

	dsc.accumulator()
}

func (dsc *Discipline[Type]) runSender() func() {
	if dsc.opts.NoDoubleBuffering {
		return func() {}
	}

	complete := make(chan struct{})

	closing := func() {
		close(dsc.interim)
		<-complete
	}

	go dsc.sender(complete)

	return closing
}

func (dsc *Discipline[Type]) sender(complete chan struct{}) {
	defer close(complete)

	for item := range dsc.interim {
		dsc.send(item)
	}
}

func (dsc *Discipline[Type]) send(item []Type) {
	item = dsc.prepareItem(item)

	dsc.output <- item

	if dsc.opts.NoCopy {
		<-dsc.release
	}
}

func (dsc *Discipline[Type]) prepareItem(item []Type) []Type {
	if dsc.opts.NoCopy {
		return item
	}

	return slices.Clone(item)
}

func (dsc *Discipline[Type]) accumulator() {
	defer dsc.passActual()

	ticker := time.NewTicker(dsc.interruptInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if dsc.isTimeouted() {
				dsc.passActual()
			}
		case item, opened := <-dsc.opts.Input:
			if !opened {
				return
			}

			dsc.process(item)
		}
	}
}

func (dsc *Discipline[Type]) accumulatorUntimeouted() {
	defer dsc.passActual()

	for item := range dsc.opts.Input {
		dsc.process(item)
	}
}

func (dsc *Discipline[Type]) process(item Type) {
	id := dsc.id.Actual()

	dsc.joins[id] = append(dsc.joins[id], item)

	if len(dsc.joins[id]) < int(dsc.opts.JoinSize) {
		return
	}

	dsc.passActual()
}

func (dsc *Discipline[Type]) passActual() {
	defer dsc.resetPassAt()

	id := dsc.id.Actual()

	if len(dsc.joins[id]) == 0 {
		return
	}

	if dsc.opts.NoDoubleBuffering {
		dsc.send(dsc.joins[id])
	} else {
		dsc.interim <- dsc.joins[id]
	}

	dsc.resetActual()
	dsc.id.Spin()
}

func (dsc *Discipline[Type]) resetActual() {
	id := dsc.id.Actual()

	dsc.joins[id] = dsc.joins[id][:0]
}

func (dsc *Discipline[Type]) resetPassAt() {
	dsc.passAt = time.Now()
}

func (dsc *Discipline[Type]) isTimeouted() bool {
	return time.Since(dsc.passAt) >= dsc.opts.Timeout
}
