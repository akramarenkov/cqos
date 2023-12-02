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
) (map[uint]uint, error) {
	before := calcDistributionQuantity(distribution)

	updated := divider(priorities, dividend, distribution)

	after := calcDistributionQuantity(updated)

	if after-before != dividend {
		return nil, ErrBadDivider
	}

	return updated, nil
}
