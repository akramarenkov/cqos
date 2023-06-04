package priority

import (
	"math"
	"sort"
)

const (
	oneHundredPercent               = 100
	referenceHandlersQuantityFactor = 1000
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

func isDistributionFilled(distribution map[uint]uint) bool {
	for _, quantity := range distribution {
		if quantity == 0 {
			return false
		}
	}

	return true
}

// inefficient implementation, but usually there are not so many priorities for
// this to be a problem
func genPriorityCombinations(priorities []uint) [][]uint {
	combinations := make([][]uint, 0)

	for _, priority := range priorities {
		combinations = append(combinations, []uint{priority})
	}

	traversed := make(map[int][][]uint)

	for _, priority := range priorities {
		for _, combination := range combinations {
			if isPriorityExists(combination, priority) {
				continue
			}

			copied := copyCombination(combination, priority)

			if isCombinationExists(traversed[len(copied)], copied) {
				continue
			}

			traversed[len(copied)] = append(traversed[len(copied)], copied)

			combinations = append(combinations, copied)
		}
	}

	return combinations
}

func copyCombination(priorities []uint, priority uint) []uint {
	copied := make([]uint, len(priorities)+1)

	copy(copied, priorities)

	copied[len(copied)-1] = priority

	sortPriorities(copied)

	return copied
}

func isPriorityExists(priorities []uint, verifiable uint) bool {
	for _, priority := range priorities {
		if priority == verifiable {
			return true
		}
	}

	return false
}

func isSortedPrioritiesEqual(left []uint, right []uint) bool {
	if len(left) != len(right) {
		return false
	}

	for id := range left {
		if right[id] != left[id] {
			return false
		}
	}

	return true
}

func isCombinationExists(combinations [][]uint, verifiable []uint) bool {
	for _, combination := range combinations {
		if isSortedPrioritiesEqual(combination, verifiable) {
			return true
		}
	}

	return false
}

// Checks if the combination of priorities, handlers quantity and divider is
// such that not causes a fatal error (for none of the priorities, the quantity is
// not equal to zero) in distributions handlers quantity by priorities
func IsNonFatalConfig(
	priorities []uint,
	divider Divider,
	handlersQuantity uint,
) bool {
	combinations := genPriorityCombinations(priorities)

	for _, combination := range combinations {
		distribution := divider(combination, handlersQuantity, nil)

		if !isDistributionFilled(distribution) {
			return false
		}
	}

	return true
}

func isDistributionHasSmallError(
	distribution map[uint]uint,
	reference map[uint]uint,
	handlersQuantity uint,
	referenceHandlersQuantity uint,
	maxDiff float64,
) bool {
	for priority, referenceQuantity := range reference {
		if referenceQuantity == 0 {
			return false
		}

		handlersQuantityRatio := float64(referenceHandlersQuantity) / float64(handlersQuantity)

		diff := 1.0 - handlersQuantityRatio*float64(distribution[priority])/float64(referenceQuantity)
		diff = oneHundredPercent * math.Abs(diff)

		if diff > maxDiff {
			return false
		}
	}

	return true
}

// Checks if the combination of priorities, handlers quantity and divider is
// such that causes a small error (in percent) in distributions handlers quantity by priorities.
func IsSmallErrorConfig(
	priorities []uint,
	divider Divider,
	handlersQuantity uint,
	maxDiff float64,
) bool {
	referenceHandlersQuantity := sumPriorities(priorities) * referenceHandlersQuantityFactor

	combinations := genPriorityCombinations(priorities)

	for _, combination := range combinations {
		distribution := divider(combination, handlersQuantity, nil)
		reference := divider(combination, referenceHandlersQuantity, nil)

		if !isDistributionFilled(distribution) {
			return false
		}

		small := isDistributionHasSmallError(
			distribution,
			reference,
			handlersQuantity,
			referenceHandlersQuantity,
			maxDiff,
		)

		if !small {
			return false
		}
	}

	return true
}
