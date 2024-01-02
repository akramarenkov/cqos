package join

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/stretchr/testify/require"
)

func TestCalcInterruptInterval(t *testing.T) {
	interval, err := calcInterruptInterval(
		10*minDefaultTimeout,
		defaultTimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, 10*consts.ReliablyMeasurableDuration, interval)

	interval, err = calcInterruptInterval(
		minDefaultTimeout,
		defaultTimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, consts.ReliablyMeasurableDuration, interval)

	interval, err = calcInterruptInterval(
		consts.ReliablyMeasurableDuration,
		defaultTimeoutInaccuracy,
	)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(
		minDefaultTimeout,
		0,
	)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(
		minDefaultTimeout,
		101,
	)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)
}

func TestCalcInterruptIntervalNonPositiveAllowed(t *testing.T) {
	interval, err := calcInterruptIntervalNonPositiveAllowed(
		0,
		defaultTimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptIntervalNonPositiveAllowed(
		-minDefaultTimeout,
		defaultTimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), interval)
}
