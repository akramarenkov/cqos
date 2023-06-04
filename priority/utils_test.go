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

func TestIsDistributionFilled(t *testing.T) {
	require.Equal(t, false, isDistributionFilled(map[uint]uint{3: 0, 2: 0, 1: 0}))
	require.Equal(t, false, isDistributionFilled(map[uint]uint{3: 1, 2: 0, 1: 0}))
	require.Equal(t, false, isDistributionFilled(map[uint]uint{3: 0, 2: 1, 1: 0}))
	require.Equal(t, false, isDistributionFilled(map[uint]uint{3: 0, 2: 0, 1: 1}))
	require.Equal(t, false, isDistributionFilled(map[uint]uint{3: 1, 2: 1, 1: 0}))
	require.Equal(t, false, isDistributionFilled(map[uint]uint{3: 1, 2: 0, 1: 1}))
	require.Equal(t, false, isDistributionFilled(map[uint]uint{3: 0, 2: 1, 1: 1}))
	require.Equal(t, true, isDistributionFilled(map[uint]uint{3: 1, 2: 1, 1: 1}))
}

func TestGenPriorityCombinations1(t *testing.T) {
	expected := [][]uint{
		{
			1,
		},
	}

	actual := genPriorityCombinations([]uint{1})
	require.ElementsMatch(t, expected, actual)
}

func TestGenPriorityCombinations21(t *testing.T) {
	expected := [][]uint{
		{
			2,
		},
		{
			1,
		},
		{
			2, 1,
		},
	}

	actual := genPriorityCombinations([]uint{2, 1})
	require.ElementsMatch(t, expected, actual)
}

func TestGenPriorityCombinations321(t *testing.T) {
	expected := [][]uint{
		{
			3,
		},
		{
			2,
		},
		{
			1,
		},
		{
			2, 1,
		},
		{
			3, 2,
		},
		{
			3, 1,
		},
		{
			3, 2, 1,
		},
	}

	actual := genPriorityCombinations([]uint{3, 2, 1})
	require.ElementsMatch(t, expected, actual)
}

func TestGenPriorityCombinations4321(t *testing.T) {
	expected := [][]uint{
		{
			4,
		},
		{
			3,
		},
		{
			2,
		},
		{
			1,
		},
		{
			2, 1,
		},
		{
			3, 2,
		},
		{
			3, 1,
		},
		{
			4, 3,
		},
		{
			4, 2,
		},
		{
			4, 1,
		},
		{
			3, 2, 1,
		},
		{
			4, 3, 1,
		},
		{
			4, 2, 1,
		},
		{
			4, 3, 2,
		},
		{
			4, 3, 2, 1,
		},
	}

	actual := genPriorityCombinations([]uint{4, 3, 2, 1})
	require.ElementsMatch(t, expected, actual)
}

func TestGenPriorityCombinations702010(t *testing.T) {
	expected := [][]uint{
		{
			70,
		},
		{
			20,
		},
		{
			10,
		},
		{
			20, 10,
		},
		{
			70, 20,
		},
		{
			70, 10,
		},
		{
			70, 20, 10,
		},
	}

	actual := genPriorityCombinations([]uint{70, 20, 10})
	require.ElementsMatch(t, expected, actual)
}

func TestGenPriorityCombinations17(t *testing.T) {
	quantity := 17
	prioritites := make([]uint, 0, quantity)

	for id := 1; id <= quantity; id++ {
		prioritites = append(prioritites, uint(id))
	}

	actual := genPriorityCombinations(prioritites)
	require.NotEqual(t, 0, len(actual))
}

func BenchmarkGenPriorityCombinations(b *testing.B) {
	quantity := 17
	prioritites := make([]uint, 0, quantity)

	for id := 1; id <= quantity; id++ {
		prioritites = append(prioritites, uint(id))
	}

	b.ResetTimer()

	_ = genPriorityCombinations(prioritites)
}

