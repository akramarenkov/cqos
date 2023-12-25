package join

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
}

func testDiscipline(t *testing.T, noCopy bool) {
	quantity := 105

	input := make(chan int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: 10,
		NoCopy:   noCopy,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	inSequence := make([]int, 0, quantity)

	go func() {
		defer close(input)

		for stage := 1; stage <= quantity; stage++ {
			inSequence = append(inSequence, stage)

			input <- stage
		}
	}()

	joins := 0

	outSequence := make([]int, 0, quantity)

	for slice := range discipline.Output() {
		require.NotEqual(t, 0, slice)

		joins++

		outSequence = append(outSequence, slice...)

		if noCopy {
			discipline.Release()
		}
	}

	expectedJoins := int(math.Ceil(float64(quantity) / float64(opts.JoinSize)))

	require.Equal(t, inSequence, outSequence)
	require.Equal(t, expectedJoins, joins)
}

func TestDiscipline(t *testing.T) {
	testDiscipline(t, false)
}

func TestDisciplineNoCopy(t *testing.T) {
	testDiscipline(t, true)
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

	inSequence := make([]int, 0, quantity)

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

	joins := 0

	outSequence := make([]int, 0, quantity)

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

func benchmarkDiscipline(b *testing.B, noCopy bool) {
	quantity := 10000000

	input := make(chan int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: 100,
		NoCopy:   noCopy,
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
		if noCopy {
			discipline.Release()
		}
	}
}

func BenchmarkDiscipline(b *testing.B) {
	benchmarkDiscipline(b, false)
}

func BenchmarkDisciplineReleased(b *testing.B) {
	benchmarkDiscipline(b, true)
}
