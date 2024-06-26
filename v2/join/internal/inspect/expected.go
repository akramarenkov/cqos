package inspect

type description struct {
	EffectiveJoinSize int
	EffectiveQuantity int
	Joins             int
	RemainderQuantity int
	UnusedJoinSize    int
}

// Returns slices expected from the discipline output channel.
func Expected(quantity int, blockSize int, joinSize uint) [][]int {
	return genExpected(1, calcDescription(quantity, blockSize, joinSize))
}

// Returns slices expected from the discipline output channel if there is one delay
// in receiving a data set blocks leading to the formation of a timeout in the
// discipline during accumulating some output slice.
func ExpectedWithTimeout(
	quantity int,
	pauseAt int,
	blockSize int,
	joinSize uint,
) [][]int {
	descs := calcDescriptionWithTimeout(quantity, pauseAt, blockSize, joinSize)

	blocks := make([][]int, 0, len(descs))

	begin := 1

	for _, desc := range descs {
		blocks = append(blocks, genExpected(begin, desc)...)

		begin += desc.EffectiveQuantity + desc.RemainderQuantity
	}

	return blocks
}

func calcDescription(quantity int, blockSize int, joinSize uint) description {
	if quantity == 0 {
		return description{}
	}

	if blockSize == 0 {
		return description{}
	}

	if joinSize == 0 {
		return description{}
	}

	desc := description{}

	blocksInJoin := int(joinSize) / blockSize

	desc.EffectiveJoinSize = blocksInJoin * blockSize
	desc.UnusedJoinSize = int(joinSize) - desc.EffectiveJoinSize

	if blocksInJoin == 0 {
		desc.EffectiveJoinSize = blockSize
		desc.UnusedJoinSize = 0
	}

	desc.Joins = quantity / desc.EffectiveJoinSize

	desc.EffectiveQuantity = desc.Joins * desc.EffectiveJoinSize
	desc.RemainderQuantity = quantity - desc.EffectiveQuantity

	if desc.Joins == 0 {
		desc.Joins = 1
		desc.EffectiveQuantity = quantity
		desc.RemainderQuantity = 0
	}

	if desc.RemainderQuantity > desc.UnusedJoinSize {
		desc.Joins++
	}

	return desc
}

func calcDescriptionWithTimeout(
	quantity int,
	pauseAt int,
	blockSize int,
	joinSize uint,
) []description {
	if blockSize == 0 {
		return nil
	}

	if pauseAt == 0 {
		return []description{calcDescription(quantity, blockSize, joinSize)}
	}

	if pauseAt > quantity {
		return []description{calcDescription(quantity, blockSize, joinSize)}
	}

	blocksInPauseAt := pauseAt / blockSize
	beforePauseAt := blocksInPauseAt * blockSize

	if beforePauseAt == pauseAt {
		beforePauseAt -= blockSize
	}

	descs := []description{
		calcDescription(beforePauseAt, blockSize, joinSize),
		calcDescription(quantity-beforePauseAt, blockSize, joinSize),
	}

	return descs
}

func genExpected(begin int, desc description) [][]int {
	blocks := make([][]int, 0, desc.Joins)

	effectiveEnd := desc.EffectiveQuantity + begin - 1

	for item := begin; item <= effectiveEnd; item++ {
		id := (item - begin) % desc.EffectiveJoinSize

		if id == 0 {
			blocks = append(blocks, make([]int, 0, desc.EffectiveJoinSize))
		}

		blocks[len(blocks)-1] = append(blocks[len(blocks)-1], item)
	}

	if desc.RemainderQuantity > desc.UnusedJoinSize {
		blocks = append(blocks, make([]int, 0, desc.RemainderQuantity))

		for base := 1; base <= desc.RemainderQuantity; base++ {
			item := base + effectiveEnd

			blocks[len(blocks)-1] = append(blocks[len(blocks)-1], item)
		}

		return blocks
	}

	for base := 1; base <= desc.RemainderQuantity; base++ {
		item := base + effectiveEnd

		blocks[len(blocks)-1] = append(blocks[len(blocks)-1], item)
	}

	return blocks
}
