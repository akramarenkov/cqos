package join

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/breaker"
	"github.com/akramarenkov/cqos/internal/consts"

	"github.com/stretchr/testify/require"
)

const (
	defaultTestTimeout = (consts.OneHundredPercent *
		consts.ReliablyMeasurableDuration) / defaultTimeoutInaccuracy
)

func TestOptsValidation(t *testing.T) {
	opts := Opts[int]{
		JoinSize: 10,
	}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input: make(chan int),
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input:    make(chan int),
		JoinSize: 10,
		Timeout:  10 * time.Millisecond,
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input:    make(chan int),
		JoinSize: 10,
	}

	_, err = New(opts)
	require.NoError(t, err)

	opts = Opts[int]{
		Input:    make(chan int),
		JoinSize: 10,
		Timeout:  100 * time.Millisecond,
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func TestDiscipline(t *testing.T) {
	testDiscipline(t, false, defaultTestTimeout)
}

func TestDisciplineReleased(t *testing.T) {
	testDiscipline(t, true, defaultTestTimeout)
}

func TestDisciplineUntimeouted(t *testing.T) {
	testDiscipline(t, false, 0)
}

func testDiscipline(t *testing.T, useReleased bool, timeout time.Duration) {
	quantity := 105

	input := make(chan int)

	released := make(chan struct{})
	defer close(released)

	opts := Opts[int]{
		Input:    input,
		JoinSize: 10,
		Timeout:  timeout,
	}

	if useReleased {
		opts.Released = released
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	joins := 0

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	go func() {
		defer close(input)

		for stage := 1; stage <= quantity; stage++ {
			inSequence = append(inSequence, stage)

			input <- stage
		}
	}()

	for slice := range discipline.Output() {
		require.NotEqual(t, 0, slice)

		joins++

		outSequence = append(outSequence, slice...)

		if useReleased {
			released <- struct{}{}
		}
	}

	expectedJoins := int(math.Ceil(float64(quantity) / float64(opts.JoinSize)))

	require.Equal(t, inSequence, outSequence)
	require.Equal(t, expectedJoins, joins)
}

func TestDisciplineTimeout(t *testing.T) {
	quantity := 105
	pauseAt := 52

	input := make(chan int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: 10,
		Timeout:  100 * time.Millisecond,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	joins := 0

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	go func() {
		defer close(input)

		for stage := 1; stage <= quantity; stage++ {
			if stage == pauseAt {
				time.Sleep(5 * opts.Timeout)
			}

			inSequence = append(inSequence, stage)

			input <- stage
		}
	}()

	for slice := range discipline.Output() {
		require.NotEqual(t, 0, slice)

		joins++

		outSequence = append(outSequence, slice...)
	}

	// plus one slice with incomplete size due to pause on write to input
	expectedJoins := int(math.Ceil(float64(quantity)/float64(opts.JoinSize))) + 1

	require.Equal(t, inSequence, outSequence)
	require.Equal(t, expectedJoins, joins)
}

func TestDisciplineStop(t *testing.T) {
	testDisciplineStop(t, false, false, defaultTestTimeout)
	testDisciplineStop(t, false, false, 0)
	testDisciplineStop(t, false, true, 0)
}

func TestDisciplineStopByCtx(t *testing.T) {
	testDisciplineStop(t, true, false, defaultTestTimeout)
	testDisciplineStop(t, true, false, 0)
	testDisciplineStop(t, true, true, 0)
}

func testDisciplineStop(
	t *testing.T,
	byCtx bool,
	useReleased bool,
	timeout time.Duration,
) {
	quantity := 105
	stopAt := 52

	input := make(chan int)

	released := make(chan struct{})
	defer close(released)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := Opts[int]{
		Ctx:      ctx,
		Input:    input,
		JoinSize: 10,
		Timeout:  timeout,
	}

	if useReleased {
		opts.Released = released
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	stop := func() {
		if byCtx {
			cancel()
			return
		}

		discipline.Stop()
	}

	interrupter := breaker.NewClosing()

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(input)

		for stage := 1; stage <= quantity; stage++ {
			if stage == stopAt {
				stop()
				return
			}

			inSequence = append(inSequence, stage)

			select {
			case input <- stage:
			case <-interrupter.Closed():
				return
			}
		}
	}()

	for slice := range discipline.Output() {
		require.NotEqual(t, 0, slice)

		outSequence = append(outSequence, slice...)

		if useReleased {
			stop()
			interrupter.Close()
		}
	}

	wg.Wait()

	require.GreaterOrEqual(t, len(outSequence), len(inSequence)*80/100)
}

func BenchmarkDiscipline(b *testing.B) {
	benchmarkDiscipline(b, false, defaultTestTimeout)
}

func BenchmarkDisciplineReleased(b *testing.B) {
	benchmarkDiscipline(b, true, defaultTestTimeout)
}

func BenchmarkDisciplineUntimeouted(b *testing.B) {
	benchmarkDiscipline(b, false, 0)
}

func benchmarkDiscipline(b *testing.B, useReleased bool, timeout time.Duration) {
	quantity := 10000000

	input := make(chan int)

	released := make(chan struct{})
	defer close(released)

	opts := Opts[int]{
		Input:    input,
		JoinSize: 100,
		Timeout:  timeout,
	}

	if useReleased {
		opts.Released = released
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	go func() {
		defer close(input)

		for stage := 1; stage <= quantity; stage++ {
			input <- stage
		}
	}()

	for range discipline.Output() {
		if useReleased {
			released <- struct{}{}
		}
	}
}
