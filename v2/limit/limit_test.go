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
}

func TestDiscipline(t *testing.T) {
	quantity := uint(5003)

	limit := Rate{
		Interval: time.Second,
		Quantity: 1000,
	}

	disciplined := testDiscipline(t, quantity, limit, true)
	undisciplined := testUndisciplined(t, quantity)

	require.Less(t, undisciplined, disciplined)
}

func TestDisciplineAlignedQuantity(t *testing.T) {
	quantity := uint(1000)

	limit := Rate{
		Interval: time.Second,
		Quantity: uint64(quantity),
	}

	disciplined := testDiscipline(t, quantity, limit, false)
	undisciplined := testUndisciplined(t, quantity)

	require.Less(t, undisciplined, disciplined)
}

func TestDisciplineUnalignedQuantity(t *testing.T) {
	quantity := uint(1001)

	limit := Rate{
		Interval: time.Second,
		Quantity: 1000,
	}

	disciplined := testDiscipline(t, quantity, limit, false)
	undisciplined := testUndisciplined(t, quantity)

	require.Less(t, undisciplined, disciplined)
}

func testDiscipline(t *testing.T, quantity uint, limit Rate, optimize bool) time.Duration {
	if optimize {
		optimized, err := limit.Optimize()
		require.NoError(t, err)

		limit = optimized
	}

	expectedDuration, expectedDeviation := calcExpectedDuration(quantity, limit, 0.1)

	input := make(chan uint, int(float64(quantity)*defaultCapacityFactor))

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

	require.Equal(t, inSequence, outSequence)
	require.InDelta(t, expectedDuration, duration, expectedDeviation)

	return duration
}

func calcExpectedDuration(
	quantity uint,
	limit Rate,
	relativeDeviation float64,
) (time.Duration, float64) {
	ticksAmount := uint64(quantity) / limit.Quantity

	if ticksAmount*limit.Quantity < uint64(quantity) {
		ticksAmount++
	}

	duration := time.Duration(ticksAmount) * limit.Interval
	deviation := relativeDeviation * float64(duration)

	return duration, deviation
}

func testUndisciplined(t *testing.T, quantity uint) time.Duration {
	input := make(chan uint)

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
	quantity := uint(100000)

	limit := Rate{
		Interval: time.Second,
		Quantity: uint64(quantity),
	}

	benchmarkDiscipline(b, quantity, limit, false)
}

func BenchmarkDisciplineOptimize(b *testing.B) {
	quantity := uint(1000000)

	limit := Rate{
		Interval: time.Second,
		Quantity: uint64(quantity),
	}

	benchmarkDiscipline(b, quantity, limit, true)
}
