package unite

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/join/internal/defaults"
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
		Timeout:  consts.ReliablyMeasurableDuration,
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input:    make(chan []int),
		JoinSize: 10,
		Timeout:  defaults.MinTimeout,
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
	for quantity := 100; quantity <= 200; quantity++ {
		for blockSize := 1; blockSize <= 25; blockSize++ {
			for joinSize := uint(1); joinSize <= 20; joinSize++ {
				testDiscipline(
					t,
					quantity,
					blockSize,
					joinSize,
					false,
					defaults.TestTimeout,
				)

				testDiscipline(
					t,
					quantity,
					blockSize,
					joinSize,
					true,
					defaults.TestTimeout,
				)

				testDiscipline(t, quantity, blockSize, joinSize, false, 0)
				testDiscipline(t, quantity, blockSize, joinSize, true, 0)
			}
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
	input := make(chan []int, joinSize)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	discipline, err := New(opts)
	require.NoError(
		t,
		err,
		"quantity: %v, block size: %v, join size: %v, no copy: %v, timeout: %v",
		quantity,
		blockSize,
		joinSize,
		noCopy,
		timeout,
	)

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	expected := inspect.Expected(quantity, blockSize, joinSize)
	output := make([][]int, 0, len(expected))

	go func() {
		defer close(input)

		for _, block := range inspect.Input(quantity, blockSize) {
			inSequence = append(inSequence, block...)

			input <- block
		}
	}()

	for join := range discipline.Output() {
		require.NotEmpty(
			t,
			join,
			"quantity: %v, block size: %v, join size: %v, no copy: %v, timeout: %v",
			quantity,
			blockSize,
			joinSize,
			noCopy,
			timeout,
		)

		output = append(output, append([]int(nil), join...))
		outSequence = append(outSequence, join...)

		if noCopy {
			discipline.Release()
		}
	}

	require.Equal(
		t,
		inSequence,
		outSequence,
		"quantity: %v, block size: %v, join size: %v, no copy: %v, timeout: %v",
		quantity,
		blockSize,
		joinSize,
		noCopy,
		timeout,
	)

	require.Equal(
		t,
		expected,
		output,
		"quantity: %v, block size: %v, join size: %v, no copy: %v, timeout: %v",
		quantity,
		blockSize,
		joinSize,
		noCopy,
		timeout,
	)
}

func TestDisciplineTimeout(t *testing.T) {
	for pauseAt := 50; pauseAt <= 70; pauseAt++ {
		for _, blockSize := range []int{3, 4} {
			t.Run(
				"",
				func(t *testing.T) {
					t.Parallel()
					testDisciplineTimeout(
						t,
						100,
						blockSize,
						10,
						false,
						500*time.Millisecond,
						pauseAt,
					)
				},
			)

			t.Run(
				"",
				func(t *testing.T) {
					t.Parallel()
					testDisciplineTimeout(
						t,
						100,
						blockSize,
						10,
						true,
						500*time.Millisecond,
						pauseAt,
					)
				},
			)
		}
	}
}

func testDisciplineTimeout(
	t *testing.T,
	quantity int,
	blockSize int,
	joinSize uint,
	noCopy bool,
	timeout time.Duration,
	pauseAt int,
) {
	input := make(chan []int, joinSize)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	require.NotZero(t, timeout)

	pauseAt = inspect.PickUpPauseAt(quantity, pauseAt, blockSize, joinSize)
	require.NotZero(
		t,
		pauseAt,
		"quantity: %v, block size: %v, join size: %v, no copy: %v, timeout: %v, "+
			"pause at: %v",
		quantity,
		blockSize,
		joinSize,
		noCopy,
		timeout,
		pauseAt,
	)

	pausetAtDuration := inspect.CalcPauseAtDuration(timeout)

	discipline, err := New(opts)
	require.NoError(
		t,
		err,
		"quantity: %v, block size: %v, join size: %v, no copy: %v, timeout: %v, "+
			"pause at: %v",
		quantity,
		blockSize,
		joinSize,
		noCopy,
		timeout,
		pauseAt,
	)

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	expected := inspect.ExpectedWithTimeout(quantity, pauseAt, blockSize, joinSize)
	output := make([][]int, 0, len(expected))

	go func() {
		defer close(input)

		for _, block := range inspect.Input(quantity, blockSize) {
			for _, item := range block {
				if item == pauseAt {
					time.Sleep(pausetAtDuration)
				}
			}

			inSequence = append(inSequence, block...)

			input <- block
		}
	}()

	for join := range discipline.Output() {
		require.NotEmpty(
			t,
			join,
			"quantity: %v, block size: %v, join size: %v, no copy: %v, timeout: %v, "+
				"pause at: %v",
			quantity,
			blockSize,
			joinSize,
			noCopy,
			timeout,
			pauseAt,
		)

		output = append(output, append([]int(nil), join...))
		outSequence = append(outSequence, join...)

		if noCopy {
			discipline.Release()
		}
	}

	require.Equal(
		t,
		inSequence,
		outSequence,
		"quantity: %v, block size: %v, join size: %v, no copy: %v, timeout: %v, "+
			"pause at: %v",
		quantity,
		blockSize,
		joinSize,
		noCopy,
		timeout,
		pauseAt,
	)

	require.Equal(
		t,
		expected,
		output,
		"quantity: %v, block size: %v, join size: %v  no copy: %v, timeout: %v, "+
			"pause at: %v",
		quantity,
		blockSize,
		joinSize,
		noCopy,
		timeout,
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

	testDisciplineByDataSet(t, nil, [][]int{}, 5, true, defaults.TestTimeout)
	testDisciplineByDataSet(t, dataSet, expected, 5, true, defaults.TestTimeout)
}

func testDisciplineByDataSet(
	t *testing.T,
	dataSet [][]int,
	expected [][]int,
	joinSize uint,
	noCopy bool,
	timeout time.Duration,
) {
	input := make(chan []int, joinSize)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		NoCopy:   noCopy,
		Timeout:  timeout,
	}

	discipline, err := New(opts)
	require.NoError(
		t,
		err,
		"data set: %v, join size: %v, no copy: %v, timeout: %v",
		dataSet,
		joinSize,
		noCopy,
		timeout,
	)

	inSequence := make([]int, 0)
	outSequence := make([]int, 0)

	output := make([][]int, 0, len(expected))

	go func() {
		defer close(input)

		for _, block := range dataSet {
			inSequence = append(inSequence, block...)

			input <- block
		}
	}()

	for join := range discipline.Output() {
		require.NotEmpty(
			t,
			join,
			"data set: %v, join size: %v, no copy: %v, timeout: %v",
			dataSet,
			joinSize,
			noCopy,
			timeout,
		)

		output = append(output, append([]int(nil), join...))
		outSequence = append(outSequence, join...)

		if noCopy {
			discipline.Release()
		}
	}

	require.Equal(
		t,
		inSequence,
		outSequence,
		"data set: %v, join size: %v, no copy: %v, timeout: %v",
		dataSet,
		joinSize,
		noCopy,
		timeout,
	)

	require.Equal(
		t,
		expected,
		output,
		"data set: %v, join size: %v, no copy: %v, timeout: %v",
		dataSet,
		joinSize,
		noCopy,
		timeout,
	)
}

func BenchmarkDiscipline(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, false, defaults.TestTimeout, 0, false)
}

func BenchmarkDisciplineInputCapIsFullJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, false, defaults.TestTimeout, 1, false)
}

func BenchmarkDisciplineNoCopy(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, true, defaults.TestTimeout, 0, false)
}

func BenchmarkDisciplineNoCopyInputCapIsFullJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, 4, true, defaults.TestTimeout, 1, false)
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
