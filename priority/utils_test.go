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
