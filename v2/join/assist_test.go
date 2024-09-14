package join

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/join/internal/common"

	"github.com/stretchr/testify/require"
)

func TestCalcInterruptInterval(t *testing.T) {
	interval, err := calcInterruptInterval(
		10*common.DefaultMinTimeout,
		common.DefaultTimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, 10*consts.ReliablyMeasurableDuration, interval)

	interval, err = calcInterruptInterval(
		common.DefaultMinTimeout,
		common.DefaultTimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, consts.ReliablyMeasurableDuration, interval)

	interval, err = calcInterruptInterval(
		consts.ReliablyMeasurableDuration,
		common.DefaultTimeoutInaccuracy,
	)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(
		common.DefaultMinTimeout,
		0,
	)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptInterval(
		common.DefaultMinTimeout,
		101,
	)
	require.Error(t, err)
	require.Equal(t, time.Duration(0), interval)
}

func TestCalcInterruptIntervalNonPositiveAllowed(t *testing.T) {
	interval, err := calcInterruptIntervalNonPositiveAllowed(
		0,
		common.DefaultTimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), interval)

	interval, err = calcInterruptIntervalNonPositiveAllowed(
		-common.DefaultMinTimeout,
		common.DefaultTimeoutInaccuracy,
	)
	require.NoError(t, err)
	require.Equal(t, time.Duration(0), interval)
}
