package durations

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateCopy(t *testing.T) {
	durations := []time.Duration{5, 4, 3, 2, 1, 0}

	copied := createCopy(durations)
	require.Equal(t, durations, copied)
	require.NotSame(t, durations, copied)

	Sort(copied)

	require.NotEqual(t, durations, copied)
	require.ElementsMatch(t, durations, copied)
}

func TestSort(t *testing.T) {
	durations := []time.Duration{5, 4, 3, 2, 1, 0}

	Sort(durations)

	require.Equal(t, []time.Duration{0, 1, 2, 3, 4, 5}, durations)
}

func TestIsSorted(t *testing.T) {
	durations := []time.Duration{5, 4, 3, 2, 1, 0}

	require.False(t, IsSorted(durations))

	Sort(durations)

	require.True(t, IsSorted(durations))
	require.True(t, IsSorted(nil))
	require.True(t, IsSorted([]time.Duration{}))
	require.True(t, IsSorted([]time.Duration{5}))
}

func TestCalcTotalDuration(t *testing.T) {
	require.Equal(t, time.Duration(0), CalcTotalDuration(nil))
	require.Equal(t, time.Duration(0), CalcTotalDuration([]time.Duration{}))
	require.Equal(t, time.Duration(5), CalcTotalDuration([]time.Duration{5, 4, 3, 2, 1, 0}))
	require.Equal(t, time.Duration(5), CalcTotalDuration([]time.Duration{5, 4, 3, 2, 1}))
	require.Equal(t, time.Duration(5), CalcTotalDuration([]time.Duration{5, 4, 3, 2}))
}
