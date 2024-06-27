package unite

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
		Input: make(chan []int),
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input:    make(chan []int),
		JoinSize: 10,
		Timeout:  general.ReliablyMeasurableDuration,
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
		for blockSize := 1; blockSize <= 21; blockSize++ {
			testDiscipline(
				t,
				quantity,
				blockSize,
				10,
				false,
				common.DefaultTestTimeout,
			)
		}
	}
}

func TestDisciplineNoCopy(t *testing.T) {
	for quantity := 100; quantity <= 125; quantity++ {
		for blockSize := 1; blockSize <= 21; blockSize++ {
			testDiscipline(
				t,
				quantity,
				blockSize,
				10,
				true,
				common.DefaultTestTimeout,
			)
		}
	}
}

func TestDisciplineUntimeouted(t *testing.T) {
	for quantity := 100; quantity <= 125; quantity++ {
		for blockSize := 1; blockSize <= 21; blockSize++ {
			testDiscipline(t, quantity, blockSize, 10, false, 0)
		}
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

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	expectedOutput := inspect.Expected(quantity, blockSize, opts.JoinSize)
	output := make([][]int, 0, len(expectedOutput))

	go func() {
		defer close(input)

		for _, slice := range inspect.Input(quantity, blockSize) {
			inSequence = append(inSequence, slice...)

			input <- slice
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
		"quantity: %v, block size: %v, join size: %v, "+
			"no copy: %v, timeout: %v",
		quantity,
		blockSize,
		joinSize,
		noCopy,
		timeout,
	)
	require.Equal(t, expectedOutput, output,
		"quantity: %v, block size: %v, join size: %v, "+
			"no copy: %v, timeout: %v",
		quantity,
		blockSize,
		joinSize,
		noCopy,
		timeout,
	)
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

	pauseAt = inspect.PickUpPauseAt(quantity, pauseAt, blockSize, opts.JoinSize)
	require.NotEqual(t, 0, pauseAt)

	pausetAtDuration := inspect.CalcPauseAtDuration(opts.Timeout)

	discipline, err := New(opts)
	require.NoError(t, err)

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	expectedOutput := inspect.ExpectedWithTimeout(
		quantity,
		pauseAt,
		blockSize,
		opts.JoinSize,
	)

	output := make([][]int, 0, len(expectedOutput))

	go func() {
		defer close(input)

		for _, slice := range inspect.Input(quantity, blockSize) {
			for _, item := range slice {
				if item == pauseAt {
					time.Sleep(pausetAtDuration)
				}
			}

			inSequence = append(inSequence, slice...)

			input <- slice
		}
	}()

	for slice := range discipline.Output() {
		require.NotEmpty(t, slice)

		outSequence = append(outSequence, slice...)

		output = append(output, append([]int(nil), slice...))
	}

	require.Equal(t, inSequence, outSequence,
		"quantity: %v, block size: %v, join size: %v, pause at: %v",
		quantity,
		blockSize,
		joinSize,
		pauseAt,
	)
	require.Equal(t, expectedOutput, output,
		"quantity: %v, block size: %v, join size: %v, pause at: %v",
		quantity,
		blockSize,
		joinSize,
		pauseAt,
	)
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

	expected := [][]int{
		{1, 2},
		{3, 4, 5, 6, 7},
		{8, 9, 10},
		{11, 12, 13, 14, 15, 16},
		{17, 18, 19},
		{20, 21, 22},
		{23, 24, 25, 26, 27},
		{28, 29, 30},
	}

	testDisciplineByDataSet(t, nil, 5, [][]int{}, true, common.DefaultTestTimeout)
	testDisciplineByDataSet(t, dataSet, 5, expected, true, common.DefaultTestTimeout)
}

func testDisciplineByDataSet(
	t *testing.T,
	dataSet [][]int,
	joinSize uint,
	expectedOutput [][]int,
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

	inSequence := make([]int, 0)
	outSequence := make([]int, 0)

	output := make([][]int, 0)

	go func() {
		defer close(input)

		for _, slice := range dataSet {
			inSequence = append(inSequence, slice...)

			input <- slice
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
		"data set: %v, join size: %v, no copy: %v, timeout: %v",
		dataSet,
		joinSize,
		noCopy,
		timeout,
	)
	require.Equal(t, expectedOutput, output,
		"data set: %v, join size: %v, no copy: %v, timeout: %v",
		dataSet,
		joinSize,
		noCopy,
		timeout,
	)
}

func BenchmarkDiscipline(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, false, common.DefaultTestTimeout, 0, false)
}

func BenchmarkDisciplineInputCapIsFullJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, false, common.DefaultTestTimeout, 1, false)
}

func BenchmarkDisciplineNoCopy(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, true, common.DefaultTestTimeout, 0, false)
}

func BenchmarkDisciplineNoCopyInputCapIsFullJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, true, common.DefaultTestTimeout, 1, false)
}

func BenchmarkDisciplineNoCopyUntimeouted(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, true, 0, 0, false)
}

func BenchmarkDisciplineNoCopyUntimeoutedInputCapIsHalfJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, true, 0, 0.5, false)
}

func BenchmarkDisciplineNoCopyUntimeoutedInputCapIsFullJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, true, 0, 1, false)
}

func BenchmarkDisciplineNoCopyUntimeoutedInputCapIsTwiceJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, true, 0, 2, false)
}

func BenchmarkDisciplineNoCopyUntimeoutedInputCapIsTripleJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, true, 0, 3, false)
}

func BenchmarkDisciplineNoCopyUntimeoutedStress(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, true, 0, 0, true)
}

func BenchmarkDisciplineNoCopyUntimeoutedInputCapIsFullJoinSizeStress(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, true, 0, 1, true)
}

func BenchmarkDisciplineUntimeoutedStress(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, false, 0, 0, true)
}

func BenchmarkDisciplineUntimeoutedInputCapIsFullJoinSizeStress(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, false, 0, 1, true)
}

func benchmarkDiscipline(
	b *testing.B,
	joinSize uint,
	blockSize int,
	noCopy bool,
	timeout time.Duration,
	inputCapFactor float64,
	stressSystem bool,
) {
	joinsQuantity := b.N
	effectiveJoinSize := blockSize * (int(joinSize) / blockSize)
	quantity := joinsQuantity * effectiveJoinSize
	inputCap := int(inputCapFactor * float64(joinSize))

	input := make(chan []int, inputCap)

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

	blocks := inspect.Input(quantity, blockSize)

	b.ResetTimer()

	go func() {
		defer close(input)

		for _, slice := range blocks {
			input <- slice
		}
	}()

	for range discipline.Output() {
		if noCopy {
			discipline.Release()
		}
	}
}
