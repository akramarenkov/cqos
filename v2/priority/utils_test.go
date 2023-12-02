package priority

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalcDistributionQuantity(t *testing.T) {
	require.Equal(t, uint(0), calcDistributionQuantity(nil))
	require.Equal(t, uint(0), calcDistributionQuantity(map[uint]uint{}))
	require.Equal(t, uint(0), calcDistributionQuantity(map[uint]uint{1: 0, 2: 0, 3: 0}))
	require.Equal(t, uint(1), calcDistributionQuantity(map[uint]uint{1: 1, 2: 0, 3: 0}))
	require.Equal(t, uint(1), calcDistributionQuantity(map[uint]uint{1: 0, 2: 1, 3: 0}))
	require.Equal(t, uint(1), calcDistributionQuantity(map[uint]uint{1: 0, 2: 0, 3: 1}))
	require.Equal(t, uint(2), calcDistributionQuantity(map[uint]uint{1: 1, 2: 1, 3: 0}))
	require.Equal(t, uint(2), calcDistributionQuantity(map[uint]uint{1: 0, 2: 1, 3: 1}))
	require.Equal(t, uint(2), calcDistributionQuantity(map[uint]uint{1: 1, 2: 0, 3: 1}))
	require.Equal(t, uint(3), calcDistributionQuantity(map[uint]uint{1: 1, 2: 1, 3: 1}))
	require.Equal(t, uint(3), calcDistributionQuantity(map[uint]uint{1: 1, 2: 1, 3: 1}))
	require.Equal(t, uint(3), calcDistributionQuantity(map[uint]uint{1: 1, 2: 1, 3: 1}))
}

func TestSafeDivide(t *testing.T) {

}
