package priority

import (
	"math"
	"sort"
)

func sortPriorities(priorities []uint) {
	less := func(i int, j int) bool {
		return priorities[j] < priorities[i]
	}

	sort.SliceStable(priorities, less)
}

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

func sumPriorities(priorities []uint) uint {
	sum := uint(0)

	for _, priority := range priorities {
		sum += priority
	}

	return sum
}

func calcCapacity(base int, factor float64, reserve int) int {
	capacity := int(math.Round(factor * float64(base)))

	if capacity == 0 {
		capacity = reserve
	}

	return capacity
}
