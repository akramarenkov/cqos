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
	output := make(chan []uint)

	opts := Opts[uint]{
		Input:  input,
		Output: output,
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
		defer close(output)
		defer discipline.GracefulStop()

		for item := range output {
			outSequence = append(outSequence, item)

			if len(outSequence) == cap(outSequence) {
				return
			}
		}
	}()

	wg.Wait()

	fmt.Println(inSequence)
	fmt.Println(outSequence)
}
