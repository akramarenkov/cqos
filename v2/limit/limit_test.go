package limit

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/general"

	"github.com/stretchr/testify/require"
)

func TestOptsValidation(t *testing.T) {
	opts := Opts[uint]{}

	_, err := New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Input: make(chan uint),
	}

	_, err = New(opts)
	require.Error(t, err)

	opts = Opts[uint]{
		Input: make(chan uint),
		Limit: Rate{Interval: time.Second, Quantity: 1},
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func TestDiscipline(t *testing.T) {
	quantity := uint(1000)

	limit := Rate{
		Interval: time.Second,
		Quantity: 1000,
	}

	disciplined := testDiscipline(t, quantity, limit, true, 0.1)
	undisciplined := testUndisciplined(t, quantity)

	require.Less(t, undisciplined, disciplined)
}

func TestDisciplineAlignedQuantity(t *testing.T) {
	quantity := uint(1000)

	limit := Rate{
		Interval: time.Second,
		Quantity: uint64(quantity),
	}

	disciplined := testDiscipline(t, quantity, limit, false, 0.1)
	undisciplined := testUndisciplined(t, quantity)

	require.Less(t, undisciplined, disciplined)
}

func TestDisciplineUnalignedQuantity(t *testing.T) {
	quantity := uint(1500)

	limit := Rate{
		Interval: time.Second,
		Quantity: 1000,
	}

	disciplined := testDiscipline(t, quantity, limit, false, 0.1)
	undisciplined := testUndisciplined(t, quantity)

	require.Less(t, undisciplined, disciplined)
}

func testDiscipline(
	t *testing.T,
	quantity uint,
	limit Rate,
	optimize bool,
	maxRelativeDurationDeviation float64,
) time.Duration {
	if optimize {
		optimized, err := limit.Optimize()
		require.NoError(t, err)

		limit = optimized
	}

	capacity := general.CalcByFactor(
		int(quantity),
		defaultCapacityFactor,
		1,
	)

	input := make(chan uint, capacity)

	opts := Opts[uint]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	inSequence := make([]uint, 0, quantity)
	outSequence := make([]uint, 0, quantity)

	startedAt := time.Now()

	go func() {
		defer close(input)

		for stage := uint(0); stage < quantity; stage++ {
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
		maxRelativeDurationDeviation,
	)

	require.Equal(t, inSequence, outSequence)
	require.InDelta(t, expectedDuration, duration, acceptableDeviation)

	return duration
}

func calcExpectedDuration(
	quantity uint,
	limit Rate,
	relativeDeviation float64,
) (time.Duration, float64) {
	duration := (time.Duration(quantity) * limit.Interval) / time.Duration(limit.Quantity)
	deviation := relativeDeviation * float64(duration)

	return duration, deviation
}

func testUndisciplined(t *testing.T, quantity uint) time.Duration {
	capacity := general.CalcByFactor(
		int(quantity),
		defaultCapacityFactor,
		1,
	)

	input := make(chan uint, capacity)

	inSequence := make([]uint, 0, quantity)
	outSequence := make([]uint, 0, quantity)

	startedAt := time.Now()

	go func() {
		defer close(input)

		for stage := uint(0); stage < quantity; stage++ {
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

func benchmarkDiscipline(b *testing.B, quantity uint, limit Rate, optimize bool) {
	if optimize {
		optimized, err := limit.Optimize()
		require.NoError(b, err)

		limit = optimized
	}

	capacity := general.CalcByFactor(
		int(quantity),
		defaultCapacityFactor,
		1,
	)

	input := make(chan uint, capacity)

	opts := Opts[uint]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	go func() {
		defer close(input)

		for stage := uint(0); stage < quantity; stage++ {
			input <- stage
		}
	}()

	for range discipline.Output() { // nolint:revive
	}
}

func BenchmarkDiscipline(b *testing.B) {
	quantity := uint(1e5)

	limit := Rate{
		Interval: time.Second,
		Quantity: uint64(quantity),
	}

	benchmarkDiscipline(b, quantity, limit, false)
}

func BenchmarkDisciplineOptimize(b *testing.B) {
	quantity := uint(1e6)

	limit := Rate{
		Interval: time.Second,
		Quantity: uint64(quantity),
	}

	benchmarkDiscipline(b, quantity, limit, true)
}
