package join

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testDiscipline(t *testing.T, useReleased bool) {
	quantity := 105

	input := make(chan uint)
	released := make(chan struct{})

	opts := Opts[uint]{
		Input:    input,
		JoinSize: 10,
	}

	if useReleased {
		opts.Released = released
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	wg := &sync.WaitGroup{}

	wg.Add(2)

	inSequence := make([]uint, 0, quantity)

	go func() {
		defer wg.Done()
		defer close(input)

		for stage := 1; stage <= quantity; stage++ {
			inSequence = append(inSequence, uint(stage))

			input <- uint(stage)
		}
	}()

	outSequence := make([]uint, 0, quantity)

	go func() {
		defer wg.Done()
		defer close(released)

		for slice := range discipline.Output() {
			require.NotEqual(t, 0, slice)

			outSequence = append(outSequence, slice...)

			if useReleased {
				released <- struct{}{}
			}
		}
	}()

	wg.Wait()

	require.Equal(t, inSequence, outSequence)
}

func TestDiscipline(t *testing.T) {
	testDiscipline(t, false)
}

func TestDisciplineReleased(t *testing.T) {
	testDiscipline(t, true)
}

func TestDisciplineOptsValidation(t *testing.T) {
	opts := Opts[uint]{
		JoinSize: 10,
	}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Input: make(chan uint),
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Input:    make(chan uint),
		JoinSize: 10,
		Timeout:  2 * time.Nanosecond,
	}

	_, err = New(opts)
	require.Error(t, err)
}

func TestDisciplineTimeout(t *testing.T) {
	quantity := 105
	pauseAt := quantity / 2

	input := make(chan uint)

	opts := Opts[uint]{
		Input:    input,
		JoinSize: 10,
		Timeout:  500 * time.Millisecond,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	wg := &sync.WaitGroup{}

	wg.Add(2)

	inSequence := make([]uint, 0, quantity)

	go func() {
		defer wg.Done()
		defer close(input)

		for stage := 1; stage <= quantity; stage++ {
			if stage == pauseAt {
				time.Sleep(4 * opts.Timeout)
			}

			inSequence = append(inSequence, uint(stage))

			input <- uint(stage)
		}
	}()

	outSequence := make([]uint, 0, quantity)

	go func() {
		defer wg.Done()

		for slice := range discipline.Output() {
			require.NotEqual(t, 0, slice)

			outSequence = append(outSequence, slice...)
		}
	}()

	wg.Wait()

	require.Equal(t, inSequence, outSequence)
}

func testDisciplineStop(t *testing.T, byCtx bool) {
	quantity := 105
	stopAt := quantity / 2

	input := make(chan uint)
	defer close(input)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := Opts[uint]{
		Ctx:      ctx,
		Input:    input,
		JoinSize: 10,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	wg := &sync.WaitGroup{}

	wg.Add(2)

	inSequence := make([]uint, 0, quantity)

	stop := func() {
		if byCtx {
			cancel()
			return
		}

		discipline.Stop()
	}

	go func() {
		defer wg.Done()

		for stage := 1; stage <= quantity; stage++ {
			if stage == stopAt {
				stop()
				return
			}

			inSequence = append(inSequence, uint(stage))

			input <- uint(stage)
		}
	}()

	outSequence := make([]uint, 0, quantity)

	go func() {
		defer wg.Done()

		for slice := range discipline.Output() {
			require.NotEqual(t, 0, slice)

			outSequence = append(outSequence, slice...)
		}
	}()

	wg.Wait()

	require.GreaterOrEqual(t, len(outSequence), len(inSequence)*80/100)
}

func TestDisciplineStop(t *testing.T) {
	testDisciplineStop(t, false)
}

func TestDisciplineStopByCtx(t *testing.T) {
	testDisciplineStop(t, true)
}

func benchmarkDiscipline(b *testing.B, useReleased bool) {
	quantity := 10000000

	input := make(chan uint)
	released := make(chan struct{})

	opts := Opts[uint]{
		Input: input,

		JoinSize: 100,
	}

	if useReleased {
		opts.Released = released
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	wg := &sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		defer close(input)

		for stage := 1; stage <= quantity; stage++ {
			input <- uint(stage)
		}
	}()

	go func() {
		defer wg.Done()
		defer close(released)

		for range discipline.Output() {
			if useReleased {
				released <- struct{}{}
			}
		}
	}()

	wg.Wait()
}

func BenchmarkDiscipline(b *testing.B) {
	benchmarkDiscipline(b, false)
}

func BenchmarkDisciplineReleased(b *testing.B) {
	benchmarkDiscipline(b, true)
}
