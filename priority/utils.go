package priority

import (
	"math"
	"sort"
)

const (
	oneHundredPercent = 100
	referenceFactor   = 1000
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

func isDistributionFilled(distribution map[uint]uint) bool {
	for _, quantity := range distribution {
		if quantity == 0 {
			return false
		}
	}

	return true
}

func isNonFatalConfig(
	combinations [][]uint,
	divider Divider,
	handlersQuantity uint,
) bool {
	for _, combination := range combinations {
		distribution := divider(combination, handlersQuantity, nil)

		if !isDistributionFilled(distribution) {
			return false
		}
	}

	return true
}

// Due to the imperfection of the division function and working with integers (since
// the number of handlers is an integer), large errors can occur when distributing
// handlers by priority, especially for small numbers of handlers. This function allows
// you to determine if a situation occurs when one or more priorities are not processed
// (for none of the priorities, the quantity is not equal to zero) with the specified
// combination of priorities, division function and number of handlers
func IsNonFatalConfig(
	priorities []uint,
	divider Divider,
	handlersQuantity uint,
) bool {
	combinations := genPriorityCombinations(priorities)

	return isNonFatalConfig(combinations, divider, handlersQuantity)
}

// Picks up the minimum number of handlers for which the division error does not
// cause stop processing of one or more priorities
func PickUpMinNonFatalHandlersQuantity(
	priorities []uint,
	divider Divider,
	maxHandlersQuantity uint,
) uint {
	combinations := genPriorityCombinations(priorities)

	for quantity := uint(1); quantity <= maxHandlersQuantity; quantity++ {
		if isNonFatalConfig(combinations, divider, quantity) {
			return quantity
		}
	}

	return 0
}

// Picks up the maximum number of handlers for which the division error does not
// cause stop processing of one or more priorities
func PickUpMaxNonFatalHandlersQuantity(
	priorities []uint,
	divider Divider,
	maxHandlersQuantity uint,
) uint {
	combinations := genPriorityCombinations(priorities)

	for quantity := maxHandlersQuantity; quantity > 0; quantity-- {
		if isNonFatalConfig(combinations, divider, quantity) {
			return quantity
		}
	}

	return 0
}

func isDistributionSuitable(
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

		ratio := float64(referenceHandlersQuantity) / float64(handlersQuantity)

		diff := 1.0 - ratio*float64(distribution[priority])/float64(referenceQuantity)

		diff = oneHundredPercent * math.Abs(diff)

		if diff > maxDiff {
			return false
		}
	}

	return true
}

func isSuitableConfig(
	combinations [][]uint,
	priorities []uint,
	divider Divider,
	handlersQuantity uint,
	maxDiff float64,
) bool {
	referenceHandlersQuantity := sumPriorities(priorities) * referenceFactor

	for _, combination := range combinations {
		distribution := divider(combination, handlersQuantity, nil)

		if !isDistributionFilled(distribution) {
			return false
		}

		reference := divider(combination, referenceHandlersQuantity, nil)

		suitable := isDistributionSuitable(
			distribution,
			reference,
			handlersQuantity,
			referenceHandlersQuantity,
			maxDiff,
		)

		if !suitable {
			return false
		}
	}

	return true
}

// Due to the imperfection of the division function and working with integers (since
// the number of handlers is an integer), large errors can occur when distributing
// handlers by priority, especially for small numbers of handlers. This function allows
// you to determine that with the specified combination of priorities, the division
// function and the number of handlers, the distribution error does not exceed
// the required value
func IsSuitableConfig(
	priorities []uint,
	divider Divider,
	handlersQuantity uint,
	maxDiff float64,
) bool {
	combinations := genPriorityCombinations(priorities)

	return isSuitableConfig(combinations, priorities, divider, handlersQuantity, maxDiff)
}

// Picks up the minimum number of handlers for which the division error does not
// exceed the specified value
func PickUpMinSuitableHandlersQuantity(
	priorities []uint,
	divider Divider,
	maxHandlersQuantity uint,
	maxDiff float64,
) uint {
	combinations := genPriorityCombinations(priorities)

	for quantity := uint(1); quantity <= maxHandlersQuantity; quantity++ {
		if isSuitableConfig(combinations, priorities, divider, quantity, maxDiff) {
			return quantity
		}
	}

	return 0
}

// Picks up the maximum number of handlers for which the division error does not
// exceed the specified value
func PickUpMaxSuitableHandlersQuantity(
	priorities []uint,
	divider Divider,
	maxHandlersQuantity uint,
	maxDiff float64,
) uint {
	combinations := genPriorityCombinations(priorities)

	for quantity := maxHandlersQuantity; quantity > 0; quantity-- {
		if isSuitableConfig(combinations, priorities, divider, quantity, maxDiff) {
			return quantity
		}
	}

	return 0
}
