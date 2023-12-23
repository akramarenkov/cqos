package limit

import (
	"errors"
	"time"

	"github.com/akramarenkov/safe"
)

var (
	ErrNegativeMinimumInterval = errors.New("minimum interval is negative")
	ErrZeroConvertedInterval   = errors.New("converted interval is zero")
	ErrZeroConvertedQuantity   = errors.New("converted quantity is zero")
	ErrZeroNegativeInterval    = errors.New("interval is zero or negative")
	ErrZeroQuantity            = errors.New("quantity is zero")
)

const (
	// the value was chosen based on studies of the graphical tests results and benchmarks
	defaultOptimizeMinimumInterval = 10 * time.Millisecond
)

// Quantity per Interval
type Rate struct {
	Interval time.Duration
	Quantity uint64
}

func (rate Rate) IsValid() error {
	if rate.Interval <= 0 {
		return ErrZeroNegativeInterval
	}

	if rate.Quantity == 0 {
		return ErrZeroQuantity
	}

	return nil
}

func (rate Rate) Flatten() (Rate, error) {
	return rate.recalc(0)
}

func (rate Rate) Optimize() (Rate, error) {
	return rate.recalc(defaultOptimizeMinimumInterval)
}

func (rate Rate) recalc(min time.Duration) (Rate, error) {
	if rate.Interval <= 0 {
		return Rate{}, ErrZeroNegativeInterval
	}

	if rate.Quantity == 0 {
		return Rate{}, ErrZeroQuantity
	}

	if min < 0 {
		return Rate{}, ErrNegativeMinimumInterval
	}

	divider, err := safe.UnsignedToSigned[uint64, time.Duration](rate.Quantity)
	if err != nil {
		return Rate{}, err
	}

	interval := rate.Interval / divider
	quantity := uint64(1)

	if interval <= min {
		if min == 0 {
			return Rate{}, ErrZeroConvertedInterval
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
