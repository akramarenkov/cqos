package priority

import (
	"errors"

	"github.com/akramarenkov/cqos/v2/priority/divider"
)

var (
	ErrBadDivider = errors.New("divider produces an incorrect distribution")
)

func calcDistributionQuantity(distribution map[uint]uint) uint {
	quantity := uint(0)

	for _, amount := range distribution {
		quantity += amount
	}

	return quantity
}

func safeDivide(
	divider divider.Divider,
	priorities []uint,
	dividend uint,
	distribution map[uint]uint,
) error {
	before := calcDistributionQuantity(distribution)

	divider(priorities, dividend, distribution)

	after := calcDistributionQuantity(distribution)

	if after == 0 {
		return nil
	}

	if after-before != dividend {
		return ErrBadDivider
	}

	return nil
}
