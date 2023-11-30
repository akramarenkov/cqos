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

func TestCalcByFactor(t *testing.T) {
	require.Equal(t, 3, CalcByFactor(10, 0.1, 3))
	require.Equal(t, 1, CalcByFactor(10, 0.1, 0))
	require.Equal(t, 10, CalcByFactor(100, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(14, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(15, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(16, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(24, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(25, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(26, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(34, 0.1, 3))
	require.Equal(t, 4, CalcByFactor(35, 0.1, 3))
	require.Equal(t, 4, CalcByFactor(36, 0.1, 3))
	require.Equal(t, 0, CalcByFactor(4, 0.1, 0))
	require.Equal(t, 1, CalcByFactor(5, 0.1, 0))
	require.Equal(t, 1, CalcByFactor(6, 0.1, 0))
}
