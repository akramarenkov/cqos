package priority

import "math"

// Distributes quantity of something by priorities. Determines how handlers are distributed among priorities.
//
// Slice of priorities is passed to this function sorted from highest to lowest.
//
// Sum of the distributed quantities must equal the original quantity.
//
// If distribution is nil then it must be created and returned, otherwise it must be updated and returned.
type Divider func(priorities []uint, dividend uint, distribution map[uint]uint) map[uint]uint

// Distributes quantity evenly among the priorities.
//
// Used for equaling.
//
// Example results:
//
//   - 6 / [3 2 1] = map[3:2, 2:2, 1:2]
//   - 100 / [70 20 10] = map[70:34, 20:33, 10:33]
func FairDivider(priorities []uint, dividend uint, distribution map[uint]uint) map[uint]uint {
	if len(priorities) == 0 {
		return nil
	}

	if distribution == nil {
		distribution = make(map[uint]uint, len(priorities))
	}

	step := float64(dividend) / float64(len(priorities))
	part := uint(math.Round(step))

	remainder := dividend

	for _, priority := range priorities {
		if remainder < part {
			distribution[priority] += remainder
			remainder = 0

			continue
		}

		distribution[priority] += part

		remainder -= part
	}

	distribution[priorities[0]] += remainder

	return distribution
}

// Distributes quantity between priorities in proportion to the priority value.
//
// Used for prioritization.
//
// Example results:
//
//   - 6 / [3 2 1] = map[3:3, 2:2, 1:1]
//   - 100 / [70 20 10] = map[70:70, 20:20, 10:10]
func RateDivider(priorities []uint, dividend uint, distribution map[uint]uint) map[uint]uint {
	if len(priorities) == 0 {
		return nil
	}

	if distribution == nil {
		distribution = make(map[uint]uint, len(priorities))
	}

	sum := uint(0)

	for _, priority := range priorities {
		sum += priority
	}

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

	return distribution
}
