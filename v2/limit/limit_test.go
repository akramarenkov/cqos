package limit

import (
	"math"
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/general"

	"github.com/stretchr/testify/require"
)

func TestOptsValidation(t *testing.T) {
	opts := Opts[int]{}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input: make(chan int),
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[int]{
		Input: make(chan int),
		Limit: Rate{Interval: time.Second, Quantity: 1},
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func TestIncreaseDelay(t *testing.T) {
	require.Equal(
		t,
		time.Duration(math.MinInt64+1),
		increaseDelay(time.Duration(math.MinInt64), 1),
	)

	require.Equal(
		t,
		time.Duration(math.MinInt64),
		increaseDelay(time.Duration(math.MinInt64+1), -1),
	)

	require.Equal(
		t,
		time.Duration(0),
		increaseDelay(time.Duration(math.MinInt64), -1),
	)

	require.Equal(
		t,
		time.Duration(math.MaxInt64-1),
		increaseDelay(time.Duration(math.MaxInt64), -1),
	)

	require.Equal(
		t,
		time.Duration(math.MaxInt64),
		increaseDelay(time.Duration(math.MaxInt64-1), 1),
	)

	require.Equal(
		t,
		time.Duration(0),
		increaseDelay(time.Duration(math.MaxInt64), 1),
	)
}

func TestDiscipline(t *testing.T) {
	quantity := 10000

	limit := Rate{
		Interval: time.Second,
		Quantity: 1000,
	}

	disciplined := testDiscipline(t, quantity, limit, false, 0.1)
	undisciplined := testUndisciplined(t, quantity)

	require.Less(t, undisciplined, disciplined)
}

func TestDisciplineOptimize(t *testing.T) {
	quantity := 10000

	limit := Rate{
		Interval: time.Second,
		Quantity: 1000,
	}

	disciplined := testDiscipline(t, quantity, limit, true, 0.1)
	undisciplined := testUndisciplined(t, quantity)

	require.Less(t, undisciplined, disciplined)
}

func BenchmarkDiscipline(b *testing.B) {
	quantity := int(11e6)

	limit := Rate{
		Interval: time.Second,
		Quantity: uint64(quantity),
	}

	benchmarkDiscipline(b, quantity, limit, false, 0.1)
}

func BenchmarkDisciplineOptimize(b *testing.B) {
	quantity := int(11e6)

	limit := Rate{
		Interval: time.Second,
		Quantity: uint64(quantity),
	}

	benchmarkDiscipline(b, quantity, limit, true, 0.1)
}

func testDiscipline(
	t *testing.T,
	quantity int,
	limit Rate,
	optimize bool,
	maxRelativeDeviation float64,
) time.Duration {
	if optimize {
		optimized, err := limit.Optimize()
		require.NoError(t, err)

		limit = optimized
	}

	capacity := general.CalcByFactor(
		quantity,
		defaultCapacityFactor,
		1,
	)

	input := make(chan int, capacity)

	opts := Opts[int]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	startedAt := time.Now()

	go func() {
		defer close(input)

		for stage := 0; stage < quantity; stage++ {
			inSequence = append(inSequence, stage)

			input <- stage
		}
	}()

	for item := range discipline.Output() {
		outSequence = append(outSequence, item)
	}

	duration := time.Since(startedAt)

	expectedDuration, acceptableDeviation := calcExpectedDuration(
		quantity,
		limit,
		maxRelativeDeviation,
	)

	require.Equal(t, inSequence, outSequence)
	require.InDelta(t, expectedDuration, duration, acceptableDeviation)

	return duration
}

func calcExpectedDuration(
	quantity int,
	limit Rate,
	relativeDeviation float64,
) (time.Duration, float64) {
	duration := (time.Duration(quantity) * limit.Interval) / time.Duration(limit.Quantity)
	deviation := relativeDeviation * float64(duration)

	return duration, deviation
}

func testUndisciplined(t *testing.T, quantity int) time.Duration {
	capacity := general.CalcByFactor(
		quantity,
		defaultCapacityFactor,
		1,
	)

	input := make(chan int, capacity)

	inSequence := make([]int, 0, quantity)
	outSequence := make([]int, 0, quantity)

	startedAt := time.Now()

	go func() {
		defer close(input)

		for stage := 0; stage < quantity; stage++ {
			inSequence = append(inSequence, stage)

			input <- stage
		}
	}()

	for item := range input {
		outSequence = append(outSequence, item)
	}

	duration := time.Since(startedAt)

	require.Equal(t, inSequence, outSequence)

	return duration
}

func benchmarkDiscipline(
	b *testing.B,
	quantity int,
	limit Rate,
	optimize bool,
	maxRelativeDurationDeviation float64,
) {
	if optimize {
		optimized, err := limit.Optimize()
		require.NoError(b, err)

		limit = optimized
	}

	capacity := general.CalcByFactor(
		quantity,
		defaultCapacityFactor,
		1,
	)

	input := make(chan int, capacity)

	opts := Opts[int]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	startedAt := time.Now()

	go func() {
		defer close(input)

		for stage := 0; stage < quantity; stage++ {
			input <- stage
		}
	}()

	for range discipline.Output() { // nolint:revive
	}

	duration := time.Since(startedAt)

	expectedDuration, acceptableDeviation := calcExpectedDuration(
		quantity,
		limit,
		maxRelativeDurationDeviation,
	)

	require.InDelta(b, expectedDuration, duration, acceptableDeviation)
}
