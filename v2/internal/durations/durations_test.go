package durations

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalcTotalDuration(t *testing.T) {
	require.Equal(t, time.Duration(0), CalcTotalDuration(nil))
	require.Equal(t, time.Duration(0), CalcTotalDuration([]time.Duration{}))
	require.Equal(t, time.Duration(5), CalcTotalDuration([]time.Duration{5, 4, 3, 2, 1, 0}))
	require.Equal(t, time.Duration(5), CalcTotalDuration([]time.Duration{5, 4, 3, 2, 1}))
	require.Equal(t, time.Duration(5), CalcTotalDuration([]time.Duration{5, 4, 3, 2}))
}
