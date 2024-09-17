package limit_test

import (
	"fmt"
	"time"

	"github.com/akramarenkov/cqos/v2/limit"
)

func ExampleDiscipline() {
	quantity := 10

	// Preferably input channel should be buffered for performance reasons.
	// Optimal capacity is in the range of 1e2 to 1e6 and should be determined
	// using benchmarks
	input := make(chan int, 10)

	opts := limit.Opts[int]{
		Input: input,
		Limit: limit.Rate{
			Interval: time.Second,
			Quantity: 1,
		},
	}

	discipline, err := limit.New(opts)
	if err != nil {
		panic(err)
	}

	outSequence := make([]int, 0, quantity)

	startedAt := time.Now()

	go func() {
		defer close(input)

		for item := 1; item <= quantity; item++ {
			input <- item
		}
	}()

	for item := range discipline.Output() {
		outSequence = append(outSequence, item)
	}

	duration := time.Since(startedAt)
	expected := (time.Duration(quantity) / time.Duration(opts.Limit.Quantity)) * opts.Limit.Interval
	deviation := 0.01

	fmt.Println(duration <= time.Duration(float64(expected)*(1.0+deviation)))
	fmt.Println(duration >= time.Duration(float64(expected)*(1.0-deviation)))
	fmt.Println(outSequence)

	// Output:
	// true
	// true
	// [1 2 3 4 5 6 7 8 9 10]
}
