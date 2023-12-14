package utils

import (
	"testing"

	"github.com/akramarenkov/cqos/v2/priority/divider"

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

func TestGenCombinations1(t *testing.T) {
	expected := [][]uint{
		{
			1,
		},
	}

	actual := genCombinations([]uint{1})
	require.ElementsMatch(t, expected, actual)
}

func TestGenCombinations21(t *testing.T) {
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

	actual := genCombinations([]uint{2, 1})
	require.ElementsMatch(t, expected, actual)
}

func TestGenCombinations321(t *testing.T) {
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

	actual := genCombinations([]uint{3, 2, 1})
	require.ElementsMatch(t, expected, actual)
}

func TestGenCombinations4321(t *testing.T) {
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

	actual := genCombinations([]uint{4, 3, 2, 1})
	require.ElementsMatch(t, expected, actual)
}

func TestGenCombinations54321(t *testing.T) {
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

	actual := genCombinations([]uint{5, 4, 3, 2, 1})
	require.ElementsMatch(t, expected, actual)
}

func TestGenCombinations702010(t *testing.T) {
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

	actual := genCombinations([]uint{70, 20, 10})
	require.ElementsMatch(t, expected, actual)
}

func TestGenCombinations(t *testing.T) {
	priorities := createPriorities(15)

	actual := genCombinations(priorities)
	require.Equal(t, calcCombinationsQuantity(priorities), len(actual))
}

func BenchmarkGenCombinations(b *testing.B) {
	priorities := createPriorities(20)

	b.ResetTimer()

	_ = genCombinations(priorities)
}

func TestCreateSortedCopy(t *testing.T) {
	priorities := []uint{1, 2, 3, 4, 5}

	copied := createSortedCopy(priorities)
	require.NotSame(t, priorities, copied)
	require.Equal(t, []uint{5, 4, 3, 2, 1}, copied)
}

func TestIsNonFatalConfig(t *testing.T) {
	require.Equal(t, false, IsNonFatalConfig([]uint{1}, divider.Fair, 0))
	require.Equal(t, true, IsNonFatalConfig([]uint{1}, divider.Fair, 1))
	require.Equal(t, true, IsNonFatalConfig([]uint{1}, divider.Fair, 2))
	require.Equal(t, true, IsNonFatalConfig([]uint{1}, divider.Fair, 3))

	require.Equal(t, false, IsNonFatalConfig([]uint{2, 1}, divider.Fair, 0))
	require.Equal(t, false, IsNonFatalConfig([]uint{2, 1}, divider.Fair, 1))
	require.Equal(t, true, IsNonFatalConfig([]uint{2, 1}, divider.Fair, 2))
	require.Equal(t, true, IsNonFatalConfig([]uint{2, 1}, divider.Fair, 3))
	require.Equal(t, true, IsNonFatalConfig([]uint{2, 1}, divider.Fair, 4))
	require.Equal(t, true, IsNonFatalConfig([]uint{2, 1}, divider.Fair, 5))

	require.Equal(t, false, IsNonFatalConfig([]uint{3, 2, 1}, divider.Fair, 0))
	require.Equal(t, false, IsNonFatalConfig([]uint{3, 2, 1}, divider.Fair, 1))
	require.Equal(t, false, IsNonFatalConfig([]uint{3, 2, 1}, divider.Fair, 2))
	require.Equal(t, true, IsNonFatalConfig([]uint{3, 2, 1}, divider.Fair, 3))
	require.Equal(t, true, IsNonFatalConfig([]uint{3, 2, 1}, divider.Fair, 4))
	require.Equal(t, true, IsNonFatalConfig([]uint{3, 2, 1}, divider.Fair, 5))
	require.Equal(t, true, IsNonFatalConfig([]uint{3, 2, 1}, divider.Fair, 6))
	require.Equal(t, true, IsNonFatalConfig([]uint{3, 2, 1}, divider.Fair, 7))

	require.Equal(t, false, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 0))
	require.Equal(t, false, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 1))
	require.Equal(t, false, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 2))
	require.Equal(t, false, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 3))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 4))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 5))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 6))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 7))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 8))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 9))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 10))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 11))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 12))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 13))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 14))
	require.Equal(t, true, IsNonFatalConfig([]uint{4, 3, 2, 1}, divider.Fair, 15))

	require.Equal(t, false, IsNonFatalConfig([]uint{3, 2, 1, 0}, divider.Rate, 100))

	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 0))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 1))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 2))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 3))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 4))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 5))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 6))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 7))
	require.Equal(t, false, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 8))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 9))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 10))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 11))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 12))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 13))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 14))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 15))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 16))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 17))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 18))
	require.Equal(t, true, IsNonFatalConfig([]uint{70, 20, 10}, divider.Rate, 19))
}

func TestPickUpMinNonFatalQuantity(t *testing.T) {
	quantity := PickUpMinNonFatalQuantity([]uint{70, 20, 10}, divider.Fair, 2)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMinNonFatalQuantity([]uint{70, 20, 10}, divider.Fair, 6)
	require.Equal(t, uint(3), quantity)

	quantity = PickUpMinNonFatalQuantity([]uint{70, 20, 10}, divider.Rate, 8)
	require.Equal(t, uint(6), quantity)
}

