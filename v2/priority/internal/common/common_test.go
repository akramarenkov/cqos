package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortPriorities(t *testing.T) {
	priorities := []uint{2, 1, 3}
	expected := []uint{3, 2, 1}

	SortPriorities(priorities)
	require.Equal(t, expected, priorities)
}

func TestSumPriorities(t *testing.T) {
	require.Equal(t, uint(1), SumPriorities([]uint{1}))
	require.Equal(t, uint(2), SumPriorities([]uint{2}))
	require.Equal(t, uint(3), SumPriorities([]uint{3}))
	require.Equal(t, uint(4), SumPriorities([]uint{4}))
	require.Equal(t, uint(3), SumPriorities([]uint{2, 1}))
	require.Equal(t, uint(5), SumPriorities([]uint{3, 2}))
	require.Equal(t, uint(4), SumPriorities([]uint{3, 1}))
	require.Equal(t, uint(7), SumPriorities([]uint{4, 3}))
	require.Equal(t, uint(6), SumPriorities([]uint{4, 2}))
	require.Equal(t, uint(5), SumPriorities([]uint{4, 1}))
	require.Equal(t, uint(6), SumPriorities([]uint{3, 2, 1}))
	require.Equal(t, uint(8), SumPriorities([]uint{4, 3, 1}))
	require.Equal(t, uint(7), SumPriorities([]uint{4, 2, 1}))
	require.Equal(t, uint(9), SumPriorities([]uint{4, 3, 2}))
	require.Equal(t, uint(10), SumPriorities([]uint{4, 3, 2, 1}))
}

func TestCalcCapacity(t *testing.T) {
	require.Equal(t, 1, CalcCapacity(10, 0.1, 100))
	require.Equal(t, 100, CalcCapacity(10, 0.01, 100))
	require.Equal(t, 2, CalcCapacity(3, 0.5, 100))
	require.Equal(t, 1, CalcCapacity(4, 0.333, 100))
	require.Equal(t, 100, CalcCapacity(1, 0.166, 100))
}
