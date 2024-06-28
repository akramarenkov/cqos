package join

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/internal/general"
	"github.com/akramarenkov/cqos/join/internal/common"
	"github.com/akramarenkov/cqos/join/internal/inspect"

	"github.com/akramarenkov/breaker/closing"
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
	for quantity := 100; quantity <= 200; quantity++ {
		for joinSize := uint(1); joinSize <= 20; joinSize++ {
			testDiscipline(t, quantity, joinSize, false, common.DefaultTestTimeout)
			testDiscipline(t, quantity, joinSize, true, common.DefaultTestTimeout)
			testDiscipline(t, quantity, joinSize, false, 0)
			testDiscipline(t, quantity, joinSize, true, 0)
		}
	}
}

func testDiscipline(
	t *testing.T,
	quantity int,
	joinSize uint,
	useReleased bool,
	timeout time.Duration,
) {
	input := make(chan int, joinSize)

	released := make(chan struct{})
	defer close(released)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		Timeout:  timeout,
	}

	if useReleased {
		opts.Released = released
	}

	discipline, err := New(opts)
	require.NoError(
		t,
		err,
		"quantity: %v, join size: %v, use released: %v, timeout: %v",
		quantity,
		joinSize,
		useReleased,
		timeout,
	)

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	expected := inspect.Expected(quantity, 1, joinSize)
	output := make([][]int, 0, len(expected))

	go func() {
		defer close(input)

		for _, block := range inspect.Input(quantity, 1) {
			for _, item := range block {
				inSequence = append(inSequence, item)

				input <- item
			}
		}
	}()

	for join := range discipline.Output() {
		require.NotEmpty(
			t,
			join,
			"quantity: %v, join size: %v, use released: %v, timeout: %v",
			quantity,
			joinSize,
			useReleased,
			timeout,
		)

		output = append(output, append([]int(nil), join...))
		outSequence = append(outSequence, join...)

		if useReleased {
			released <- struct{}{}
		}
	}

	require.Equal(
		t,
		inSequence,
		outSequence,
		"quantity: %v, join size: %v, use released: %v, timeout: %v",
		quantity,
		joinSize,
		useReleased,
		timeout,
	)

	require.Equal(
		t,
		expected,
		output,
		"quantity: %v, join size: %v, use released: %v, timeout: %v",
		quantity,
		joinSize,
		useReleased,
		timeout,
	)
}

func TestDisciplineTimeout(t *testing.T) {
	for pauseAt := 50; pauseAt <= 70; pauseAt++ {
		t.Run(
			"",
			func(t *testing.T) {
				t.Parallel()
				testDisciplineTimeout(t, 100, 10, false, 500*time.Millisecond, pauseAt)
			},
		)

		t.Run(
			"",
			func(t *testing.T) {
				t.Parallel()
				testDisciplineTimeout(t, 100, 10, true, 500*time.Millisecond, pauseAt)
			},
		)
	}
}

func testDisciplineTimeout(
	t *testing.T,
	quantity int,
	joinSize uint,
	useReleased bool,
	timeout time.Duration,
	pauseAt int,
) {
	input := make(chan int, joinSize)

	released := make(chan struct{})
	defer close(released)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		Timeout:  timeout,
	}

	if useReleased {
		opts.Released = released
	}

	require.NotZero(t, timeout)

	pauseAt = inspect.PickUpPauseAt(quantity, pauseAt, 1, joinSize)
	require.NotZero(
		t,
		pauseAt,
		"quantity: %v, join size: %v, use released: %v, timeout: %v, pause at: %v",
		quantity,
		joinSize,
		useReleased,
		timeout,
		pauseAt,
	)

	pausetAtDuration := inspect.CalcPauseAtDuration(timeout)

	discipline, err := New(opts)
	require.NoError(
		t,
		err,
		"quantity: %v, join size: %v, use released: %v, timeout: %v, pause at: %v",
		quantity,
		joinSize,
		useReleased,
		timeout,
		pauseAt,
	)

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	expected := inspect.ExpectedWithTimeout(quantity, pauseAt, 1, joinSize)
	output := make([][]int, 0, len(expected))

	go func() {
		defer close(input)

		for _, block := range inspect.Input(quantity, 1) {
			for _, item := range block {
				if item == pauseAt {
					time.Sleep(pausetAtDuration)
				}

				inSequence = append(inSequence, item)

				input <- item
			}
		}
	}()

	for join := range discipline.Output() {
		require.NotEmpty(
			t,
			join,
			"quantity: %v, join size: %v, use released: %v, timeout: %v, pause at: %v",
			quantity,
			joinSize,
			useReleased,
			timeout,
			pauseAt,
		)

		output = append(output, append([]int(nil), join...))
		outSequence = append(outSequence, join...)

		if useReleased {
			released <- struct{}{}
		}
	}

	require.Equal(
		t,
		inSequence,
		outSequence,
		"quantity: %v, join size: %v, use released: %v, timeout: %v, pause at: %v",
		quantity,
		joinSize,
		useReleased,
		timeout,
		pauseAt,
	)

	require.Equal(
		t,
		expected,
		output,
		"quantity: %v, join size: %v, use released: %v, timeout: %v, pause at: %v",
		quantity,
		joinSize,
		useReleased,
		timeout,
		pauseAt,
	)
}

