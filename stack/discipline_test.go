package stack

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscipline(t *testing.T) {
	quantity := 105

	input := make(chan uint)

	opts := Opts[uint]{
		Input:     input,
		StackSize: 10,
	}

	discipline, err := New(opts)
	require.NoError(t, err)

	wg := &sync.WaitGroup{}

	wg.Add(2)

	inSequence := make([]uint, 0, quantity)

	go func() {
		defer wg.Done()
		defer close(input)

		for stage := 1; stage <= quantity; stage++ {
			inSequence = append(inSequence, uint(stage))

			input <- uint(stage)
		}
	}()

	outSequence := make([]uint, 0, quantity)

	go func() {
		defer wg.Done()

		for stack := range discipline.Output() {
			require.NotEqual(t, 0, stack)
			outSequence = append(outSequence, stack...)
		}
	}()

	wg.Wait()

	require.Equal(t, inSequence, outSequence)
}
