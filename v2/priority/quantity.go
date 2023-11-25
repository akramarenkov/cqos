package priority

import "math"

const (
	oneHundredPercent = 100
	referenceFactor   = 1000
)

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

			newed := newCombination(combination, priority)

			if isCombinationExists(traversed[len(newed)], newed) {
				continue
			}

			traversed[len(newed)] = append(traversed[len(newed)], newed)

			combinations = append(combinations, newed)
		}
	}

	return combinations
}

func isPriorityExists(priorities []uint, verifiable uint) bool {
	for _, priority := range priorities {
		if priority == verifiable {
			return true
		}
	}

	return false
}

func newCombination(combination []uint, added uint) []uint {
	newed := make([]uint, len(combination)+1)

	copy(newed, combination)

	newed[len(newed)-1] = added

	sortPriorities(newed)

	return newed
}

func isCombinationExists(combinations [][]uint, verifiable []uint) bool {
	for _, combination := range combinations {
		if isSortedPrioritiesEqual(combination, verifiable) {
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
	quantity uint,
) bool {
	for _, combination := range combinations {
		distribution := divider(combination, quantity, nil)

		if !isDistributionFilled(distribution) {
			return false
		}
	}

	return true
}

// Due to the imperfection of the dividing function and working with integers (since
// the quantity of handlers is an integer), large errors can occur when distributing
// handlers by priority, especially for small quantity of handlers. This function allows
// you to determine that with the specified combination of priorities, the dividing
// function and the quantity of handlers, the distribution error does not cause stop
// processing of one or more priorities (for none of the priorities, the quantity is
// not equal to zero)
func IsNonFatalConfig(
	priorities []uint,
	divider Divider,
	quantity uint,
) bool {
	combinations := genPriorityCombinations(priorities)

	return isNonFatalConfig(combinations, divider, quantity)
}

// Picks up the minimum quantity of handlers for which the division error does not
// cause stop processing of one or more priorities
func PickUpMinNonFatalQuantity(
	priorities []uint,
	divider Divider,
	maxQuantity uint,
) uint {
	combinations := genPriorityCombinations(priorities)

	for quantity := uint(1); quantity <= maxQuantity; quantity++ {
		if isNonFatalConfig(combinations, divider, quantity) {
			return quantity
		}
	}

	return 0
}

// Picks up the maximum quantity of handlers for which the division error does not
// cause stop processing of one or more priorities
func PickUpMaxNonFatalQuantity(
	priorities []uint,
	divider Divider,
	maxQuantity uint,
) uint {
	combinations := genPriorityCombinations(priorities)

	for quantity := maxQuantity; quantity > 0; quantity-- {
		if isNonFatalConfig(combinations, divider, quantity) {
			return quantity
		}
	}

	return 0
}

func isDistributionSuitable(
	distribution map[uint]uint,
	reference map[uint]uint,
	totalQuantity uint,
	referenceTotalQuantity uint,
	limit float64,
) bool {
	ratio := float64(referenceTotalQuantity) / float64(totalQuantity)

	for priority, referenceQuantity := range reference {
		// a bug is assumed in the dividing function, in which it always returns 0,
		// even with large quantities
		// or strange combinations of priorities and dividing function are used
		if referenceQuantity == 0 {
			return false
		}

		diff := 1.0 - ratio*float64(distribution[priority])/float64(referenceQuantity)

		diff = oneHundredPercent * math.Abs(diff)

		if diff > limit {
			return false
		}
	}

	return true
}

func isSuitableConfig(
	combinations [][]uint,
	priorities []uint,
	divider Divider,
	quantity uint,
	limit float64,
) bool {
	referenceTotalQuantity := referenceFactor * sumPriorities(priorities)

	for _, combination := range combinations {
		distribution := divider(combination, quantity, nil)

		if !isDistributionFilled(distribution) {
			return false
		}

		reference := divider(combination, referenceTotalQuantity, nil)

		suitable := isDistributionSuitable(
			distribution,
			reference,
			quantity,
			referenceTotalQuantity,
			limit,
		)

		if !suitable {
			return false
		}
	}

	return true
}

// Due to the imperfection of the dividing function and working with integers (since
// the quantity of handlers is an integer), large errors can occur when distributing
// handlers by priority, especially for small quantity of handlers. This function allows
// you to determine that with the specified combination of priorities, the dividing
// function and the quantity of handlers, the distribution error does not exceed
// the limit
func IsSuitableConfig(
	priorities []uint,
	divider Divider,
	quantity uint,
	limit float64,
) bool {
	combinations := genPriorityCombinations(priorities)

	return isSuitableConfig(combinations, priorities, divider, quantity, limit)
}

// Picks up the minimum quantity of handlers for which the division error does not
// exceed the limit
func PickUpMinSuitableQuantity(
	priorities []uint,
	divider Divider,
	maxQuantity uint,
	limit float64,
) uint {
	combinations := genPriorityCombinations(priorities)

	for quantity := uint(1); quantity <= maxQuantity; quantity++ {
		if isSuitableConfig(combinations, priorities, divider, quantity, limit) {
			return quantity
		}
	}

	return 0
}

// Picks up the maximum quantity of handlers for which the division error does not
// exceed the limit
func PickUpMaxSuitableQuantity(
	priorities []uint,
	divider Divider,
	maxQuantity uint,
	limit float64,
) uint {
	combinations := genPriorityCombinations(priorities)

	for quantity := maxQuantity; quantity > 0; quantity-- {
		if isSuitableConfig(combinations, priorities, divider, quantity, limit) {
			return quantity
		}
	}

	return 0
}
