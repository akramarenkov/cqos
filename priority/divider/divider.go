package divider

import "math"

type Divider func(priorities []uint, dividend uint, distribution map[uint]uint) map[uint]uint

func Fair(priorities []uint, dividend uint, distribution map[uint]uint) map[uint]uint {
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

func Rate(priorities []uint, dividend uint, distribution map[uint]uint) map[uint]uint {
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
