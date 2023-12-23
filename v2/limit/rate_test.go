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
	regular := Rate{
		Interval: time.Second,
		Quantity: 10,
	}

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

	negativeConvertedQuantity := Rate{
		Interval: time.Second,
		Quantity: math.MaxUint64,
	}

	zeroConvertedInterval := Rate{
		Interval: time.Second,
		Quantity: math.MaxInt64,
	}

	flatten, err := regular.Flatten()
	require.NoError(t, err)
	require.Equal(t, Rate{Interval: 100 * time.Millisecond, Quantity: 1}, flatten)

	_, err = zeroInterval.Flatten()
	require.Error(t, err)

	_, err = negativeInterval.Flatten()
	require.Error(t, err)

	_, err = zeroQuantity.Flatten()
	require.Error(t, err)

	_, err = negativeConvertedQuantity.Flatten()
	require.Error(t, err)

	_, err = zeroConvertedInterval.Flatten()
	require.Error(t, err)
}

func TestRecalc(t *testing.T) {
	regular := Rate{
		Interval: time.Second,
		Quantity: 10,
	}

	regularLessMinimum := Rate{
		Interval: time.Second,
		Quantity: 1e4,
	}

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

	negativeConvertedQuantity := Rate{
		Interval: time.Second,
		Quantity: math.MaxUint64,
	}

	zeroConvertedInterval := Rate{
		Interval: time.Second,
		Quantity: math.MaxInt64,
	}

	recalced, err := regular.recalc(time.Millisecond)
	require.NoError(t, err)
	require.Equal(t, Rate{Interval: 100 * time.Millisecond, Quantity: 1}, recalced)

	recalced, err = regularLessMinimum.recalc(time.Millisecond)
	require.NoError(t, err)
	require.Equal(t, Rate{Interval: 1 * time.Millisecond, Quantity: 10}, recalced)

	_, err = regular.recalc(-time.Second)
	require.Error(t, err)

	_, err = zeroInterval.recalc(0)
	require.Error(t, err)

	_, err = negativeInterval.recalc(0)
	require.Error(t, err)

	_, err = zeroQuantity.recalc(0)
	require.Error(t, err)

	_, err = negativeConvertedQuantity.recalc(0)
	require.Error(t, err)

	_, err = zeroConvertedInterval.recalc(0)
	require.Error(t, err)

	_, err = zeroConvertedInterval.recalc(time.Millisecond)
	require.Error(t, err)
}
