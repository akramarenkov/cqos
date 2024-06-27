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
	for quantity := 100; quantity <= 125; quantity++ {
		testDiscipline(t, quantity, 10, false, common.DefaultTestTimeout)
	}
}

func TestDisciplineRelease(t *testing.T) {
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
	useReleased bool,
	timeout time.Duration,
) {
	input := make(chan int)

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

		if useReleased {
			released <- struct{}{}
		}
	}

	require.Equal(t, inSequence, outSequence,
		"quantity: %v, join size: %v, use released: %v, timeout: %v",
		quantity,
		joinSize,
		useReleased,
		timeout,
	)
	require.Equal(t, expectedOutput, output,
		"quantity: %v, join size: %v, use released: %v, timeout: %v",
		quantity,
		joinSize,
		useReleased,
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

func TestDisciplineStop(t *testing.T) {
	testDisciplineStop(t, 105, 10, 52, false, false, common.DefaultTestTimeout)
	testDisciplineStop(t, 105, 10, 52, false, false, 0)
	testDisciplineStop(t, 105, 10, 52, false, true, 0)
}

func TestDisciplineStopByCtx(t *testing.T) {
	testDisciplineStop(t, 105, 10, 52, true, false, common.DefaultTestTimeout)
	testDisciplineStop(t, 105, 10, 52, true, false, 0)
	testDisciplineStop(t, 105, 10, 52, true, true, 0)
}

func testDisciplineStop(
	t *testing.T,
	quantity int,
	joinSize uint,
	stopAt int,
	byCtx bool,
	useReleased bool,
	timeout time.Duration,
) {
	input := make(chan int)

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
	require.NoError(t, err)

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

		for _, slice := range inspect.Input(quantity, 1) {
			for _, item := range slice {
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

	for slice := range discipline.Output() {
		require.NotEmpty(t, slice)

		outSequence = append(outSequence, slice...)

		if useReleased {
			stop()
			interrupter.Close()
		}
	}

	wg.Wait()

	if useReleased {
		require.Less(t, len(outSequence), quantity)
		require.GreaterOrEqual(t, len(outSequence), int(opts.JoinSize))
	} else {
		require.Less(
			t,
			len(inSequence)-len(outSequence),
			2*int(opts.JoinSize),
			"in: %v, out: %v",
			inSequence, outSequence,
		)

		require.GreaterOrEqual(t, len(inSequence)-len(outSequence), 0)
	}

	require.Subset(t, inSequence, outSequence)
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
