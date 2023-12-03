package priority

import (
	"math"
	"testing"

	"github.com/akramarenkov/cqos/v2/priority/divider"
	"github.com/stretchr/testify/require"
)

func TestCalcDistributionQuantity(t *testing.T) {
	quantity, err := calcDistributionQuantity(nil)
	require.Equal(t, uint(0), quantity)
	require.NoError(t, err)

	quantity, err = calcDistributionQuantity(map[uint]uint{})
	require.Equal(t, uint(0), quantity)
	require.NoError(t, err)

	quantity, err = calcDistributionQuantity(map[uint]uint{1: 0, 2: 0, 3: 0})
	require.Equal(t, uint(0), quantity)
	require.NoError(t, err)

	quantity, err = calcDistributionQuantity(map[uint]uint{1: 1, 2: 0, 3: 0})
	require.Equal(t, uint(1), quantity)
	require.NoError(t, err)

	quantity, err = calcDistributionQuantity(map[uint]uint{1: 0, 2: 1, 3: 0})
	require.Equal(t, uint(1), quantity)
	require.NoError(t, err)

	quantity, err = calcDistributionQuantity(map[uint]uint{1: 0, 2: 0, 3: 1})
	require.Equal(t, uint(1), quantity)
	require.NoError(t, err)

	quantity, err = calcDistributionQuantity(map[uint]uint{1: 1, 2: 1, 3: 0})
	require.Equal(t, uint(2), quantity)
	require.NoError(t, err)

	quantity, err = calcDistributionQuantity(map[uint]uint{1: 0, 2: 1, 3: 1})
	require.Equal(t, uint(2), quantity)
	require.NoError(t, err)

	quantity, err = calcDistributionQuantity(map[uint]uint{1: 1, 2: 0, 3: 1})
	require.Equal(t, uint(2), quantity)
	require.NoError(t, err)

	quantity, err = calcDistributionQuantity(map[uint]uint{1: 1, 2: 1, 3: 1})
	require.Equal(t, uint(3), quantity)
	require.NoError(t, err)

	quantity, err = calcDistributionQuantity(map[uint]uint{1: 1, 2: 1, 3: 1})
	require.Equal(t, uint(3), quantity)
	require.NoError(t, err)

	quantity, err = calcDistributionQuantity(map[uint]uint{1: 1, 2: 1, 3: 1})
	require.Equal(t, uint(3), quantity)
	require.NoError(t, err)
}

func TestCalcDistributionQuantityError(t *testing.T) {
	distribution := map[uint]uint{1: math.MaxUint - 2, 2: 2, 3: 1}

	quantity, err := calcDistributionQuantity(distribution)
	require.Equal(t, uint(0), quantity)
	require.Error(t, err)
}

func TestSafeDivide(t *testing.T) {
	badDivider := func(
		priorities []uint,
		dividend uint,
		distribution map[uint]uint,
	) {
		divider.Fair(priorities, dividend, distribution)

		for priority := range distribution {
			distribution[priority] *= 2
		}
	}

	distribution := make(map[uint]uint)
	err := safeDivide(divider.Fair, []uint{3, 2, 1}, 6, distribution)
	require.NoError(t, err)

	distribution = map[uint]uint{3: 1, 2: 2, 1: 0}
	err = safeDivide(divider.Fair, []uint{3, 2, 1}, 6, distribution)
	require.NoError(t, err)

	distribution = make(map[uint]uint)
	err = safeDivide(divider.Fair, nil, 6, distribution)
	require.NoError(t, err)

	distribution = make(map[uint]uint)
	err = safeDivide(divider.Fair, []uint{3, 2, 1}, 0, distribution)
	require.NoError(t, err)

	distribution = make(map[uint]uint)
	err = safeDivide(badDivider, []uint{3, 2, 1}, 6, distribution)
	require.Error(t, err)

	distribution = map[uint]uint{3: 1, 2: 2, 1: 0}
	err = safeDivide(badDivider, []uint{3, 2, 1}, 6, distribution)
	require.Error(t, err)

	distribution = map[uint]uint{3: math.MaxUint - 2, 2: 2, 1: 0}
	err = safeDivide(badDivider, []uint{3, 2, 1}, 1, distribution)
	require.Error(t, err)

	distribution = map[uint]uint{3: math.MaxUint - 2, 2: 2, 1: 1}
	err = safeDivide(badDivider, []uint{3, 2, 1}, 1, distribution)
	require.Error(t, err)
}
