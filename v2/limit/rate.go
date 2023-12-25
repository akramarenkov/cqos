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

// Quantity per Interval
type Rate struct {
	Interval time.Duration
	Quantity uint64
}

func (rate Rate) IsValid() error {
	if rate.Interval <= 0 {
		return ErrIntervalZeroNegative
	}

	if rate.Quantity == 0 {
		return ErrQuantityZero
	}

	return nil
}

func (rate Rate) Flatten() (Rate, error) {
	return rate.recalc(0)
}

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
