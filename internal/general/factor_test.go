package general

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalcByFactor(t *testing.T) {
	require.Equal(t, 3, CalcByFactor(10, 0.1, 3))
	require.Equal(t, 1, CalcByFactor(10, 0.1, 0))
	require.Equal(t, 10, CalcByFactor(100, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(14, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(15, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(16, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(24, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(25, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(26, 0.1, 3))
	require.Equal(t, 3, CalcByFactor(34, 0.1, 3))
	require.Equal(t, 4, CalcByFactor(35, 0.1, 3))
	require.Equal(t, 4, CalcByFactor(36, 0.1, 3))
	require.Equal(t, 0, CalcByFactor(4, 0.1, 0))
	require.Equal(t, 1, CalcByFactor(5, 0.1, 0))
	require.Equal(t, 1, CalcByFactor(6, 0.1, 0))
}
