package stack

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscipline(t *testing.T) {
	quantity := 100

	input := make(chan uint)

	opts := Opts[uint]{
		Input: input,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	wg := &sync.WaitGroup{}

	wg.Add(2)

	inSequence := make([]uint, 0, quantity)

	go func() {
		defer wg.Done()
		defer close(input)

		period := uint(1)

		for stage := 0; stage < quantity; stage++ {
			inSequence = append(inSequence, period)

			input <- period

			period++

			if period == 4 {
				period = 1
			}
		}
	}()

	outSequence := make([][]uint, 0, quantity/10)

	go func() {
		defer wg.Done()

		for item := range discipline.Output() {
			outSequence = append(outSequence, item)
		}
	}()

	wg.Wait()

	fmt.Println(inSequence)
	fmt.Println(outSequence)
}
