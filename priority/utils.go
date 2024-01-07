package priority

import (
	"math"

	"github.com/akramarenkov/cqos/internal/consts"
	"github.com/akramarenkov/cqos/priority/internal/common"
)

const (
	referenceFactor = 1000
)

func createPriorities(quantity int) []uint {
	priorities := make([]uint, 0, quantity)

	for id := quantity; id != 0; id-- {
		priorities = append(priorities, uint(id))
	}

	return priorities
}

// If the number of combinations for n priorities is m, then
// for n+1 priorities the number of combinations is 2m+1
// accordingly, the increment is m+1.
func calcCombinationsQuantitySlow(priorities []uint) int {
	quantity := 0

	for length := 0; length < len(priorities); length++ {
		quantity += quantity + 1
	}

	return quantity
}

// It is easy to see that this corresponds to the function 2^n - 1.
func calcCombinationsQuantity(priorities []uint) int {
	const base = 2

	return int(math.Pow(base, float64(len(priorities)))) - 1
}

// An inefficient implementation, but simple and usually there are not so many
// priorities that this would be a problem.
//
// Slice of priorities must be sorted similar to how it does common.SortPriorities()
// if it is necessary that the priorities also get into the divider being sorted.
func genPriorityCombinations(priorities []uint) [][]uint {
	combinations := make([][]uint, 0, calcCombinationsQuantity(priorities))

	for _, priority := range priorities {
		for _, combination := range combinations {
			combinations = append(combinations, addToCombination(combination, priority))
		}

		combinations = append(combinations, addToCombination(nil, priority))
	}

	return combinations
}

func addToCombination(combination []uint, priority uint) []uint {
	created := make([]uint, len(combination)+1)

	copy(created, combination)

	created[len(created)-1] = priority

	return created
}

func createSortedCopy(priorities []uint) []uint {
	copied := make([]uint, len(priorities))

	copy(copied, priorities)

	common.SortPriorities(copied)

	return copied
}

func isNonFatalConfig(
	combinations [][]uint,
	divider Divider,
	quantity uint,
) bool {
	for _, combination := range combinations {
		distribution := divider(combination, quantity, nil)

		if !common.IsDistributionFilled(distribution) {
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
// not equal to zero).
func IsNonFatalConfig(
	priorities []uint,
	divider Divider,
	quantity uint,
) bool {
	priorities = createSortedCopy(priorities)

	combinations := genPriorityCombinations(priorities)

	return isNonFatalConfig(combinations, divider, quantity)
}

// Picks up the minimum quantity of handlers for which the division error does not
// cause stop processing of one or more priorities.
func PickUpMinNonFatalQuantity(
	priorities []uint,
	divider Divider,
	maxQuantity uint,
) uint {
	priorities = createSortedCopy(priorities)

	combinations := genPriorityCombinations(priorities)

	for quantity := uint(1); quantity <= maxQuantity; quantity++ {
		if isNonFatalConfig(combinations, divider, quantity) {
			return quantity
		}
	}

	return 0
}

// Picks up the maximum quantity of handlers for which the division error does not
// cause stop processing of one or more priorities.
func PickUpMaxNonFatalQuantity(
	priorities []uint,
	divider Divider,
	maxQuantity uint,
) uint {
	priorities = createSortedCopy(priorities)

	combinations := genPriorityCombinations(priorities)

	for quantity := maxQuantity; quantity != 0; quantity-- {
		if isNonFatalConfig(combinations, divider, quantity) {
			return quantity
		}
	}

	return 0
}

// diffLimit is specified as a percentage.
func isDistributionSuitable(
	distribution map[uint]uint,
	reference map[uint]uint,
	totalQuantity uint,
	referenceTotalQuantity uint,
	diffLimit float64,
) bool {
	ratio := float64(referenceTotalQuantity) / float64(totalQuantity)

	for priority, referenceQuantity := range reference {
		// a bug is assumed in the dividing function, in which it always returns 0,
		// even with large quantities
		// or strange combinations of priorities and dividing function are used
		if referenceQuantity == 0 {
			return false
		}

		diff := 1.0 - (ratio*float64(distribution[priority]))/float64(referenceQuantity)

		diff = consts.OneHundredPercent * math.Abs(diff)

		if diff > diffLimit {
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
	diffLimit float64,
) bool {
	referenceTotalQuantity := referenceFactor * common.SumPriorities(priorities)

	for _, combination := range combinations {
		distribution := divider(combination, quantity, nil)

		if !common.IsDistributionFilled(distribution) {
			return false
		}

		reference := divider(combination, referenceTotalQuantity, nil)

		suitable := isDistributionSuitable(
			distribution,
			reference,
			quantity,
			referenceTotalQuantity,
			diffLimit,
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
// the limit, specified as a percentage.
func IsSuitableConfig(
	priorities []uint,
	divider Divider,
	quantity uint,
	limit float64,
) bool {
	priorities = createSortedCopy(priorities)

	combinations := genPriorityCombinations(priorities)

	return isSuitableConfig(combinations, priorities, divider, quantity, limit)
}

// Picks up the minimum quantity of handlers for which the division error does not
// exceed the limit, specified as a percentage.
func PickUpMinSuitableQuantity(
	priorities []uint,
	divider Divider,
	maxQuantity uint,
	limit float64,
) uint {
	priorities = createSortedCopy(priorities)

	combinations := genPriorityCombinations(priorities)

	for quantity := uint(1); quantity <= maxQuantity; quantity++ {
		if isSuitableConfig(combinations, priorities, divider, quantity, limit) {
			return quantity
		}
	}

	return 0
}

// Picks up the maximum quantity of handlers for which the division error does not
// exceed the limit, specified as a percentage.
func PickUpMaxSuitableQuantity(
	priorities []uint,
	divider Divider,
	maxQuantity uint,
	limit float64,
) uint {
	priorities = createSortedCopy(priorities)

	combinations := genPriorityCombinations(priorities)

	for quantity := maxQuantity; quantity != 0; quantity-- {
		if isSuitableConfig(combinations, priorities, divider, quantity, limit) {
			return quantity
		}
	}

	return 0
}
