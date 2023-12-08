package common

import "math"

func CalcByFactor(base int, factor float64, min int) int {
	capacity := int(math.Round(factor * float64(base)))

	if capacity < min {
		return min
	}

	return capacity
}
