package priority

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestCalcCapacity(t *testing.T) {
	require.Equal(t, 1, calcCapacity(10, 0.1, 100))
	require.Equal(t, 100, calcCapacity(10, 0.01, 100))
	require.Equal(t, 2, calcCapacity(3, 0.5, 100))
	require.Equal(t, 1, calcCapacity(4, 0.333, 100))
	require.Equal(t, 100, calcCapacity(1, 0.166, 100))
}
