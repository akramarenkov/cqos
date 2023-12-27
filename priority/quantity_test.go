package priority

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatePriorities(t *testing.T) {
	require.Equal(t, []uint{}, createPriorities(0))
	require.Equal(t, []uint{1}, createPriorities(1))
	require.Equal(t, []uint{2, 1}, createPriorities(2))
	require.Equal(t, []uint{3, 2, 1}, createPriorities(3))
	require.Equal(t, []uint{4, 3, 2, 1}, createPriorities(4))
	require.Equal(t, []uint{5, 4, 3, 2, 1}, createPriorities(5))
	require.Equal(t, []uint{6, 5, 4, 3, 2, 1}, createPriorities(6))
	require.Equal(t, []uint{7, 6, 5, 4, 3, 2, 1}, createPriorities(7))
	require.Equal(t, []uint{8, 7, 6, 5, 4, 3, 2, 1}, createPriorities(8))
	require.Equal(t, []uint{9, 8, 7, 6, 5, 4, 3, 2, 1}, createPriorities(9))
	require.Equal(t, []uint{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, createPriorities(10))
}

func TestCalcCombinationsQuantitySlow(t *testing.T) {
	require.Equal(t, 0, calcCombinationsQuantitySlow([]uint{}))
	require.Equal(t, 1, calcCombinationsQuantitySlow([]uint{1}))
	require.Equal(t, 3, calcCombinationsQuantitySlow([]uint{2, 1}))
	require.Equal(t, 7, calcCombinationsQuantitySlow([]uint{3, 2, 1}))
	require.Equal(t, 15, calcCombinationsQuantitySlow([]uint{4, 3, 2, 1}))
	require.Equal(t, 31, calcCombinationsQuantitySlow([]uint{5, 4, 3, 2, 1}))
	require.Equal(t, 63, calcCombinationsQuantitySlow([]uint{6, 5, 4, 3, 2, 1}))
	require.Equal(t, 127, calcCombinationsQuantitySlow([]uint{7, 6, 5, 4, 3, 2, 1}))
	require.Equal(t, 255, calcCombinationsQuantitySlow([]uint{8, 7, 6, 5, 4, 3, 2, 1}))
	require.Equal(t, 511, calcCombinationsQuantitySlow([]uint{9, 8, 7, 6, 5, 4, 3, 2, 1}))
	require.Equal(t, 1023, calcCombinationsQuantitySlow([]uint{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}))
}

func TestCalcCombinationsQuantity(t *testing.T) {
	require.Equal(t, 0, calcCombinationsQuantity([]uint{}))
	require.Equal(t, 1, calcCombinationsQuantity([]uint{1}))
	require.Equal(t, 3, calcCombinationsQuantity([]uint{2, 1}))
	require.Equal(t, 7, calcCombinationsQuantity([]uint{3, 2, 1}))
	require.Equal(t, 15, calcCombinationsQuantity([]uint{4, 3, 2, 1}))
	require.Equal(t, 31, calcCombinationsQuantity([]uint{5, 4, 3, 2, 1}))
	require.Equal(t, 63, calcCombinationsQuantity([]uint{6, 5, 4, 3, 2, 1}))
	require.Equal(t, 127, calcCombinationsQuantity([]uint{7, 6, 5, 4, 3, 2, 1}))
	require.Equal(t, 255, calcCombinationsQuantity([]uint{8, 7, 6, 5, 4, 3, 2, 1}))
	require.Equal(t, 511, calcCombinationsQuantity([]uint{9, 8, 7, 6, 5, 4, 3, 2, 1}))
	require.Equal(t, 1023, calcCombinationsQuantity([]uint{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}))
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

func TestGenPriorityCombinations54321(t *testing.T) {
	expected := [][]uint{
		{
			5,
		},
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
			5, 4,
		},
		{
			5, 3,
		},
		{
			5, 2,
		},
		{
			5, 1,
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
		{
			5, 2, 1,
		},
		{
			5, 3, 2,
		},
		{
			5, 3, 1,
		},
		{
			5, 4, 3,
		},
		{
			5, 4, 2,
		},
		{
			5, 4, 1,
		},
		{
			5, 3, 2, 1,
		},
		{
			5, 4, 3, 1,
		},
		{
			5, 4, 2, 1,
		},
		{
			5, 4, 3, 2,
		},
		{
			5, 4, 3, 2, 1,
		},
	}

	actual := genPriorityCombinations([]uint{5, 4, 3, 2, 1})
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

func TestGenPriorityCombinations(t *testing.T) {
	priorities := createPriorities(15)

	actual := genPriorityCombinations(priorities)
	require.Len(t, actual, calcCombinationsQuantity(priorities))
}

func BenchmarkGenPriorityCombinations(b *testing.B) {
	priorities := createPriorities(20)

	b.ResetTimer()

	_ = genPriorityCombinations(priorities)
}

func TestCreateSortedCopy(t *testing.T) {
	priorities := []uint{1, 2, 3, 4, 5}

	copied := createSortedCopy(priorities)
	require.NotSame(t, priorities, copied)
	require.Equal(t, []uint{5, 4, 3, 2, 1}, copied)
}

func TestIsNonFatalConfig(t *testing.T) {
	require.False(t, IsNonFatalConfig([]uint{1}, FairDivider, 0))
	require.True(t, IsNonFatalConfig([]uint{1}, FairDivider, 1))
	require.True(t, IsNonFatalConfig([]uint{1}, FairDivider, 2))
	require.True(t, IsNonFatalConfig([]uint{1}, FairDivider, 3))

	require.False(t, IsNonFatalConfig([]uint{2, 1}, FairDivider, 0))
	require.False(t, IsNonFatalConfig([]uint{2, 1}, FairDivider, 1))
	require.True(t, IsNonFatalConfig([]uint{2, 1}, FairDivider, 2))
	require.True(t, IsNonFatalConfig([]uint{2, 1}, FairDivider, 3))
	require.True(t, IsNonFatalConfig([]uint{2, 1}, FairDivider, 4))
	require.True(t, IsNonFatalConfig([]uint{2, 1}, FairDivider, 5))

	require.False(t, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 0))
	require.False(t, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 1))
	require.False(t, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 2))
	require.True(t, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 3))
	require.True(t, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 4))
	require.True(t, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 5))
	require.True(t, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 6))
	require.True(t, IsNonFatalConfig([]uint{3, 2, 1}, FairDivider, 7))

	require.False(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 0))
	require.False(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 1))
	require.False(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 2))
	require.False(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 3))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 4))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 5))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 6))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 7))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 8))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 9))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 10))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 11))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 12))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 13))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 14))
	require.True(t, IsNonFatalConfig([]uint{4, 3, 2, 1}, FairDivider, 15))

	require.False(t, IsNonFatalConfig([]uint{3, 2, 1, 0}, RateDivider, 100))

	require.False(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 0))
	require.False(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 1))
	require.False(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 2))
	require.False(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 3))
	require.False(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 4))
	require.False(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 5))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 6))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 7))
	require.False(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 8))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 9))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 10))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 11))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 12))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 13))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 14))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 15))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 16))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 17))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 18))
	require.True(t, IsNonFatalConfig([]uint{70, 20, 10}, RateDivider, 19))
}

