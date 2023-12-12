// Several auxiliary functions for pickup and checking the quantity of
// handlers are implemented here
package utils

import (
	"fmt"
	"math"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/priority/divider"
	"github.com/akramarenkov/cqos/v2/priority/internal/common"
)

const (
	referenceFactor = 1000
)

// if the number of combinations for n priorities is m, then
// for n+1 priorities the number of combinations is 2m+1
// accordingly, the increment is m+1
func calcCombinationsQuantitySlow(priorities []uint) int {
	quantity := 0

	for length := 0; length < len(priorities); length++ {
		quantity += quantity + 1
	}

	return quantity
}

// it is easy to see that this corresponds to the function 2^n - 1
func calcCombinationsQuantity(priorities []uint) int {
	const base = 2

	return int(math.Pow(base, float64(len(priorities)))) - 1
}

func genPriorityCombinations2(priorities []uint) [][]uint {
	combinations := make([][]uint, 0, calcCombinationsQuantity(priorities))

	for window := len(priorities); window != 0; window-- {
		for shift := 0; shift <= len(priorities)-window; shift++ {
			base := priorities[shift : shift+window]
			remainder := priorities[shift+window:]

			combinations = append(combinations, createCopy(base))

			for _, substituted := range remainder {
				fmt.Println(base, base[1:], substituted)

				for id := range base[1:] {
					copied := createCopy(base)

					copied[id+1] = substituted

					common.SortPriorities(copied)

					combinations = append(combinations, copied)
				}
			}
		}
	}

	return combinations
}

func createCopy(src []uint) []uint {
	newed := make([]uint, len(src))

	copy(newed, src)

	return newed
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

	common.SortPriorities(newed)

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
	divider divider.Divider,
	quantity uint,
) bool {
	for _, combination := range combinations {
		distribution := make(map[uint]uint)

		divider(combination, quantity, distribution)

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
	divider divider.Divider,
	quantity uint,
) bool {
	combinations := genPriorityCombinations(priorities)

	return isNonFatalConfig(combinations, divider, quantity)
}

// Picks up the minimum quantity of handlers for which the division error does not
// cause stop processing of one or more priorities
func PickUpMinNonFatalQuantity(
	priorities []uint,
	divider divider.Divider,
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
	divider divider.Divider,
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

		diff = consts.OneHundredPercent * math.Abs(diff)

		if diff > limit {
			return false
		}
	}

	return true
}

func isSuitableConfig(
	combinations [][]uint,
	priorities []uint,
	divider divider.Divider,
	quantity uint,
	limit float64,
) bool {
	referenceTotalQuantity := referenceFactor * common.SumPriorities(priorities)

	for _, combination := range combinations {
		distribution := make(map[uint]uint)
		reference := make(map[uint]uint)

		divider(combination, quantity, distribution)

		if !isDistributionFilled(distribution) {
			return false
		}

		divider(combination, referenceTotalQuantity, reference)

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
	divider divider.Divider,
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
	divider divider.Divider,
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
	divider divider.Divider,
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
