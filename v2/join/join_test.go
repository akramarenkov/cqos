package join

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/general"
	"github.com/akramarenkov/cqos/v2/join/internal/common"
	"github.com/akramarenkov/cqos/v2/join/internal/inspect"

	"github.com/akramarenkov/stressor"
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
		Timeout:  general.ReliablyMeasurableDuration,
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input:    make(chan int),
		JoinSize: 10,
		Timeout:  common.DefaultMinTimeout,
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
	for quantity := 100; quantity <= 125; quantity++ {
		testDiscipline(t, quantity, 10, false, common.DefaultTestTimeout)
	}
}

func TestDisciplineNoCopy(t *testing.T) {
	for quantity := 100; quantity <= 125; quantity++ {
		testDiscipline(t, quantity, 10, true, common.DefaultTestTimeout)
	}
}

func TestDisciplineUntimeouted(t *testing.T) {
	for quantity := 100; quantity <= 125; quantity++ {
		testDiscipline(t, quantity, 10, false, 0)
	}
}

func testDiscipline(
	t *testing.T,
	quantity int,
	joinSize uint,
	noCopy bool,
	timeout time.Duration,
) {
	input := make(chan int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	expectedOutput := inspect.Expected(quantity, 1, opts.JoinSize)
	output := make([][]int, 0, len(expectedOutput))

	go func() {
		defer close(input)

		for _, slice := range inspect.Input(quantity, 1) {
			for _, item := range slice {
				inSequence = append(inSequence, item)

				input <- item
			}
		}
	}()

	for slice := range discipline.Output() {
		require.NotEmpty(t, slice)

		outSequence = append(outSequence, slice...)

		output = append(output, append([]int(nil), slice...))

		if noCopy {
			discipline.Release()
		}
	}

	require.Equal(t, inSequence, outSequence,
		"quantity: %v, join size: %v, no copy: %v, timeout: %v",
		quantity,
		joinSize,
		noCopy,
		timeout,
	)
	require.Equal(t, expectedOutput, output,
		"quantity: %v, join size: %v, no copy: %v, timeout: %v",
		quantity,
		joinSize,
		noCopy,
		timeout,
	)
}

func TestDisciplineTimeout(t *testing.T) {
	for quantity := 100; quantity <= 125; quantity++ {
		testDisciplineTimeout(t, quantity, 10, 53)
	}
}

func testDisciplineTimeout(
	t *testing.T,
	quantity int,
	joinSize uint,
	pauseAt int,
) {
	input := make(chan int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		Timeout:  100 * time.Millisecond,
	}

	pauseAt = inspect.PickUpPauseAt(quantity, pauseAt, 1, opts.JoinSize)
	require.NotEqual(t, 0, pauseAt)

	pausetAtDuration := inspect.CalcPauseAtDuration(opts.Timeout)

	discipline, err := New(opts)
	require.NoError(t, err)

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	expectedOutput := inspect.ExpectedWithTimeout(
		quantity,
		pauseAt,
		1,
		opts.JoinSize,
	)

	output := make([][]int, 0, len(expectedOutput))

	go func() {
		defer close(input)

		for _, slice := range inspect.Input(quantity, 1) {
			for _, item := range slice {
				if item == pauseAt {
					time.Sleep(pausetAtDuration)
				}

				inSequence = append(inSequence, item)

				input <- item
			}
		}
	}()

	for slice := range discipline.Output() {
		require.NotEmpty(t, slice)

		outSequence = append(outSequence, slice...)

		output = append(output, append([]int(nil), slice...))
	}

	require.Equal(t, inSequence, outSequence,
		"quantity: %v, join size: %v, pause at: %v",
		quantity,
		joinSize,
		pauseAt,
	)
	require.Equal(t, expectedOutput, output,
		"quantity: %v, join size: %v, pause at: %v",
		quantity,
		joinSize,
		pauseAt,
	)
}

func BenchmarkDiscipline(b *testing.B) {
	benchmarkDiscipline(b, false, common.DefaultTestTimeout, 0, false)
}

func BenchmarkDisciplineNoCopy(b *testing.B) {
	benchmarkDiscipline(b, true, common.DefaultTestTimeout, 0, false)
}

func BenchmarkDisciplineUntimeouted(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0, false)
}

func benchmarkDiscipline(
	b *testing.B,
	noCopy bool,
	timeout time.Duration,
	outputDelayFactor float64,
	stressSystem bool,
) {
	joinsQuantity := b.N
	joinSize := uint(10)
	// Accuracy of this delay is sufficient for the benchmark and
	// general.ReliablyMeasurableDuration is too large to perform a representative
	// number of iterations
	inputDelayBase := 1 * time.Millisecond

	quantity := joinsQuantity * int(joinSize)

	inputDelay, outputDelay := calcBenchmarkDelays(
		inputDelayBase,
		joinSize,
		outputDelayFactor,
	)

	input := make(chan int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	if stressSystem {
		stress := stressor.New(stressor.Opts{})
		defer stress.Stop()

		time.Sleep(time.Second)
	}

	b.ResetTimer()

	go func() {
		defer close(input)

		for item := 1; item <= quantity; item++ {
			time.Sleep(inputDelay)

			input <- item
		}
	}()

	for range discipline.Output() {
		time.Sleep(outputDelay)

		if noCopy {
			discipline.Release()
		}
	}
}

func calcBenchmarkDelays(
	inputDelay time.Duration,
	joinSize uint,
	outputDelayFactor float64,
) (time.Duration, time.Duration) {
	outputDelay := time.Duration(
		outputDelayFactor *
			float64(inputDelay) *
			float64(joinSize),
	)

	if outputDelay == 0 {
		inputDelay = 0
	}

	return inputDelay, outputDelay
}
