// Discipline that is used to limits the speed of passing data elements from the
// input channel to the output channel.
package limit

import (
	"errors"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
)

var (
	ErrInputEmpty = errors.New("input channel was not specified")
)

// Options of the created discipline.
type Opts[Type any] struct {
	// Input data channel. For terminate discipline it is necessary and sufficient to
	// close the input channel. Preferably input channel should be buffered for
	// performance reasons. Optimal capacity is in the range of 1e2 to 1e6 and
	// should be determined using benchmarks
	//
	// Note that if the number of data elements written to the input channel before it
	// is closed is a multiple of the Quantity field in the rate limit structure, the
	// discipline will still perform a delay after the last data element is transmitted.
	// This, with large values ​​of the Interval field in the rate limit structure, will
	// result in a long discipline completion time
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

// Limit discipline.
type Discipline[Type any] struct {
	opts Opts[Type]

	output chan Type
}

// Creates and runs discipline.
func New[Type any](opts Opts[Type]) (*Discipline[Type], error) {
	if err := opts.isValid(); err != nil {
		return nil, err
	}

	dsc := &Discipline[Type]{
		opts: opts,

		output: make(chan Type, 1+cap(opts.Input)),
	}

	go dsc.main()

	return dsc, nil
}

// Returns output channel.
//
// If this channel is closed, it means that the discipline is terminated.
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
		duration, stop := dsc.transfer()
		if stop {
			return
		}

		delay = dsc.delay(delay, duration)
	}
}

func (dsc *Discipline[Type]) transfer() (time.Duration, bool) {
	startedAt := time.Now()

	if stop := dsc.pass(); stop {
		return 0, true
	}

	// This duration is the time difference of monotonic clock, so it is always
	// at least non-negative
	return time.Since(startedAt), false
}

func (dsc *Discipline[Type]) pass() bool {
	for range dsc.opts.Limit.Quantity {
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

func (dsc *Discipline[Type]) delay(delay, duration time.Duration) time.Duration {
	// Integer overflow is impossible, because values of Interval field in rate limit
	// structure and transfer duration are greater than zero
	remainder := dsc.opts.Limit.Interval - duration

	// Integer overflow is impossible in practice since it can only occur with large
	// values ​​of Interval field in rate limit structure and transfer duration that
	// cannot be tested
	delay += remainder

	// Reset to zero is performed to prevent situations where the rate limit is not
	// present after the data transfer duration has been exceeding the value of
	// Interval field in rate limit structure for a long time (for example, due to an
	// external load on the system) and the subsequent reduction of the data transfer
	// duration compared to the value of Interval field in rate limit structure
	// (after the external load on the system has decreased)
	//
	// Integer overflow when negating is impossible, because value of Interval field
	// in rate limit structure are greater than zero i.e. less than the minimum
	// negative value for given type in absolute value
	if delay < -dsc.opts.Limit.Interval {
		return 0
	}

	// Delay accumulates and is deferred until it becomes sufficient to execute with
	// acceptable accuracy
	if delay < consts.ReliablyMeasurableDuration {
		return delay
	}

	time.Sleep(delay)

	return 0
}
