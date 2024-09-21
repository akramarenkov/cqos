package limit

import (
	"testing"
	"time"

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
		Limit: Rate{
			Interval: time.Second,
			Quantity: 1,
		},
	}

	_, err = New(opts)
	require.NoError(t, err)
}

func TestDiscipline(t *testing.T) {
	quantity := 10000

	limit := Rate{
		Interval: time.Second,
		Quantity: 1000,
	}

	duration := testDiscipline(t, quantity, limit)
	expected := calcExpectedDuration(quantity, limit)
	require.InEpsilon(t, expected, duration, 0.1)
}

func testDiscipline(t *testing.T, quantity int, limit Rate) time.Duration {
	input := make(chan int, quantity)

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

		for item := range quantity {
			inSequence = append(inSequence, item)

			input <- item
		}
	}()

	for item := range discipline.Output() {
		outSequence = append(outSequence, item)
	}

	duration := time.Since(startedAt)

	require.Equal(t, inSequence, outSequence)

	return duration
}

func calcExpectedDuration(quantity int, limit Rate) time.Duration {
	// Accuracy of calculations is deliberately roughened (first division is performed
	// and only then multiplication) because such a calculation corresponds to the work
	// of the discipline when closing the input channel: if the number of data elements
	// written to the input channel is not a multiple of the Quantity field in rate
	// limit structure, then the delay after the transmission of the last data is not
	// performed
	ratio := time.Duration(quantity) / time.Duration(limit.Quantity)

	return ratio * limit.Interval
}

func BenchmarkDisciplineInputCapacity0(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 0)
}

func BenchmarkDisciplineInputCapacity1e0(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1)
}

func BenchmarkDisciplineInputCapacity1e1(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e1)
}

func BenchmarkDisciplineInputCapacity1e2(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e2)
}

func BenchmarkDisciplineInputCapacity1e3(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e3)
}

func BenchmarkDisciplineInputCapacity1e4(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e4)
}

func BenchmarkDisciplineInputCapacity1e5(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e5)
}

func BenchmarkDisciplineInputCapacity1e6(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e6)
}

func BenchmarkDisciplineInputCapacity1e7(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, 1e7)
}

func BenchmarkDisciplineInputCapacityQuantity(b *testing.B) {
	benchmarkDisciplineInputCapacity(b, -1)
}

// This benchmark is used to test the impact of input channel capacity on
// performance. Therefore, the value of Quantity field in rate limit structure is
// always set to be greater than the number of data elements written to the input
// channel so that there is no delay after data transfer.
func benchmarkDisciplineInputCapacity(b *testing.B, capacity int) {
	quantity := b.N

	limit := Rate{
		Interval: time.Second,
		Quantity: uint64(b.N) + 1,
	}

	if capacity < 0 {
		capacity = quantity
	}

	input := make(chan int, capacity)

	opts := Opts[int]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	b.ResetTimer()

	go func() {
		defer close(input)

		for item := range quantity {
			input <- item
		}
	}()

	for item := range discipline.Output() {
		_ = item
	}
}

// Here we model the worst case: when the number of operations for measuring time and
// calculating delays is equal to the number of operations for transmitting data
// elements.
func BenchmarkDiscipline(b *testing.B) {
	quantity := b.N

	limit := Rate{
		Interval: time.Nanosecond,
		Quantity: 1,
	}

	input := make(chan int, b.N)

	opts := Opts[int]{
		Input: input,
		Limit: limit,
	}

	discipline, err := New(opts)
	require.NoError(b, err)

	b.ResetTimer()

	go func() {
		defer close(input)

		for item := range quantity {
			input <- item
		}
	}()

	for item := range discipline.Output() {
		_ = item
	}
}
