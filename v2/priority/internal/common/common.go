// Internal package with functions and data types used in other packages
package common

import (
	"math"
	"sort"
)

func SortPriorities(priorities []uint) {
	less := func(i int, j int) bool {
		return priorities[j] < priorities[i]
	}

	sort.SliceStable(priorities, less)
}

func SumPriorities(priorities []uint) uint {
	sum := uint(0)

	for _, priority := range priorities {
		sum += priority
	}

	return sum
}

func CalcByFactor(base int, factor float64, min int) int {
	capacity := int(math.Round(factor * float64(base)))

	if capacity < min {
		return min
	}

	return capacity
}
