package unite

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/join/internal/defaults"

	"github.com/stretchr/testify/require"
)

func TestCalcInterruptInterval(t *testing.T) {
	interval, err := calcInterruptInterval(
		10*defaults.MinTimeout,
		defaults.TimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, 10*consts.ReliablyMeasurableDuration, interval)

	interval, err = calcInterruptInterval(
		defaults.MinTimeout,
		defaults.TimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, consts.ReliablyMeasurableDuration, interval)

	interval, err = calcInterruptInterval(
		consts.ReliablyMeasurableDuration,
		defaults.TimeoutInaccuracy,
	)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(
		defaults.MinTimeout,
		0,
	)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(
		defaults.MinTimeout,
		101,
	)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)
}

func TestCalcInterruptIntervalNonPositiveAllowed(t *testing.T) {
	interval, err := calcInterruptIntervalNonPositiveAllowed(
		0,
		defaults.TimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptIntervalNonPositiveAllowed(
		-defaults.MinTimeout,
		defaults.TimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), interval)
}
