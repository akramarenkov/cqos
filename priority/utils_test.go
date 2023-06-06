package priority

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortPriorities(t *testing.T) {
	priorities := []uint{2, 1, 3}
	expected := []uint{3, 2, 1}

	sortPriorities(priorities)
	require.Equal(t, expected, priorities)
}

func TestRemovePriority(t *testing.T) {
	priorities := []uint{4, 3, 2, 1}

	priorities = removePriority(priorities, 5)
	require.Equal(t, priorities, priorities)

	priorities = removePriority(priorities, 2)
	require.Equal(t, []uint{4, 3, 1}, priorities)

	priorities = removePriority(priorities, 4)
	require.Equal(t, []uint{3, 1}, priorities)

	priorities = removePriority(priorities, 1)
	require.Equal(t, []uint{3}, priorities)

	priorities = removePriority(priorities, 3)
	require.Equal(t, []uint{}, priorities)

	priorities = removePriority(priorities, 3)
	require.Equal(t, []uint{}, priorities)
}

func TestSumPriorities(t *testing.T) {
	require.Equal(t, uint(1), sumPriorities([]uint{1}))
	require.Equal(t, uint(2), sumPriorities([]uint{2}))
	require.Equal(t, uint(3), sumPriorities([]uint{3}))
	require.Equal(t, uint(4), sumPriorities([]uint{4}))
	require.Equal(t, uint(3), sumPriorities([]uint{2, 1}))
	require.Equal(t, uint(5), sumPriorities([]uint{3, 2}))
	require.Equal(t, uint(4), sumPriorities([]uint{3, 1}))
	require.Equal(t, uint(7), sumPriorities([]uint{4, 3}))
	require.Equal(t, uint(6), sumPriorities([]uint{4, 2}))
	require.Equal(t, uint(5), sumPriorities([]uint{4, 1}))
	require.Equal(t, uint(6), sumPriorities([]uint{3, 2, 1}))
	require.Equal(t, uint(8), sumPriorities([]uint{4, 3, 1}))
	require.Equal(t, uint(7), sumPriorities([]uint{4, 2, 1}))
	require.Equal(t, uint(9), sumPriorities([]uint{4, 3, 2}))
	require.Equal(t, uint(10), sumPriorities([]uint{4, 3, 2, 1}))
}