func TestPickUpMinNonFatalQuantity(t *testing.T) {
	quantity := PickUpMinNonFatalQuantity([]uint{70, 20, 10}, FairDivider, 2)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMinNonFatalQuantity([]uint{70, 20, 10}, FairDivider, 6)
	require.Equal(t, uint(3), quantity)

	quantity = PickUpMinNonFatalQuantity([]uint{70, 20, 10}, RateDivider, 8)
	require.Equal(t, uint(6), quantity)
}

func TestPickUpMaxNonFatalQuantity(t *testing.T) {
	quantity := PickUpMaxNonFatalQuantity([]uint{70, 20, 10}, FairDivider, 2)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMaxNonFatalQuantity([]uint{70, 20, 10}, FairDivider, 6)
	require.Equal(t, uint(6), quantity)

	quantity = PickUpMaxNonFatalQuantity([]uint{70, 20, 10}, RateDivider, 8)
	require.Equal(t, uint(7), quantity)
}

func TestIsDistributionSuitable(t *testing.T) {
	quantity := uint(100)

	distribution := make(map[uint]uint)
	reference := make(map[uint]uint)

	RateDivider([]uint{3, 2, 0}, quantity, distribution)
	RateDivider([]uint{3, 2, 0}, referenceFactor*quantity, reference)

	suitable := isDistributionSuitable(distribution, reference, quantity, referenceFactor*quantity, 10.0)
	require.False(t, suitable)

	distribution = make(map[uint]uint)
	reference = make(map[uint]uint)

	RateDivider([]uint{3, 2, 1}, quantity, nil)
	RateDivider([]uint{3, 2, 1}, referenceFactor*quantity, nil)

	suitable = isDistributionSuitable(distribution, reference, quantity, referenceFactor*quantity, 10.0)
	require.True(t, suitable)
}

