package unite_test

import (
	"fmt"
	"time"

	"github.com/akramarenkov/cqos/v2/join/unite"
)

func ExampleDiscipline() {
	data := [][]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
		{17, 18, 19, 20},
		{21, 22, 23, 24},
		{25, 26, 27},
	}

	// Preferably input channel should be buffered for performance reasons.
	// Optimal capacity is in the range of one to three JoinSize
	input := make(chan []int, 10)

	opts := unite.Opts[int]{
		Input:    input,
		JoinSize: 10,
		Timeout:  time.Second,
	}

	discipline, err := unite.New(opts)
	if err != nil {
		panic(err)
	}

	go func() {
		defer close(input)

		for _, item := range data {
			input <- item
		}
	}()

	for join := range discipline.Output() {
		fmt.Println(join)
	}

	// Output:
	// [1 2 3 4 5 6 7 8]
	// [9 10 11 12 13 14 15 16]
	// [17 18 19 20 21 22 23 24]
	// [25 26 27]
}

func ExampleDiscipline_Release() {
	data := [][]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
		{17, 18, 19, 20},
		{21, 22, 23, 24},
		{25, 26, 27},
	}

	// Preferably input channel should be buffered for performance reasons.
	// Optimal capacity is in the range of one to three JoinSize
	input := make(chan []int, 10)

	opts := unite.Opts[int]{
		Input:    input,
		JoinSize: 10,
		NoCopy:   true,
		Timeout:  time.Second,
	}

	discipline, err := unite.New(opts)
	if err != nil {
		panic(err)
	}

	go func() {
		defer close(input)

		for _, item := range data {
			input <- item
		}
	}()

	for join := range discipline.Output() {
		fmt.Println(join)

		discipline.Release()
	}

	// Output:
	// [1 2 3 4 5 6 7 8]
	// [9 10 11 12 13 14 15 16]
	// [17 18 19 20 21 22 23 24]
	// [25 26 27]
}
