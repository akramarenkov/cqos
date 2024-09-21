package limit

import (
	"errors"
	"math/big"
	"time"
)

var (
	ErrConvertedIntervalZero            = errors.New("converted interval is zero, quantity is too large")
	ErrConvertedQuantityUnrepresentable = errors.New("converted quantity is unrepresentable by used type")
	ErrIntervalNegative                 = errors.New("interval is negative")
	ErrIntervalZero                     = errors.New("interval is zero")
	ErrMinimumIntervalNegative          = errors.New("minimum interval is negative")
	ErrQuantityZero                     = errors.New("quantity is zero")
)

var (
	// Deprecated.
	ErrConvertedQuantityZero = errors.New("converted quantity is zero")
)

// Quantity of data elements passed per time Interval.
type Rate struct {
	Interval time.Duration
	Quantity uint64
}

// Validates field values. Interval cannot be negative or equal to zero. Quantity
// cannot be equal to zero.
func (rt Rate) IsValid() error {
	if rt.Interval < 0 {
		return ErrIntervalNegative
	}

	if rt.Interval == 0 {
		return ErrIntervalZero
	}

	if rt.Quantity == 0 {
		return ErrQuantityZero
	}

	return nil
}

// Recalculates the units of measurement of the Interval so that the Quantity is
// equal to 1.
//
// Maximizes the uniformity of the distribution of output data elements over time by
// reducing the productivity of the discipline.
func (rt Rate) Flatten() (Rate, error) {
	return rt.Recalculate(0)
}

// Recalculates the units of measurement of the Interval so that the Quantity is
// as small as possible but the Interval is not less than the recommended value.
//
// Increases the uniformity of the distribution of output data elements over time,
// almost without reducing the productivity of the discipline.
func (rt Rate) Optimize() (Rate, error) {
	return rt.Recalculate(OptimizationInterval)
}

// Recalculates the units of measurement of an Interval with a limitation on its
// minimum value.
func (rt Rate) Recalculate(minimum time.Duration) (Rate, error) {
	if err := rt.IsValid(); err != nil {
		return Rate{}, err
	}

	if minimum < 0 {
		return Rate{}, ErrMinimumIntervalNegative
	}

	// integer overflows are not possible given the checks above
	interval := time.Duration(uint64(rt.Interval) / rt.Quantity)

	if interval > minimum {
		recalculated := Rate{
			Interval: interval,
			Quantity: 1,
		}

		return recalculated, nil
	}

	if minimum == 0 {
		return Rate{}, ErrConvertedIntervalZero
	}

	quantity, err := recalculateQuantity(rt.Quantity, minimum, rt.Interval)
	if err != nil {
		return Rate{}, err
	}

	recalculated := Rate{
		Interval: minimum,
		Quantity: quantity,
	}

	return recalculated, nil
}

func recalculateQuantity(
	quantity uint64,
	minimum time.Duration,
	interval time.Duration,
) (uint64, error) {
	qb := new(big.Int).SetUint64(quantity)
	mb := new(big.Int).SetInt64(int64(minimum))
	ib := new(big.Int).SetInt64(int64(interval))

	product := new(big.Int).Mul(qb, mb)
	quotient := new(big.Int).Quo(product, ib)

	if !quotient.IsUint64() {
		return 0, ErrConvertedQuantityUnrepresentable
	}

	return quotient.Uint64(), nil
}
