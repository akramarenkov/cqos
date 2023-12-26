package join

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalcInterruptInterval(t *testing.T) {
	interval, err := calcInterruptInterval(100*time.Millisecond, 25)
	require.NoError(t, err)
	require.Equal(t, 25*time.Millisecond, interval)

	interval, err = calcInterruptInterval(40*time.Millisecond, 25)
	require.NoError(t, err)
	require.Equal(t, 10*time.Millisecond, interval)

	interval, err = calcInterruptInterval(10*time.Millisecond, 25)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(100*time.Millisecond, 0)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(100*time.Millisecond, 101)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)
}

func TestCalcInterruptIntervalZeroAllowed(t *testing.T) {
	interval, err := calcInterruptIntervalZeroAllowed(100*time.Millisecond, 25)
	require.NoError(t, err)
	require.Equal(t, 25*time.Millisecond, interval)

	interval, err = calcInterruptIntervalZeroAllowed(40*time.Millisecond, 25)
	require.NoError(t, err)
	require.Equal(t, 10*time.Millisecond, interval)

	interval, err = calcInterruptIntervalZeroAllowed(0, 25)
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptIntervalZeroAllowed(-100*time.Millisecond, 25)
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptIntervalZeroAllowed(10*time.Millisecond, 25)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptIntervalZeroAllowed(100*time.Millisecond, 0)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptIntervalZeroAllowed(100*time.Millisecond, 101)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)
}
