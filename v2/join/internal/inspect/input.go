package inspect

// Returns a sequence of numbers starting from 1 to value of 'quantity' inclusive
// divided into blocks which should be supplied to the input of the discipline.
func Input(quantity int, blockSize int) [][]int {
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

	for base := 1; base <= quantity; base += blockSize {
		block := make([]int, 0, blockSize)

		for id := range blockSize {
			item := base + id

			if item > quantity {
				break
			}

			block = append(block, item)
		}

		blocks = append(blocks, block)
	}

	return blocks
}
