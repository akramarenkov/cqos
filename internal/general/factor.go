package general

import (
	"math"

	"github.com/akramarenkov/safe"

	"golang.org/x/exp/constraints"
)

func CalcByFactor[Type constraints.Integer](
	base Type,
	factor float64,
	min Type,
) (Type, error) {
	product := math.Round(factor * float64(base))

	converted, err := safe.FloatToInt[float64, Type](product)
	if err != nil {
		return 0, err
	}

	if converted < min {
		return min, nil
	}

	return converted, nil
}
