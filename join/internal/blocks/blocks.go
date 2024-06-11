// Internal package used to divide sequence of numbers into blocks and
// predictions of how they will be processed by disciplines. Used for
// testing disciplines.
package blocks

const (
	// It is possible to pause only between blocks of the same join.
	minBlocksNumberToPause = 2
)

// Divides a sequence of numbers starting from 1 to value of 'quantity' inclusive
// into blocks.
func DivideSequence(quantity int, blockSize int) [][]int {
	if quantity == 0 {
		return nil
	}

	if blockSize == 0 {
		return nil
	}

	blocksNumber := quantity / blockSize

	if blocksNumber*blockSize != quantity {
		blocksNumber++
	}

	blocks := make([][]int, 0, blocksNumber)

	for first := 1; first <= quantity; first += blockSize {
		block := make([]int, 0, blockSize)

		for id := range blockSize {
			value := first + id

			if value > quantity {
				break
			}

			block = append(block, value)
		}

		blocks = append(blocks, block)
	}

	return blocks
}

// Calculates the number of output slices of the discipline.
func CalcExpectedJoins(quantity int, blockSize int, joinSize uint) int {
	if quantity == 0 {
		return 0
	}

	if blockSize == 0 {
		return 0
	}

	if joinSize == 0 {
		return 0
	}

	blocksInJoin := int(joinSize) / blockSize
	effectiveJoin := blocksInJoin * blockSize
	unusedJoin := int(joinSize) - effectiveJoin

	if blocksInJoin == 0 {
		effectiveJoin = blockSize
		unusedJoin = 0
	}

	expectedJoins := quantity / effectiveJoin

	if expectedJoins == 0 {
		expectedJoins = 1
	}

	effectiveQuantity := expectedJoins * effectiveJoin
	remainder := quantity - effectiveQuantity

	if remainder > unusedJoin {
		expectedJoins++
	}

	return expectedJoins
}

// Calculates the number of output slices of discipline if there is one delay
// in receiving a blocks leading to the formation of a timeout in the discipline
// during accumulating some output slice.
func CalcExpectedJoinsWithTimeout(
	quantity int,
	pauseAt int,
	blockSize int,
	joinSize uint,
) int {
	if blockSize == 0 {
		return 0
	}

	if pauseAt == 0 {
		return CalcExpectedJoins(quantity, blockSize, joinSize)
	}

	if pauseAt > quantity {
		return CalcExpectedJoins(quantity, blockSize, joinSize)
	}

	blocksInPauseAt := pauseAt / blockSize
	beforePauseAt := blocksInPauseAt * blockSize

	if beforePauseAt == pauseAt {
		beforePauseAt -= blockSize
	}

	expectedJoins := CalcExpectedJoins(beforePauseAt, blockSize, joinSize)
	expectedJoins += CalcExpectedJoins(quantity-beforePauseAt, blockSize, joinSize)

	return expectedJoins
}

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

	estimated += blockSize - scope + 1

	if estimated > quantity {
		return 0
	}

	return estimated
}
