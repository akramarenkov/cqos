package join_test

import (
	"fmt"
	"sync"

	"github.com/akramarenkov/cqos/v2/join"
)

func ExampleDiscipline() {
	quantity := 27

	input := make(chan uint)

	opts := join.Opts[uint]{
		Input:    input,
		JoinSize: 5,
	}

	discipline, err := join.New(opts)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		defer close(input)

		for stage := 1; stage <= quantity; stage++ {
			input <- uint(stage)
		}
	}()

	outSequence := make([]uint, 0, quantity)

	go func() {
		defer wg.Done()

		for slice := range discipline.Output() {
			outSequence = append(outSequence, slice...)
		}
	}()

	wg.Wait()

	fmt.Println(outSequence)
	// Output:[1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27]
}
