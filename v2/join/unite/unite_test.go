package unite

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/internal/stressor"
	"github.com/akramarenkov/cqos/v2/join/internal/blocks"
	"github.com/akramarenkov/cqos/v2/join/internal/common"

	"github.com/stretchr/testify/require"
)

func TestOptsValidation(t *testing.T) {
	opts := Opts[int]{
		JoinSize: 10,
	}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input: make(chan []int),
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input:    make(chan []int),
		JoinSize: 10,
		Timeout:  consts.ReliablyMeasurableDuration,
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input:    make(chan []int),
		JoinSize: 10,
		Timeout:  common.DefaultMinTimeout,
	}

	_, err = New(opts)
	require.NoError(t, err)

	opts = Opts[int]{
		Input:    make(chan []int),
		JoinSize: 10,
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func TestDiscipline(t *testing.T) {
	for quantity := 100; quantity <= 125; quantity++ {
		testDiscipline(t, quantity, 1, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 2, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 3, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 4, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 5, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 6, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 7, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 8, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 9, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 10, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 11, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 12, 10, false, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 13, 10, false, common.DefaultTestTimeout)
	}
}

func TestDisciplineNoCopy(t *testing.T) {
	for quantity := 100; quantity <= 125; quantity++ {
		testDiscipline(t, quantity, 1, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 2, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 3, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 4, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 5, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 6, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 7, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 8, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 9, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 10, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 11, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 12, 10, true, common.DefaultTestTimeout)
		testDiscipline(t, quantity, 13, 10, true, common.DefaultTestTimeout)
	}
}

func TestDisciplineUntimeouted(t *testing.T) {
	for quantity := 100; quantity <= 125; quantity++ {
		testDiscipline(t, quantity, 1, 10, false, 0)
		testDiscipline(t, quantity, 2, 10, false, 0)
		testDiscipline(t, quantity, 3, 10, false, 0)
		testDiscipline(t, quantity, 4, 10, false, 0)
		testDiscipline(t, quantity, 5, 10, false, 0)
		testDiscipline(t, quantity, 6, 10, false, 0)
		testDiscipline(t, quantity, 7, 10, false, 0)
		testDiscipline(t, quantity, 8, 10, false, 0)
		testDiscipline(t, quantity, 9, 10, false, 0)
		testDiscipline(t, quantity, 10, 10, false, 0)
		testDiscipline(t, quantity, 11, 10, false, 0)
		testDiscipline(t, quantity, 12, 10, false, 0)
		testDiscipline(t, quantity, 13, 10, false, 0)
	}
}

func testDiscipline(
	t *testing.T,
	quantity int,
	blockSize int,
	joinSize uint,
	noCopy bool,
	timeout time.Duration,
) {
	input := make(chan []int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
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

		for _, slice := range blocks.DivideSequence(quantity, blockSize) {
			inSequence = append(inSequence, slice...)

			input <- slice
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

	expectedJoins := blocks.CalcExpectedJoins(quantity, blockSize, opts.JoinSize)

	require.Equal(t, inSequence, outSequence)
	require.Equal(t, expectedJoins, joins)
}

func TestDisciplineTimeout(t *testing.T) {
	for quantity := 100; quantity <= 125; quantity++ {
		testDisciplineTimeout(t, quantity, 4, 10, 53)
	}
}

func testDisciplineTimeout(
	t *testing.T,
	quantity int,
	blockSize int,
	joinSize uint,
	pauseAt int,
) {
	input := make(chan []int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		Timeout:  100 * time.Millisecond,
	}

	pauseAt = blocks.PickUpPauseAt(quantity, pauseAt, blockSize, opts.JoinSize)
	require.NotEqual(t, 0, pauseAt)

	discipline, err := New(opts)
	require.NoError(t, err)

	joins := 0

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	go func() {
		defer close(input)

		for _, slice := range blocks.DivideSequence(quantity, blockSize) {
			for _, item := range slice {
				if item == pauseAt {
					time.Sleep(5 * opts.Timeout)
				}
			}

			inSequence = append(inSequence, slice...)

			input <- slice
		}
	}()

	for slice := range discipline.Output() {
		require.NotEqual(t, 0, slice)

		joins++

		outSequence = append(outSequence, slice...)
	}

	expectedJoins := blocks.CalcExpectedJoinsWithTimeout(
		quantity,
		pauseAt,
		blockSize,
		opts.JoinSize,
	)

	require.Equal(t, inSequence, outSequence)
	require.Equal(t, expectedJoins, joins)
}

func TestDisciplineManually(t *testing.T) {
	dataSet := [][]int{
		{},                       // nothing has been done
		{1, 2},                   // add this slice into join
		{3, 4, 5, 6, 7},          // pass join and then this slice (2)
		{8, 9, 10},               // add this slice into join
		{11, 12, 13, 14, 15, 16}, // pass join and then this slice (2)
		{17, 18, 19},             // add this slice into join
		{20, 21, 22},             // pass join and add this slice into join (1)
		{},                       // nothing has been done
		{23, 24, 25},             // pass join and add this slice into join (1)
		{26, 27},                 // add this slice into join and pass join (1)
		{28, 29, 30},             // add this slice into join and pass join at close input (1)
	}

	testDisciplineByDataSet(t, nil, 5, 0, true, common.DefaultTestTimeout)
	testDisciplineByDataSet(t, dataSet, 5, 2+2+1+1+1+1, true, common.DefaultTestTimeout)
}

func testDisciplineByDataSet(
	t *testing.T,
	dataSet [][]int,
	joinSize uint,
	expectedJoins int,
	noCopy bool,
	timeout time.Duration,
) {
	input := make(chan []int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	joins := 0

	inSequence := make([]int, 0)
	outSequence := make([]int, 0)

	go func() {
		defer close(input)

		for _, slice := range dataSet {
			inSequence = append(inSequence, slice...)

			input <- slice
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

	require.Equal(t, inSequence, outSequence)
	require.Equal(t, expectedJoins, joins)
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

func BenchmarkDisciplineStress(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0, true)
}

func BenchmarkDisciplineNoStress(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0, false)
}

func BenchmarkDisciplineOutputDelayIsSame(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 1, false)
}

func BenchmarkDisciplineOutputDelayIs4TimesLess(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0.25, false)
}

func BenchmarkDisciplineOutputDelayIs2TimesLess(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 0.5, false)
}

func BenchmarkDisciplineOutputDelayIs2TimesLonger(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 2, false)
}

func BenchmarkDisciplineOutputDelayIs4TimesLonger(b *testing.B) {
	benchmarkDiscipline(b, true, 0, 4, false)
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
	blockSize := 4
	// Accuracy of this delay is sufficient for the benchmark and
	// consts.ReliablyMeasurableDuration is too large to perform a representative
	// number of iterations
	inputDelayBase := 1 * time.Millisecond

	effectiveJoinSize := blockSize * (int(joinSize) / blockSize)
	quantity := joinsQuantity * effectiveJoinSize

	inputDelay, outputDelay := calcBenchmarkDelays(
		inputDelayBase,
		joinSize,
		blockSize,
		outputDelayFactor,
	)

	blocks := blocks.DivideSequence(quantity, blockSize)

	input := make(chan []int)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	if stressSystem {
		stress, err := stressor.New(0, 0)
		require.NoError(b, err)

		defer stress.Stop()
	}

	b.ResetTimer()

	go func() {
		defer close(input)

		for _, slice := range blocks {
			time.Sleep(inputDelay)

			input <- slice
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
	blockSize int,
	outputDelayFactor float64,
) (time.Duration, time.Duration) {
	blocksInJoin := int(joinSize) / blockSize

	outputDelay := time.Duration(
		outputDelayFactor *
			float64(inputDelay) *
			float64(blocksInJoin),
	)

	if outputDelay == 0 {
		inputDelay = 0
	}

	return inputDelay, outputDelay
}