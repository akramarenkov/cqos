package limit

import (
	"errors"
	"time"

	"github.com/akramarenkov/safe"
)

var (
	ErrConvertedIntervalIsZero = errors.New("converted interval is zero")
	ErrConvertedQuantityIsZero = errors.New("converted quantity is zero")
	ErrInvalidInterval         = errors.New("invalid interval")
	ErrInvalidMinimumInterval  = errors.New("invalid minimum interval")
	ErrInvalidQuantity         = errors.New("invalid quantity")
)

const (
	// the value was chosen based on studies of the results of graphical tests
	defaultMinimumInterval = 1 * time.Millisecond
)

// Quantity per Interval
type Rate struct {
	Interval time.Duration
	Quantity uint64
}

func (rate Rate) IsValid() error {
	if rate.Interval <= 0 {
		return ErrInvalidInterval
	}

	if rate.Quantity == 0 {
		return ErrInvalidQuantity
	}

	return nil
}

func (rate Rate) Flatten() (Rate, error) {
	return rate.recalc(0)
}

func (rate Rate) Optimize() (Rate, error) {
	return rate.recalc(defaultMinimumInterval)
}

func (rate Rate) recalc(min time.Duration) (Rate, error) {
	if rate.Interval <= 0 {
		return Rate{}, ErrInvalidInterval
	}

	if rate.Quantity == 0 {
		return Rate{}, ErrInvalidQuantity
	}

	if min < 0 {
		return Rate{}, ErrInvalidMinimumInterval
	}

	divider, err := safe.UnsignedToSigned[uint64, time.Duration](rate.Quantity)
	if err != nil {
		return Rate{}, err
	}

	interval := rate.Interval / divider
	quantity := uint64(1)

	if interval <= min {
		if min == 0 {
			return Rate{}, ErrConvertedIntervalIsZero
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
