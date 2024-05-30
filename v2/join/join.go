// Discipline that used to accumulates elements from the input channel into a slice and
// writes it to the output channel when the size or timeout is reached.
package join

import (
	"errors"
	"slices"
	"time"

	"github.com/akramarenkov/cqos/v2/join/internal/spinner"
)

var (
	ErrInputEmpty   = errors.New("input channel was not specified")
	ErrJoinSizeZero = errors.New("join size is zero")
)

const (
	// The number 2 in indicates the number of goroutines involved in processing
	// slices. In this case, there are 2 of them - the accumulation and send
	// goroutines.
	involvedInProcessing = 2

	// Defines buffers quantity. Cannot be less than number of goroutines involved
	// in processing slices, but there may be more than this value, although this
	// does not make sense.
	buffersQuantity = involvedInProcessing

	// To prevent from using still unsent slices by the accumulation goroutine,
	// it must be blocked from writing to the interim channel on the last one of
	// the unsent slices.
	interimCapacity = buffersQuantity - involvedInProcessing
)

// Options of the created discipline.
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

// Join discipline.
type Discipline[Type any] struct {
	opts Opts[Type]

	id                *spinner.Spinner
	interim           chan int
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

		id:                spinner.New(0, buffersQuantity-1),
		interim:           make(chan int, interimCapacity),
		interruptInterval: interval,
		joins:             make([][]Type, buffersQuantity),
		output:            make(chan []Type, 1),
		release:           make(chan struct{}),
	}

	dsc.initJoins()
	dsc.resetPassAt()

	go dsc.main()

	return dsc, nil
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

	complete := make(chan struct{})
	defer func() { <-complete }()

	defer close(dsc.interim)

	go dsc.sender(complete)

	if dsc.interruptInterval == 0 {
		dsc.loopUntimeouted()
		return
	}

	dsc.loop()
}

func (dsc *Discipline[Type]) sender(complete chan struct{}) {
	defer close(complete)

	for id := range dsc.interim {
		dsc.send(id)
	}
}

func (dsc *Discipline[Type]) send(id int) {
	join := dsc.prepareJoin(id)

	if len(join) == 0 {
		return
	}

	dsc.output <- join

	if dsc.opts.NoCopy {
		<-dsc.release
	}

	dsc.resetJoin(id)
}

func (dsc *Discipline[Type]) prepareJoin(id int) []Type {
	if dsc.opts.NoCopy {
		return dsc.joins[id]
	}

	return dsc.copyJoin(id)
}

func (dsc *Discipline[Type]) copyJoin(id int) []Type {
	return slices.Clone(dsc.joins[id])
}

func (dsc *Discipline[Type]) resetJoin(id int) {
	dsc.joins[id] = dsc.joins[id][:0]
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

func (dsc *Discipline[Type]) process(item Type) {
	id := dsc.id.Actual()

	dsc.joins[id] = append(dsc.joins[id], item)

	if len(dsc.joins[id]) < int(dsc.opts.JoinSize) {
		return
	}

	dsc.pass()
}

func (dsc *Discipline[Type]) pass() {
	dsc.interim <- dsc.id.Actual()
	dsc.id.Spin()
	dsc.resetPassAt()
}

func (dsc *Discipline[Type]) resetPassAt() {
	dsc.passAt = time.Now()
}

func (dsc *Discipline[Type]) isTimeouted() bool {
	return time.Since(dsc.passAt) >= dsc.opts.Timeout
}
