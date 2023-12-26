package priority

import (
	"errors"

	"github.com/akramarenkov/cqos/v2/priority/divider"
	"github.com/akramarenkov/safe"
)

var (
	ErrDividerBad = errors.New("divider produces an incorrect distribution")
)

func calcDistributionQuantity(distribution map[uint]uint) uint {
	quantity := uint(0)

	for _, amount := range distribution {
		quantity += amount
	}

	return quantity
}

func safeCalcDistributionQuantity(distribution map[uint]uint) (uint, error) {
	quantity := uint(0)

	for _, amount := range distribution {
		sum, err := safe.SumInt(quantity, amount)
		if err != nil {
			return 0, err
		}

		quantity = sum
	}

	return quantity, nil
}

func safeDivide(
	divider divider.Divider,
	priorities []uint,
	dividend uint,
	distribution map[uint]uint,
) error {
	before, err := safeCalcDistributionQuantity(distribution)
	if err != nil {
		return err
	}

	divider(priorities, dividend, distribution)

	after, err := safeCalcDistributionQuantity(distribution)
	if err != nil {
		return err
	}

	if after == 0 {
		return nil
	}

	if after-before != dividend {
		return ErrDividerBad
	}

	return nil
}
