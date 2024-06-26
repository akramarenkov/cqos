package inspect

// Used in ongoing tests for brevity.
func seq(begin int, end int) []int {
	seq := make([]int, 0, end-begin+1)

	for number := begin; number <= end; number++ {
		seq = append(seq, number)
	}

	return seq
}
