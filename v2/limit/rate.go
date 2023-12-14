package limit

import (
	"errors"
	"time"
)

var (
	ErrInvalidInterval = errors.New("invalid interval")
	ErrInvalidQuantity = errors.New("invalid quantity")
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

func (rate Rate) Flatten() (Rate, bool) {
	if rate.Interval <= 0 {
		return rate, false
	}

	if rate.Quantity == 0 {
		return rate, false
	}

	interval := rate.Interval / time.Duration(rate.Quantity)

	if interval <= 0 {
		return rate, false
	}

	rate.Interval = interval
	rate.Quantity = 1

	return rate, true
}