func TestIsNonFatalConfig(t *testing.T) {
	require.Equal(t, false, IsNonFatalConfig([]uint{1}, FairDivider, 0))
	require.Equal(t, true, IsNonFatalConfig([]uint{1}, FairDivider, 1))
	require.Equal(t, true, IsNonFatalConfig([]uint{1}, FairDivider, 2))
	require.Equal(t, true, IsNonFatalConfig([]uint{1}, FairDivider, 3))

	require.Equal(t, false, IsNonFatalConfig([]uint{2, 1}, FairDivider, 0))
	require.Equal(t, false, IsNonFatalConfig([]uint{2, 1}, FairDivider, 1))
	require.Equal(t, true, IsNonFatalConfig([]uint{2, 1}, FairDivider, 2))
	require.Equal(t, true, IsNonFatalConfig([]uint{2, 1}, FairDivider, 3))
	require.Equal(t, true, IsNonFatalConfig([]uint{2, 1}, FairDivider, 4))
	require.Equal(t, true, IsNonFatalConfig([]uint{2, 1}, FairDivider, 5))

	require.Equal(t, false, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 0))
	require.Equal(t, false, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 1))
	require.Equal(t, false, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 2))
	require.Equal(t, true, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 3))
	require.Equal(t, true, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 4))
	require.Equal(t, true, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 5))
	require.Equal(t, true, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 6))
	require.Equal(t, true, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 7))

	require.Equal(t, false, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 0))
	require.Equal(t, false, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 1))
	require.Equal(t, false, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 2))
	require.Equal(t, false, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 3))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 4))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 5))
	require.Equal(t, false, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 6))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 7))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 8))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 9))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 10))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 11))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 12))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 13))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 14))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 15))

	require.Equal(t, false, IsNonFatalConfig([]uint{3, 2, 1, 0}, RateDivider, 100))

	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 0))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 1))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 2))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 3))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 4))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 5))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 6))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 7))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 8))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 9))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 10))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 11))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 12))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 13))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 14))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 15))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 16))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 17))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 18))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 19))
}

func TestIsSmallErrorConfig(t *testing.T) {
	require.Equal(t, false, IsSmallErrorConfig([]uint{1}, FairDivider, 0, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{1}, FairDivider, 1, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{1}, FairDivider, 2, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{1}, FairDivider, 3, 10.0))

	require.Equal(t, false, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 0, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 1, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 2, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 3, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 4, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 5, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 6, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 7, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 8, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 9, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 30, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 50, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 70, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{2, 1}, FairDivider, 90, 10.0))

	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 0, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 1, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 2, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 3, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 4, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 5, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 6, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 7, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 8, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 9, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 10, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 11, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{3, 2, 1}, FairDivider, 12, 10.0))

	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 0, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 1, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 2, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 3, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 4, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 5, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 6, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 7, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 8, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 9, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 10, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 11, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 12, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 13, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 14, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{4, 3, 2, 1}, FairDivider, 15, 10.0))

	require.Equal(t, false, IsSmallErrorConfig([]uint{3, 2, 1, 0}, RateDivider, 100, 10.0))

	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 0, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 1, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 2, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 3, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 4, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 5, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 6, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 7, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 8, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 9, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 10, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 11, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 12, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 13, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 14, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 15, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 16, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 17, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 18, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 19, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 20, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 21, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 22, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 23, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 24, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 25, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 26, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 27, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 28, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 29, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 30, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 31, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 32, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 33, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 34, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 35, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 36, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 37, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 38, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 39, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 40, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 41, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 42, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 43, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 44, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 45, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 46, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 47, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 48, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 49, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 50, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 51, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 52, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 53, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 54, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 55, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 56, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 57, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 58, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 59, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 60, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 61, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 62, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 63, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 64, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 65, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 66, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 67, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 68, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 69, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 70, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 71, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 72, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 73, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 74, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 75, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 76, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 77, 10.0))
	require.Equal(t, false, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 78, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 79, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 80, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 81, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 82, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 83, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 84, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 85, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 86, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 87, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 88, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 89, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 90, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 91, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 92, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 93, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 94, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 95, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 96, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 97, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 98, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 99, 10.0))
	require.Equal(t, true, IsSmallErrorConfig([]uint{70, 20, 10}, RateDivider, 100, 10.0))
}