func TestDisciplineStop(t *testing.T) {
	quantity := 200

	for stopAt := 1; stopAt <= quantity; stopAt++ {
		testDisciplineStop(t, quantity, 10, false, common.DefaultTestTimeout, stopAt, false)
		testDisciplineStop(t, quantity, 10, true, common.DefaultTestTimeout, stopAt, false)
		testDisciplineStop(t, quantity, 10, false, 0, stopAt, false)
		testDisciplineStop(t, quantity, 10, true, 0, stopAt, false)

		testDisciplineStop(t, quantity, 10, false, common.DefaultTestTimeout, stopAt, true)
		testDisciplineStop(t, quantity, 10, true, common.DefaultTestTimeout, stopAt, true)
		testDisciplineStop(t, quantity, 10, false, 0, stopAt, true)
		testDisciplineStop(t, quantity, 10, true, 0, stopAt, true)
	}
}

func testDisciplineStop(
	t *testing.T,
	quantity int,
	joinSize uint,
	useReleased bool,
	timeout time.Duration,
	stopAt int,
	byCtx bool,
) {
	input := make(chan int, joinSize)

	released := make(chan struct{})
	defer close(released)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := Opts[int]{
		Ctx:      ctx,
		Input:    input,
		JoinSize: joinSize,
		Timeout:  timeout,
	}

	if useReleased {
		opts.Released = released
	}

	discipline, err := New(opts)
	require.NoError(
		t,
		err,
		"quantity: %v, join size: %v, use released: %v, timeout: %v, stopAt: %v"+
			"byCtx: %v",
		quantity,
		joinSize,
		useReleased,
		timeout,
		stopAt,
		byCtx,
	)

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	stop := func() {
		if byCtx {
			cancel()
			return
		}

		discipline.Stop()
	}

	interrupter := closing.New()

	wg := &sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(input)

		for _, block := range inspect.Input(quantity, 1) {
			for _, item := range block {
				if item == stopAt {
					stop()
					return
				}

				inSequence = append(inSequence, item)

				select {
				case <-interrupter.IsClosed():
					return
				case input <- item:
				}
			}
		}
	}()

	for join := range discipline.Output() {
		require.NotEmpty(
			t,
			join,
			"quantity: %v, join size: %v, use released: %v, timeout: %v, stopAt: %v"+
				"byCtx: %v",
			quantity,
			joinSize,
			useReleased,
			timeout,
			stopAt,
			byCtx,
		)

		outSequence = append(outSequence, join...)

		if useReleased {
			stop()
			interrupter.Close()
		}
	}

	wg.Wait()

	require.Subset(
		t,
		inSequence,
		outSequence,
		"quantity: %v, join size: %v, use released: %v, timeout: %v, stopAt: %v"+
			"byCtx: %v",
		quantity,
		joinSize,
		useReleased,
		timeout,
		stopAt,
		byCtx,
	)
}

func BenchmarkDiscipline(b *testing.B) {
	benchmarkDiscipline(b, 10, false, common.DefaultTestTimeout, 0, false)
}

func BenchmarkDisciplineInputCapIsFullJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, false, common.DefaultTestTimeout, 1, false)
}

func BenchmarkDisciplineNoCopy(b *testing.B) {
	benchmarkDiscipline(b, 10, true, common.DefaultTestTimeout, 0, false)
}

func BenchmarkDisciplineNoCopyInputCapIsFullJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, true, common.DefaultTestTimeout, 1, false)
}

func BenchmarkDisciplineNoCopyUntimeouted(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 0, false)
}

func BenchmarkDisciplineNoCopyUntimeoutedInputCapIsHalfJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 0.5, false)
}

func BenchmarkDisciplineNoCopyUntimeoutedInputCapIsFullJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 1, false)
}

func BenchmarkDisciplineNoCopyUntimeoutedInputCapIsTwiceJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 2, false)
}

func BenchmarkDisciplineNoCopyUntimeoutedInputCapIsTripleJoinSize(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 3, false)
}

func BenchmarkDisciplineNoCopyUntimeoutedStress(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 0, true)
}

func BenchmarkDisciplineNoCopyUntimeoutedInputCapIsFullJoinSizeStress(b *testing.B) {
	benchmarkDiscipline(b, 10, true, 0, 1, true)
}

func BenchmarkDisciplineUntimeoutedStress(b *testing.B) {
	benchmarkDiscipline(b, 10, false, 0, 0, true)
}

func BenchmarkDisciplineUntimeoutedInputCapIsFullJoinSizeStress(b *testing.B) {
	benchmarkDiscipline(b, 10, false, 0, 1, true)
}

func benchmarkDiscipline(
	b *testing.B,
	joinSize uint,
	useReleased bool,
	timeout time.Duration,
	inputCapFactor float64,
	stressSystem bool,
) {
	joinsQuantity := b.N
	quantity := joinsQuantity * int(joinSize)
	inputCap := int(inputCapFactor * float64(joinSize))

	input := make(chan int, inputCap)

	released := make(chan struct{})
	defer close(released)

	opts := Opts[int]{
		Input:    input,
		JoinSize: joinSize,
		Timeout:  timeout,
	}

	if useReleased {
		opts.Released = released
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
			input <- item
		}
	}()

	for range discipline.Output() {
		if useReleased {
			released <- struct{}{}
		}
	}
}
