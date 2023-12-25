package limit_test

import (
	"fmt"
	"time"

	"github.com/akramarenkov/cqos/v2/limit"
)

func ExampleDiscipline() {
	quantity := 10

	input := make(chan int)

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

		for stage := 1; stage <= quantity; stage++ {
			input <- stage
		}
	}()

	for item := range discipline.Output() {
		outSequence = append(outSequence, item)
	}

	duration := time.Since(startedAt)
	expected := (time.Duration(quantity) * opts.Limit.Interval) / time.Duration(opts.Limit.Quantity)
	deviation := 0.01

	fmt.Println(duration <= time.Duration(float64(expected)*(1.0+deviation)))
	fmt.Println(duration >= time.Duration(float64(expected)*(1.0-deviation)))
	fmt.Println(outSequence)
	// Output:
	// true
	// true
	// [1 2 3 4 5 6 7 8 9 10]
}
