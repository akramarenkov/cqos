// Several dividers are implemented here that distributes quantity of something by
// priorities
package divider

import (
	"math"

	"github.com/akramarenkov/cqos/v2/priority/internal/common"
)

// Distributes quantity of something by priorities. Determines how handlers are
// distributed among priorities.
//
// Slice of priorities is passed to this function sorted from highest to lowest.
//
// Sum of the distributed quantities must equal the original quantity
type Divider func(priorities []uint, dividend uint, distribution map[uint]uint)

// Distributes quantity evenly among the priorities.
//
// Used for equaling.
//
// Example results:
//
//   - 6 / [3 2 1] = map[3:2, 2:2, 1:2]
//   - 100 / [70 20 10] = map[70:34, 20:33, 10:33]
func Fair(priorities []uint, dividend uint, distribution map[uint]uint) {
	if len(priorities) == 0 {
		return
	}

	if distribution == nil {
		return
	}

	divider := uint(len(priorities))
	base := dividend / divider
	remainder := dividend - base*divider

	// max value of remainder is len(priorities), so we simply increase distribution by one
	for _, priority := range priorities {
		distribution[priority] += base

		if remainder == 0 {
			continue
		}

		distribution[priority]++
		remainder--
	}
}

// Distributes quantity between priorities in proportion to the priority value.
//
// Used for prioritization.
//
// Example results:
//
//   - 6 / [3 2 1] = map[3:3, 2:2, 1:1]
//   - 100 / [70 20 10] = map[70:70, 20:20, 10:10]
func Rate(priorities []uint, dividend uint, distribution map[uint]uint) {
	if len(priorities) == 0 {
		return
	}

	if distribution == nil {
		return
	}

	sum := common.SumPriorities(priorities)

	step := float64(dividend) / float64(sum)

	remainder := dividend

	for _, priority := range priorities {
		part := uint(math.Round(step * float64(priority)))

		if remainder < part {
			distribution[priority] += remainder
			remainder = 0

			continue
		}

		distribution[priority] += part

		remainder -= part
	}

	distribution[priorities[0]] += remainder
}
