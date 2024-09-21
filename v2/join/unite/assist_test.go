package unite

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/join/internal/defaults"

	"github.com/stretchr/testify/require"
)

func TestCalcInterruptInterval(t *testing.T) {
	interval, err := calcInterruptInterval(-1, defaults.TimeoutInaccuracy)
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(0, defaults.TimeoutInaccuracy)
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(time.Second, defaults.TimeoutInaccuracy)
	require.NoError(t, err)
	require.Equal(t, 250*time.Millisecond, interval)
}

func TestCalcInterruptIntervalError(t *testing.T) {
	interval, err := calcInterruptInterval(time.Second, 0)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(time.Second, consts.HundredPercent+1)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(3*time.Nanosecond, defaults.TimeoutInaccuracy)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)
}
