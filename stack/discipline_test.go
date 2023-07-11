package stack

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDiscipline(t *testing.T) {
	quantity := 105

	input := make(chan uint)

	opts := Opts[uint]{
		Input: input,
		Size:  10,
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

		for stack := range discipline.Output() {
			require.NotEqual(t, 0, stack)

			outSequence = append(outSequence, stack...)
		}
	}()

	wg.Wait()

	require.Equal(t, inSequence, outSequence)
}

func TestDisciplineReleased(t *testing.T) {
	quantity := 105

	input := make(chan uint)
	released := make(chan struct{})

	opts := Opts[uint]{
		Input:    input,
		Released: released,
		Size:     10,
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

		for stack := range discipline.Output() {
			require.NotEqual(t, 0, stack)

			outSequence = append(outSequence, stack...)

			released <- struct{}{}
		}
	}()

	wg.Wait()

	require.Equal(t, inSequence, outSequence)
}

func TestDisciplineOptsValidation(t *testing.T) {
	opts := Opts[uint]{
		Size: 10,
	}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Input: make(chan uint),
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Input:   make(chan uint),
		Size:    10,
		Timeout: 2 * time.Nanosecond,
	}

	_, err = New(opts)
	require.Error(t, err)
}

func TestDisciplineTimeout(t *testing.T) {
	quantity := 105
	pauseAt := quantity / 2

	input := make(chan uint)

	opts := Opts[uint]{
		Input:   input,
		Size:    10,
		Timeout: 500 * time.Millisecond,
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

		for stack := range discipline.Output() {
			require.NotEqual(t, 0, stack)
			outSequence = append(outSequence, stack...)
		}
	}()

	wg.Wait()

	require.Equal(t, inSequence, outSequence)
}

func TestDisciplineStop(t *testing.T) {
	quantity := 105
	stopAt := quantity / 2

	input := make(chan uint)

	opts := Opts[uint]{
		Input: input,
		Size:  10,
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
			if stage == stopAt {
				discipline.Stop()
				return
			}

			inSequence = append(inSequence, uint(stage))

			input <- uint(stage)
		}
	}()

	outSequence := make([]uint, 0, quantity)

	go func() {
		defer wg.Done()

		for stack := range discipline.Output() {
			require.NotEqual(t, 0, stack)
			outSequence = append(outSequence, stack...)
		}
	}()

	wg.Wait()

	require.GreaterOrEqual(t, len(outSequence), len(inSequence)*80/100)
}

func TestDisciplineCtx(t *testing.T) {
	quantity := 105
	stopAt := quantity / 2

	input := make(chan uint)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := Opts[uint]{
		Ctx:   ctx,
		Input: input,
		Size:  10,
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
			if stage == stopAt {
				cancel()
				return
			}

			inSequence = append(inSequence, uint(stage))

			input <- uint(stage)
		}
	}()

	outSequence := make([]uint, 0, quantity)

	go func() {
		defer wg.Done()

		for stack := range discipline.Output() {
			require.NotEqual(t, 0, stack)
			outSequence = append(outSequence, stack...)
		}
	}()

	wg.Wait()

	require.GreaterOrEqual(t, len(outSequence), len(inSequence)*80/100)
}

func BenchmarkDiscipline(b *testing.B) {
	quantity := 10000000

	input := make(chan uint)

	opts := Opts[uint]{
		Input: input,
		Size:  100,
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

	count := 0

	go func() {
		defer wg.Done()

		for range discipline.Output() {
			count++
		}
	}()

	wg.Wait()
}

func BenchmarkDisciplineReleased(b *testing.B) {
	quantity := 10000000

	input := make(chan uint)
	released := make(chan struct{})

	opts := Opts[uint]{
		Input:    input,
		Released: released,
		Size:     100,
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
			released <- struct{}{}
		}
	}()

	wg.Wait()
}
