package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDistributionFilled(t *testing.T) {
	require.False(t, IsDistributionFilled(map[uint]uint{3: 0, 2: 0, 1: 0}))
	require.False(t, IsDistributionFilled(map[uint]uint{3: 1, 2: 0, 1: 0}))
	require.False(t, IsDistributionFilled(map[uint]uint{3: 0, 2: 1, 1: 0}))
	require.False(t, IsDistributionFilled(map[uint]uint{3: 0, 2: 0, 1: 1}))
	require.False(t, IsDistributionFilled(map[uint]uint{3: 1, 2: 1, 1: 0}))
	require.False(t, IsDistributionFilled(map[uint]uint{3: 1, 2: 0, 1: 1}))
	require.False(t, IsDistributionFilled(map[uint]uint{3: 0, 2: 1, 1: 1}))
	require.True(t, IsDistributionFilled(map[uint]uint{3: 1, 2: 1, 1: 1}))
}
