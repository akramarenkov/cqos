package priority

import (
	"math"
)

func removePriority(priorities []uint, removed uint) []uint {
	kept := 0

	for _, priority := range priorities {
		if priority == removed {
			continue
		}

		priorities[kept] = priority
		kept++
	}

	return priorities[:kept]
}

func calcCapacity(base int, factor float64, reserve int) int {
	capacity := int(math.Round(factor * float64(base)))

	if capacity == 0 {
		capacity = reserve
	}

	return capacity
}