func TestPickUpMaxNonFatalQuantity(t *testing.T) {
	quantity := PickUpMaxNonFatalQuantity([]uint{70, 20, 10}, divider.Fair, 2)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMaxNonFatalQuantity([]uint{70, 20, 10}, divider.Fair, 6)
	require.Equal(t, uint(6), quantity)

	quantity = PickUpMaxNonFatalQuantity([]uint{70, 20, 10}, divider.Rate, 8)
	require.Equal(t, uint(7), quantity)
}

func TestIsDistributionSuitable(t *testing.T) {
	quantity := uint(100)

	distribution := make(map[uint]uint)
	reference := make(map[uint]uint)

	divider.Rate([]uint{3, 2, 0}, quantity, distribution)
	divider.Rate([]uint{3, 2, 0}, referenceFactor*quantity, reference)

	suitable := isDistributionSuitable(distribution, reference, quantity, referenceFactor*quantity, 10.0)
	require.Equal(t, false, suitable)

	distribution = make(map[uint]uint)
	reference = make(map[uint]uint)

	divider.Rate([]uint{3, 2, 1}, quantity, nil)
	divider.Rate([]uint{3, 2, 1}, referenceFactor*quantity, nil)

	suitable = isDistributionSuitable(distribution, reference, quantity, referenceFactor*quantity, 10.0)
	require.Equal(t, true, suitable)
}

func TestIsSuitableConfig(t *testing.T) {
	require.Equal(t, false, IsSuitableConfig([]uint{1}, divider.Fair, 0, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{1}, divider.Fair, 1, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{1}, divider.Fair, 2, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{1}, divider.Fair, 3, 10.0))

	require.Equal(t, false, IsSuitableConfig([]uint{2, 1}, divider.Fair, 0, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{2, 1}, divider.Fair, 1, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{2, 1}, divider.Fair, 2, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{2, 1}, divider.Fair, 3, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{2, 1}, divider.Fair, 4, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{2, 1}, divider.Fair, 5, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{2, 1}, divider.Fair, 6, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{2, 1}, divider.Fair, 7, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{2, 1}, divider.Fair, 8, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{2, 1}, divider.Fair, 9, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{2, 1}, divider.Fair, 30, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{2, 1}, divider.Fair, 50, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{2, 1}, divider.Fair, 70, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{2, 1}, divider.Fair, 90, 10.0))

	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 0, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 1, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 2, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 3, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 4, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 5, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 6, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 7, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 8, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 9, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 10, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 11, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{3, 2, 1}, divider.Fair, 12, 10.0))

	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 0, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 1, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 2, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 3, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 4, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 5, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 6, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 7, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 8, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 9, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 10, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 11, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 12, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 13, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 14, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{4, 3, 2, 1}, divider.Fair, 15, 10.0))

	require.Equal(t, false, IsSuitableConfig([]uint{3, 2, 1, 0}, divider.Rate, 100, 10.0))

	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 0, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 1, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 2, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 3, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 4, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 5, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 6, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 7, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 8, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 9, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 10, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 11, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 12, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 13, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 14, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 15, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 16, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 17, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 18, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 19, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 20, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 21, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 22, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 23, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 24, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 25, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 26, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 27, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 28, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 29, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 30, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 31, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 32, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 33, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 34, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 35, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 36, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 37, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 38, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 39, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 40, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 41, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 42, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 43, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 44, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 45, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 46, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 47, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 48, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 49, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 50, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 51, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 52, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 53, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 54, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 55, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 56, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 57, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 58, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 59, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 60, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 61, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 62, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 63, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 64, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 65, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 66, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 67, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 68, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 69, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 70, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 71, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 72, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 73, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 74, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 75, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 76, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 77, 10.0))
	require.Equal(t, false, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 78, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 79, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 80, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 81, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 82, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 83, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 84, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 85, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 86, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 87, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 88, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 89, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 90, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 91, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 92, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 93, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 94, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 95, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 96, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 97, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 98, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 99, 10.0))
	require.Equal(t, true, IsSuitableConfig([]uint{70, 20, 10}, divider.Rate, 100, 10.0))
}

func TestPickUpMinSuitableQuantity(t *testing.T) {
	quantity := PickUpMinSuitableQuantity([]uint{70, 20, 10}, divider.Fair, 5, 10)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMinSuitableQuantity([]uint{70, 20, 10}, divider.Fair, 1000, 10)
	require.Equal(t, uint(6), quantity)

	quantity = PickUpMinSuitableQuantity([]uint{70, 20, 10}, divider.Rate, 21, 10)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMinSuitableQuantity([]uint{70, 20, 10}, divider.Rate, 1000, 10)
	require.Equal(t, uint(22), quantity)
}

func TestPickUpMaxSuitableQuantity(t *testing.T) {
	quantity := PickUpMaxSuitableQuantity([]uint{70, 20, 10}, divider.Fair, 5, 10)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMaxSuitableQuantity([]uint{70, 20, 10}, divider.Fair, 7, 10)
	require.Equal(t, uint(6), quantity)

	quantity = PickUpMaxSuitableQuantity([]uint{70, 20, 10}, divider.Fair, 1000, 10)
	require.Equal(t, uint(1000), quantity)

	quantity = PickUpMaxSuitableQuantity([]uint{70, 20, 10}, divider.Rate, 21, 10)
	require.Equal(t, uint(0), quantity)

	quantity = PickUpMaxSuitableQuantity([]uint{70, 20, 10}, divider.Rate, 23, 10)
	require.Equal(t, uint(22), quantity)

	quantity = PickUpMaxSuitableQuantity([]uint{70, 20, 10}, divider.Rate, 1000, 10)
	require.Equal(t, uint(1000), quantity)
}
