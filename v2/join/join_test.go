package join

import (
	"math"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"

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
		Timeout:  consts.ReliablyMeasurableDuration,
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input:    make(chan int),
		JoinSize: 10,
		Timeout:  minDefaultTimeout,
	}

	_, err = New(opts)
	require.NoError(t, err)

	opts = Opts[int]{
		Input:    make(chan int),
		JoinSize: 10,
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func TestDiscipline(t *testing.T) {
	testDiscipline(t, false, minDefaultTimeout)
}

func TestDisciplineNoCopy(t *testing.T) {
	testDiscipline(t, true, minDefaultTimeout)
}

func TestDisciplineUntimeouted(t *testing.T) {
	testDiscipline(t, false, 0)
}

func testDiscipline(t *testing.T, noCopy bool, timeout time.Duration) {
	quantity := 105

	input := make(chan int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: 10,
		NoCopy:   noCopy,
		Timeout:  timeout,
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

		if noCopy {
			discipline.Release()
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

func BenchmarkDiscipline(b *testing.B) {
	benchmarkDiscipline(b, false, minDefaultTimeout, 0)
}

func BenchmarkDisciplineRelease(b *testing.B) {
	benchmarkDiscipline(b, true, minDefaultTimeout, 0)
}

func BenchmarkDisciplineUntimeouted(b *testing.B) {
	benchmarkDiscipline(b, false, 0, 0)
}

func BenchmarkDisciplineOutputDelayIsSame(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 1)
}

func BenchmarkDisciplineOutputDelayIsLess(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0.25)
}

func BenchmarkDisciplineOutputDelayIsLonger(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 4)
}

func benchmarkDiscipline(
	b *testing.B,
	noCopy bool,
	timeout time.Duration,
	outputDelayFactor float64,
) {
	quantity, joinSize := benchmarkCalcQuantityJoinSize(b)
	inputDelay, outputDelay := benchmarkCalcDelays(joinSize, outputDelayFactor)

	input := make(chan int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	go func() {
		defer close(input)

		for stage := 1; stage <= quantity; stage++ {
			time.Sleep(inputDelay)

			input <- stage
		}
	}()

	for range discipline.Output() {
		time.Sleep(outputDelay)

		if noCopy {
			discipline.Release()
		}
	}
}

func benchmarkCalcQuantityJoinSize(b *testing.B) (int, uint) {
	const joinsQuantity = 100

	quantity := b.N
	joinSize := quantity / joinsQuantity

	if joinSize == 0 {
		joinSize = 1
	}

	return quantity, uint(joinSize)
}

func benchmarkCalcDelays(
	joinSize uint,
	outputDelayFactor float64,
) (time.Duration, time.Duration) {
	inputDelay := 1 * time.Millisecond
	outputDelay := time.Duration(outputDelayFactor * float64(inputDelay) * float64(joinSize))

	if outputDelay == 0 {
		inputDelay = 0
	}

	return inputDelay, outputDelay
}
