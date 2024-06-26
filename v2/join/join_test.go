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
		testDiscipline(t, quantity, 10, false, common.DefaultTestTimeout, false)
		testDiscipline(t, quantity, 10, false, common.DefaultTestTimeout, true)
	}
}

func TestDisciplineNoCopy(t *testing.T) {
	for quantity := 100; quantity <= 125; quantity++ {
		testDiscipline(t, quantity, 10, true, common.DefaultTestTimeout, false)
		testDiscipline(t, quantity, 10, true, common.DefaultTestTimeout, true)
	}
}

func TestDisciplineUntimeouted(t *testing.T) {
	for quantity := 100; quantity <= 125; quantity++ {
		testDiscipline(t, quantity, 10, false, 0, false)
		testDiscipline(t, quantity, 10, false, 0, true)
	}
}

func testDiscipline(
	t *testing.T,
	quantity int,
	joinSize uint,
	noCopy bool,
	timeout time.Duration,
	noDoubleBuffering bool,
) {
	input := make(chan int)

	opts := Opts[int]{
		Input:             input,
		JoinSize:          joinSize,
		NoCopy:            noCopy,
		NoDoubleBuffering: noDoubleBuffering,
		Timeout:           timeout,
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
		"quantity: %v, join size: %v, no copy: %v, "+
			"timeout: %v, no double buffering: %v",
		quantity,
		joinSize,
		noCopy,
		timeout,
		noDoubleBuffering,
	)
	require.Equal(t, expectedOutput, output,
		"quantity: %v, join size: %v, no copy: %v, "+
			"timeout: %v, no double buffering: %v",
		quantity,
		joinSize,
		noCopy,
		timeout,
		noDoubleBuffering,
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
	benchmarkDiscipline(b, false, common.DefaultTestTimeout, 0, false, false)
}

func BenchmarkDisciplineNoDoubleBuffering(b *testing.B) {
	benchmarkDiscipline(b, false, common.DefaultTestTimeout, 0, false, true)
}

func BenchmarkDisciplineNoCopy(b *testing.B) {
	benchmarkDiscipline(b, true, common.DefaultTestTimeout, 0, false, false)
}

func BenchmarkDisciplineNoCopyNoDoubleBuffering(b *testing.B) {
	benchmarkDiscipline(b, true, common.DefaultTestTimeout, 0, false, true)
}

func BenchmarkDisciplineUntimeouted(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0, false, false)
}

func BenchmarkDisciplineUntimeoutedNoDoubleBuffering(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0, false, true)
}

func BenchmarkDisciplineStress(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0, true, false)
}

func BenchmarkDisciplineStressNoDoubleBuffering(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0, true, true)
}

func BenchmarkDisciplineNoStress(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0, false, false)
}

func BenchmarkDisciplineNoStressNoDoubleBuffering(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0, false, true)
}

func BenchmarkDisciplineOutputDelayIsSame(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 1, false, false)
}

func BenchmarkDisciplineOutputDelayIsSameNoDoubleBuffering(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 1, false, true)
}

func BenchmarkDisciplineOutputDelayIsSameStress(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 1, true, false)
}

func BenchmarkDisciplineOutputDelayIsSameNoDoubleBufferingStress(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 1, true, true)
}

func BenchmarkDisciplineOutputDelayIs4TimesLess(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0.25, false, false)
}

func BenchmarkDisciplineOutputDelayIs4TimesLessNoDoubleBuffering(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0.25, false, true)
}

func BenchmarkDisciplineOutputDelayIs2TimesLess(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0.5, false, false)
}

func BenchmarkDisciplineOutputDelayIs2TimesLessNoDoubleBuffering(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0.5, false, true)
}

func BenchmarkDisciplineOutputDelayIs2TimesLonger(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 2, false, false)
}

func BenchmarkDisciplineOutputDelayIs2TimesLongerNoDoubleBuffering(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 2, false, true)
}

func BenchmarkDisciplineOutputDelayIs4TimesLonger(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 4, false, false)
}

func BenchmarkDisciplineOutputDelayIs4TimesLongerNoDoubleBuffering(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 4, false, true)
}

func benchmarkDiscipline(
	b *testing.B,
	noCopy bool,
	timeout time.Duration,
	outputDelayFactor float64,
	stressSystem bool,
	noDoubleBuffering bool,
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
		Input:             input,
		JoinSize:          joinSize,
		NoCopy:            noCopy,
		NoDoubleBuffering: noDoubleBuffering,
		Timeout:           timeout,
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	if stressSystem {
		stress := stressor.New(stressor.Opts{})
		defer stress.Stop()
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
