package join

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalcTickerDuration(t *testing.T) {
	duration, err := calcTickerDuration(100*time.Millisecond, 25)
	require.NoError(t, err)
	require.Equal(t, 25*time.Millisecond, duration)

	duration, err = calcTickerDuration(40*time.Millisecond, 25)
	require.NoError(t, err)
	require.Equal(t, 10*time.Millisecond, duration)

	duration, err = calcTickerDuration(10*time.Millisecond, 25)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), duration)

	duration, err = calcTickerDuration(100*time.Millisecond, 0)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), duration)

	duration, err = calcTickerDuration(100*time.Millisecond, 101)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), duration)
}
