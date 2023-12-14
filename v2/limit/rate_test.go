package limit

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateIsValid(t *testing.T) {
	zeroInterval := Rate{
		Interval: 0,
		Quantity: 10,
	}

	negativeInterval := Rate{
		Interval: -time.Second,
		Quantity: 10,
	}

	zeroQuantity := Rate{
		Interval: time.Second,
		Quantity: 0,
	}

	regular := Rate{
		Interval: time.Second,
		Quantity: 10,
	}

	require.Error(t, zeroInterval.IsValid())
	require.Error(t, negativeInterval.IsValid())
	require.Error(t, zeroQuantity.IsValid())
	require.NoError(t, regular.IsValid())
}

func TestFlatten(t *testing.T) {
	zeroInterval := Rate{
		Interval: 0,
		Quantity: 10,
	}

	negativeInterval := Rate{
		Interval: -time.Second,
		Quantity: 10,
	}

	zeroQuantity := Rate{
		Interval: time.Second,
		Quantity: 0,
	}

	regular := Rate{
		Interval: time.Second,
		Quantity: 10,
	}

	negativeConversion := Rate{
		Interval: time.Second,
		Quantity: math.MaxUint64,
	}

	zeroConversion := Rate{
		Interval: time.Second,
		Quantity: math.MaxInt64,
	}

	flattened, done := zeroInterval.Flatten()
	require.Equal(t, false, done)
	require.Equal(t, zeroInterval, flattened)

	flattened, done = negativeInterval.Flatten()
	require.Equal(t, false, done)
	require.Equal(t, negativeInterval, flattened)

	flattened, done = zeroQuantity.Flatten()
	require.Equal(t, false, done)
	require.Equal(t, zeroQuantity, flattened)

	flattened, done = regular.Flatten()
	require.Equal(t, true, done)
	require.Equal(t, Rate{Interval: 100 * time.Millisecond, Quantity: 1}, flattened)

	flattened, done = negativeConversion.Flatten()
	require.Equal(t, false, done)
	require.Equal(t, negativeConversion, flattened)

	flattened, done = zeroConversion.Flatten()
	require.Equal(t, false, done)
	require.Equal(t, zeroConversion, flattened)
}
