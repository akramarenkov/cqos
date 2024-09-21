package limit

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateIsValid(t *testing.T) {
	rate := Rate{
		Interval: time.Second,
		Quantity: 10,
	}
	require.NoError(t, rate.IsValid())

	rate = Rate{
		Interval: -time.Second,
		Quantity: 10,
	}
	require.Error(t, rate.IsValid())

	rate = Rate{
		Interval: 0,
		Quantity: 10,
	}
	require.Error(t, rate.IsValid())

	rate = Rate{
		Interval: time.Second,
		Quantity: 0,
	}
	require.Error(t, rate.IsValid())
}

func TestRateFlatten(t *testing.T) {
	rate := Rate{
		Interval: time.Second,
		Quantity: 10,
	}
	flatten, err := rate.Flatten()
	require.NoError(t, err)
	require.Equal(
		t,
		Rate{
			Interval: 100 * time.Millisecond,
			Quantity: 1,
		},
		flatten,
	)

	rate = Rate{
		Interval: time.Second,
		Quantity: math.MaxInt64,
	}
	flatten, err = rate.Flatten()
	require.Error(t, err)
	require.Equal(t, Rate{}, flatten)
}

func TestRateOptimize(t *testing.T) {
	rate := Rate{
		Interval: time.Second,
		Quantity: 10,
	}
	optimized, err := rate.Optimize()
	require.NoError(t, err)
	require.Equal(
		t,
		Rate{
			Interval: 100 * time.Millisecond,
			Quantity: 1,
		},
		optimized,
	)

	rate = Rate{
		Interval: time.Second,
		Quantity: 1e4,
	}
	optimized, err = rate.Optimize()
	require.NoError(t, err)
	require.Equal(
		t,
		Rate{
			Interval: OptimizationInterval,
			Quantity: 100,
		},
		optimized,
	)

	rate = Rate{
		Interval: time.Second,
		Quantity: math.MaxUint64,
	}
	optimized, err = rate.Optimize()
	require.NoError(t, err)
	require.Equal(
		t,
		Rate{
			Interval: OptimizationInterval,
			Quantity: 184467440737095516,
		},
		optimized,
	)

	rate = Rate{
		Interval: time.Millisecond,
		Quantity: math.MaxUint64,
	}
	optimized, err = rate.Optimize()
	require.Error(t, err)
	require.Equal(t, Rate{}, optimized)
}

func TestRateRecalculate(t *testing.T) {
	rate := Rate{
		Interval: -time.Second,
		Quantity: 10,
	}
	recalculated, err := rate.Recalculate(time.Millisecond)
	require.Error(t, err)
	require.Equal(t, Rate{}, recalculated)

	rate = Rate{
		Interval: time.Second,
		Quantity: 10,
	}
	recalculated, err = rate.Recalculate(-time.Millisecond)
	require.Error(t, err)
	require.Equal(t, Rate{}, recalculated)
}
