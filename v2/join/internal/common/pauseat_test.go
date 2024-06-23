package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalcPauseAtDuration(t *testing.T) {
	require.Equal(t, 275*time.Millisecond, CalcPauseAtDuration(100*time.Millisecond))
}
