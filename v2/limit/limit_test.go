package limit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDiscipline(t *testing.T) {
	quantity := uint(10000)

	limit := Rate{
		Interval: time.Second,
		Quantity: 1000,
	}

	expectedDuration := (time.Duration(quantity) * limit.Interval) / time.Duration(limit.Quantity)
	expectedDeviation := float64(expectedDuration) * 0.1

	limit, err := limit.Optimize()
	require.NoError(t, err)

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

func TestDisciplineError(t *testing.T) {
	opts := Opts[uint]{}

	_, err := New(opts)
	require.Error(t, err)
}

func TestDisciplineAligned(t *testing.T) {
	quantity := uint(1000)

	limit := Rate{
		Interval: time.Second,
		Quantity: 1000,
	}

	expectedDuration := (time.Duration(quantity) * limit.Interval) / time.Duration(limit.Quantity)
	expectedDeviation := float64(expectedDuration) * 0.1

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
