package general

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalcByFactor(t *testing.T) {
	calced, err := CalcByFactor(10, 0.1, 3)
	require.NoError(t, err)
	require.Equal(t, 3, calced)

	calced, err = CalcByFactor(10, 0.1, 0)
	require.NoError(t, err)
	require.Equal(t, 1, calced)

	calced, err = CalcByFactor(100, 0.1, 3)
	require.NoError(t, err)
	require.Equal(t, 10, calced)

	calced, err = CalcByFactor(14, 0.1, 3)
	require.NoError(t, err)
	require.Equal(t, 3, calced)

	calced, err = CalcByFactor(15, 0.1, 3)
	require.NoError(t, err)
	require.Equal(t, 3, calced)

	calced, err = CalcByFactor(16, 0.1, 3)
	require.NoError(t, err)
	require.Equal(t, 3, calced)

	calced, err = CalcByFactor(24, 0.1, 3)
	require.NoError(t, err)
	require.Equal(t, 3, calced)

	calced, err = CalcByFactor(25, 0.1, 3)
	require.NoError(t, err)
	require.Equal(t, 3, calced)

	calced, err = CalcByFactor(26, 0.1, 3)
	require.NoError(t, err)
	require.Equal(t, 3, calced)

	calced, err = CalcByFactor(34, 0.1, 3)
	require.NoError(t, err)
	require.Equal(t, 3, calced)

	calced, err = CalcByFactor(35, 0.1, 3)
	require.NoError(t, err)
	require.Equal(t, 4, calced)

	calced, err = CalcByFactor(36, 0.1, 3)
	require.NoError(t, err)
	require.Equal(t, 4, calced)

	calced, err = CalcByFactor(4, 0.1, 0)
	require.NoError(t, err)
	require.Equal(t, 0, calced)

	calced, err = CalcByFactor(5, 0.1, 0)
	require.NoError(t, err)
	require.Equal(t, 1, calced)

	calced, err = CalcByFactor(6, 0.1, 0)
	require.NoError(t, err)
	require.Equal(t, 1, calced)
}

func TestCalcByFactorOverflow(t *testing.T) {
	calced, err := CalcByFactor(1, math.MaxFloat64, 0)
	require.Error(t, err)
	require.Equal(t, 0, calced)

	calced, err = CalcByFactor(6, math.MaxFloat64, 0)
	require.Error(t, err)
	require.Equal(t, 0, calced)

	calced, err = CalcByFactor(6, math.MaxFloat64/10, 0)
	require.Error(t, err)
	require.Equal(t, 0, calced)
}