func TestIsSuitableConfig(t *testing.T) {
	require.False(t, IsSuitableConfig([]uint{1}, FairDivider, 0, 10.0))
	require.True(t, IsSuitableConfig([]uint{1}, FairDivider, 1, 10.0))
	require.True(t, IsSuitableConfig([]uint{1}, FairDivider, 2, 10.0))
	require.True(t, IsSuitableConfig([]uint{1}, FairDivider, 3, 10.0))

	require.False(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 0, 10.0))
	require.False(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 1, 10.0))
	require.True(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 2, 10.0))
	require.False(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 3, 10.0))
	require.True(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 4, 10.0))
	require.False(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 5, 10.0))
	require.True(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 6, 10.0))
	require.False(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 7, 10.0))
	require.True(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 8, 10.0))
	require.False(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 9, 10.0))
	require.True(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 30, 10.0))
	require.True(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 50, 10.0))
	require.True(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 70, 10.0))
	require.True(t, IsSuitableConfig([]uint{2, 1}, FairDivider, 90, 10.0))

	require.False(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 0, 10.0))
	require.False(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 1, 10.0))
	require.False(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 2, 10.0))
	require.False(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 3, 10.0))
	require.False(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 4, 10.0))
	require.False(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 5, 10.0))
	require.True(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 6, 10.0))
	require.False(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 7, 10.0))
	require.False(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 8, 10.0))
	require.False(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 9, 10.0))
	require.False(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 10, 10.0))
	require.False(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 11, 10.0))
	require.True(t, IsSuitableConfig([]uint{3, 2, 1}, FairDivider, 12, 10.0))

	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 0, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 1, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 2, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 3, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 4, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 5, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 6, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 7, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 8, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 9, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 10, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 11, 10.0))
	require.True(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 12, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 13, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 14, 10.0))
	require.False(t, IsSuitableConfig([]uint{4, 3, 2, 1}, FairDivider, 15, 10.0))

	require.False(t, IsSuitableConfig([]uint{3, 2, 1, 0}, RateDivider, 100, 10.0))

	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 0, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 1, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 2, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 3, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 4, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 5, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 6, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 7, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 8, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 9, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 10, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 11, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 12, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 13, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 14, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 15, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 16, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 17, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 18, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 19, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 20, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 21, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 22, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 23, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 24, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 25, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 26, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 27, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 28, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 29, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 30, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 31, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 32, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 33, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 34, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 35, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 36, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 37, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 38, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 39, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 40, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 41, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 42, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 43, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 44, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 45, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 46, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 47, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 48, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 49, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 50, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 51, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 52, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 53, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 54, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 55, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 56, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 57, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 58, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 59, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 60, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 61, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 62, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 63, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 64, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 65, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 66, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 67, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 68, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 69, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 70, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 71, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 72, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 73, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 74, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 75, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 76, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 77, 10.0))
	require.False(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 78, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 79, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 80, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 81, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 82, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 83, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 84, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 85, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 86, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 87, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 88, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 89, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 90, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 91, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 92, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 93, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 94, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 95, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 96, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 97, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 98, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 99, 10.0))
	require.True(t, IsSuitableConfig([]uint{70, 20, 10}, RateDivider, 100, 10.0))
}

func TestPickUpMinSuitableQuantity(t *testing.T) {
	quantity := PickUpMinSuitableQuantity([]uint{70, 20, 10}, FairDivider, 5, 10)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMinSuitableQuantity([]uint{70, 20, 10}, FairDivider, 1000, 10)
	require.Equal(t, uint(6), quantity)

	quantity = PickUpMinSuitableQuantity([]uint{70, 20, 10}, RateDivider, 21, 10)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMinSuitableQuantity([]uint{70, 20, 10}, RateDivider, 1000, 10)
	require.Equal(t, uint(22), quantity)
}

func TestPickUpMaxSuitableQuantity(t *testing.T) {
	quantity := PickUpMaxSuitableQuantity([]uint{70, 20, 10}, FairDivider, 5, 10)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMaxSuitableQuantity([]uint{70, 20, 10}, FairDivider, 7, 10)
	require.Equal(t, uint(6), quantity)

	quantity = PickUpMaxSuitableQuantity([]uint{70, 20, 10}, FairDivider, 1000, 10)
	require.Equal(t, uint(1000), quantity)

	quantity = PickUpMaxSuitableQuantity([]uint{70, 20, 10}, RateDivider, 21, 10)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMaxSuitableQuantity([]uint{70, 20, 10}, RateDivider, 23, 10)
	require.Equal(t, uint(22), quantity)

	quantity = PickUpMaxSuitableQuantity([]uint{70, 20, 10}, RateDivider, 1000, 10)
	require.Equal(t, uint(1000), quantity)
}
