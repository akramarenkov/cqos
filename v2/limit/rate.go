package limit

import (
	"errors"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/safe"
)

var (
	ErrConvertedIntervalZero   = errors.New("converted interval is zero")
	ErrConvertedQuantityZero   = errors.New("converted quantity is zero")
	ErrIntervalZeroNegative    = errors.New("interval is zero or negative")
	ErrMinimumIntervalNegative = errors.New("minimum interval is negative")
	ErrQuantityZero            = errors.New("quantity is zero")
)

// Quantity data elements per time Interval
type Rate struct {
	Interval time.Duration
	Quantity uint64
}

// Validates field values.
//
// Interval cannot be negative or equal to zero.
//
// Quantity cannot be equal to zero
func (rate Rate) IsValid() error {
	if rate.Interval <= 0 {
		return ErrIntervalZeroNegative
	}

	if rate.Quantity == 0 {
		return ErrQuantityZero
	}

	return nil
}

// Recalculates the units of measurement of the Interval so that the Quantity is
// equal to 1.
//
// Maximizes the uniformity of the distribution of data elements over time by
// reducing the productivity of the discipline
func (rate Rate) Flatten() (Rate, error) {
	return rate.recalc(0)
}

// Recalculates the units of measurement of the interval so that the Quantity is
// as small as possible but the Interval is not less than the recommended value.
//
// Increases the uniformity of the distribution of data elements over time,
// almost without reducing the productivity of the discipline
func (rate Rate) Optimize() (Rate, error) {
	return rate.recalc(consts.ReliablyMeasurableDuration)
}

func (rate Rate) recalc(min time.Duration) (Rate, error) {
	if rate.Interval <= 0 {
		return Rate{}, ErrIntervalZeroNegative
	}

	if rate.Quantity == 0 {
		return Rate{}, ErrQuantityZero
	}

	if min < 0 {
		return Rate{}, ErrMinimumIntervalNegative
	}

	divider, err := safe.UnsignedToSigned[uint64, time.Duration](rate.Quantity)
	if err != nil {
		return Rate{}, err
	}

	interval := rate.Interval / divider
	quantity := uint64(1)

	if interval <= min {
		if min == 0 {
			return Rate{}, ErrConvertedIntervalZero
		}

		interval = min

		product, err := safe.ProductInt(rate.Quantity, uint64(min))
		if err != nil {
			return Rate{}, err
		}

		quantity = product / uint64(rate.Interval)
	}

	rate.Interval = interval
	rate.Quantity = quantity

	return rate, nil
}
