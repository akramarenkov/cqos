package limit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptsValidation(t *testing.T) {
	opts := Opts[uint]{}

	_, err := New(opts)
	require.Error(t, err)
}

func TestDiscipline(t *testing.T) {
	quantity := uint(10000)

	limit := Rate{
		Interval: time.Second,
		Quantity: 1000,
	}

	testDiscipline(t, quantity, limit, true)
}

func TestDisciplineAlignedQuantity(t *testing.T) {
	quantity := uint(1000)

	limit := Rate{
		Interval: time.Second,
		Quantity: uint64(quantity),
	}

	testDiscipline(t, quantity, limit, false)
}

func testDiscipline(t *testing.T, quantity uint, limit Rate, optimize bool) {
	expectedDuration := (time.Duration(quantity) * limit.Interval) / time.Duration(limit.Quantity)
	expectedDeviation := float64(expectedDuration) * 0.1

	if optimize {
		optimized, err := limit.Optimize()
		require.NoError(t, err)

		limit = optimized
	}

	input := make(chan uint)

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

	require.Equal(t, inSequence, outSequence)
	require.InDelta(t, expectedDuration, time.Since(startedAt), expectedDeviation)
}
