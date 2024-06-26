package inspect

import "time"

const (
	// It is possible to pause only between blocks of the same join.
	minBlocksNumberToPause = 2
	// For the reasons for choosing the value, see doc/pauseat-duration-select.svg.
	pausetAtRelativeDuration = 2.75
)

// Pickups correct number of item in sequence of numbers at which you need to
// pause the transmission of input blocks so that the discipline has a timeout and
// she sends an incomplete output slice.
//
// Return value will be located in the block located after the first block in same join.
//
// To pickup the correct value, you must specify the estimated (approximate) value. When
// pickups, the value only increases.
//
// A zero return value means that the number could not be pickuped.
func PickUpPauseAt(quantity int, estimated int, blockSize int, joinSize uint) int {
	if blockSize == 0 {
		return 0
	}

	if joinSize == 0 {
		return 0
	}

	if estimated > quantity {
		return 0
	}

	blocksInJoin := int(joinSize) / blockSize

	if blocksInJoin < minBlocksNumberToPause {
		return 0
	}

	effectiveJoin := blocksInJoin * blockSize
	discarded := estimated / effectiveJoin
	scope := estimated - discarded*effectiveJoin

	if scope == 0 {
		return estimated
	}

	if scope > blockSize {
		return estimated
	}

	pauseAt := estimated + blockSize - scope + 1

	if pauseAt > quantity {
		return 0
	}

	return pauseAt
}

// Calculates pause duration value for timeout tests with which the tests will
// run reliably.
func CalcPauseAtDuration(timeout time.Duration) time.Duration {
	return time.Duration(pausetAtRelativeDuration * float64(timeout))
}
