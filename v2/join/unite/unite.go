// Discipline used to accumulate slices elements from an input channel into a one slice
// and write that slice to an output channel when the maximum slice size or timeout for
// its accumulation is reached. It works like a join discipline but accepts slices as
// input and unite their elements into one slice. Moreover, the input slices are not
// divided between the output slices.
package unite

import (
	"errors"
	"slices"
	"time"

	"github.com/akramarenkov/cqos/v2/join/internal/defaults"
)

var (
	ErrInputEmpty   = errors.New("input channel was not specified")
	ErrJoinSizeZero = errors.New("join size is zero")
)

// Options of the created discipline.
type Opts[Type any] struct {
	// Input data channel. For terminate discipline it is necessary and sufficient to
	// close the input channel. Preferably input channel should be buffered for
	// performance reasons. Optimal capacity is in the range of one to three JoinSize
	Input <-chan []Type
	// Maximum size of the output slice. Actual size of the output slice may be
	// smaller due to the timeout or closure of the input channel and the fact
	// that the input slices accumulate entirely. Also, the actual size of the output
	// slice may be larger if an slice larger than the maximum size is received at
	// the input
	JoinSize uint
	// By default, to the output channel is written a copy of the accumulated slice
	// If the NoCopy is set to true, then to the output channel will be directly
	// written the accumulated slice. In this case, after the accumulated slice is
	// no longer used it is necessary to inform the discipline about it by calling
	// Release() method
	NoCopy bool
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
	// this parameter in percents. The lower this value, the lower the performance of
	// the discipline (due to frequent interruptions to check for timeout expiration)
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
		opts.TimeoutInaccuracy = defaults.TimeoutInaccuracy
	}

	return opts
}

// Unite discipline.
type Discipline[Type any] struct {
	opts Opts[Type]

	interruptInterval time.Duration
	join              []Type
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

	interval, err := calcInterruptInterval(opts.Timeout, opts.TimeoutInaccuracy)
	if err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		interruptInterval: interval,
		join:              make([]Type, 0, opts.JoinSize),
		// Value returned by the cap() function is always positive and, in the case of
		// integer overflow due to adding one, the resulting value can only become
		// negative, which will cause a panic when executing make() as same as when
		// specifying a large positive value
		output:  make(chan []Type, 1+cap(opts.Input)),
		release: make(chan struct{}),
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

// Marks accumulated slice as no longer used.
//
// Must be used only if NoCopy option is set to true.
func (dsc *Discipline[Type]) Release() {
	dsc.release <- struct{}{}
}

func (dsc *Discipline[Type]) main() {
	defer close(dsc.output)
	defer close(dsc.release)

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

	for item := range dsc.opts.Input {
		dsc.process(item)
	}
}

func (dsc *Discipline[Type]) process(item []Type) {
	if uint(len(item)) >= dsc.opts.JoinSize {
		dsc.pass()
		dsc.forward(item)

		return
	}

	// Integer overflow is impossible because len() function returns only positive
	// values ​​for the int type and the sum of the two maximum values ​​for the int type is
	// less than the maximum value for the uint type by one
	if uint(len(item))+uint(len(dsc.join)) > dsc.opts.JoinSize {
		dsc.pass()
	}

	dsc.join = append(dsc.join, item...)

	// Integer overflow is impossible because len() function returns only positive
	// values ​​for type int and the maximum value for type int is less than the
	// maximum value for type uint
	if uint(len(dsc.join)) < dsc.opts.JoinSize {
		return
	}

	dsc.pass()
}

func (dsc *Discipline[Type]) pass() {
	if len(dsc.join) == 0 {
		// defer statement is not used to allow inlining of the current function
		dsc.resetPassAt()
		return
	}

	dsc.send(dsc.join)
	dsc.resetJoin()
	dsc.resetPassAt()
}

func (dsc *Discipline[Type]) forward(item []Type) {
	dsc.send(item)
	dsc.resetPassAt()
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

func (dsc *Discipline[Type]) resetJoin() {
	dsc.join = dsc.join[:0]
}

func (dsc *Discipline[Type]) resetPassAt() {
	dsc.passAt = time.Now()
}

func (dsc *Discipline[Type]) isTimeouted() bool {
	return time.Since(dsc.passAt) >= dsc.opts.Timeout
}
