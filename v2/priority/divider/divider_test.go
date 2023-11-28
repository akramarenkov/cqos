package divider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFairDivider(t *testing.T) {
	priorities := []uint{3, 2, 1}

	distribution := Fair(nil, 3, nil)
	require.Equal(t, map[uint]uint(nil), distribution)

	distribution = Fair(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0, 2: 0, 1: 0}, distribution)

	distribution = Fair(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1, 2: 0, 1: 0}, distribution)

	distribution = Fair(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 1, 2: 1, 1: 0}, distribution)

	distribution = Fair(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 1, 2: 1, 1: 1}, distribution)

	distribution = Fair(priorities, 4, nil)
	require.Equal(t, map[uint]uint{3: 2, 2: 1, 1: 1}, distribution)

	distribution = Fair(priorities, 5, nil)
	require.Equal(t, map[uint]uint{3: 2, 2: 2, 1: 1}, distribution)

	distribution = Fair(priorities, 6, nil)
	require.Equal(t, map[uint]uint{3: 2, 2: 2, 1: 2}, distribution)

	distribution = Fair(priorities, 7, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 2, 1: 2}, distribution)

	distribution = Fair(priorities, 8, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 3, 1: 2}, distribution)

	distribution = Fair(priorities, 9, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 3, 1: 3}, distribution)

	distribution = Fair(priorities, 10, nil)
	require.Equal(t, map[uint]uint{3: 4, 2: 3, 1: 3}, distribution)

	distribution = Fair(priorities, 11, nil)
	require.Equal(t, map[uint]uint{3: 4, 2: 4, 1: 3}, distribution)

	distribution = Fair(priorities, 12, nil)
	require.Equal(t, map[uint]uint{3: 4, 2: 4, 1: 4}, distribution)

	distribution = Fair(priorities, 100, nil)
	require.Equal(t, map[uint]uint{3: 34, 2: 33, 1: 33}, distribution)
}

func TestFairDividerEven(t *testing.T) {
	priorities := []uint{4, 3, 2, 1}

	distribution := Fair(priorities, 0, nil)
	require.Equal(t, map[uint]uint{4: 0, 3: 0, 2: 0, 1: 0}, distribution)

	distribution = Fair(priorities, 1, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 0, 2: 0, 1: 0}, distribution)

	distribution = Fair(priorities, 2, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 1, 2: 0, 1: 0}, distribution)

	distribution = Fair(priorities, 3, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 0}, distribution)

	distribution = Fair(priorities, 4, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 1}, distribution)

	distribution = Fair(priorities, 5, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 1, 2: 1, 1: 1}, distribution)

	// sequence of low priorities is not monotonic due to rounding in highest priorities
	distribution = Fair(priorities, 6, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 0}, distribution)

	distribution = Fair(priorities, 7, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 1}, distribution)

	distribution = Fair(priorities, 8, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 2}, distribution)

	distribution = Fair(priorities, 9, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 2, 2: 2, 1: 2}, distribution)

	// sequence of low priorities is not monotonic due to rounding in highest priorities
	distribution = Fair(priorities, 10, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 3, 2: 3, 1: 1}, distribution)

	distribution = Fair(priorities, 11, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 3, 2: 3, 1: 2}, distribution)

	distribution = Fair(priorities, 12, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 3, 2: 3, 1: 3}, distribution)

	distribution = Fair(priorities, 100, nil)
	require.Equal(t, map[uint]uint{4: 25, 3: 25, 2: 25, 1: 25}, distribution)
}

func TestFairDividerSingle(t *testing.T) {
	priorities := []uint{3}

	distribution := Fair(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0}, distribution)

	distribution = Fair(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1}, distribution)

	distribution = Fair(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 2}, distribution)

	distribution = Fair(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 3}, distribution)
}

func TestFairDividerAdd(t *testing.T) {
	priorities := []uint{3, 2, 1}

	distribution := map[uint]uint{3: 0, 1: 0}

	Fair(priorities, 0, distribution)
	require.Equal(t, map[uint]uint{3: 0, 2: 0, 1: 0}, distribution)

	Fair(priorities, 1, distribution)
	require.Equal(t, map[uint]uint{3: 1, 2: 0, 1: 0}, distribution)

	Fair(priorities, 2, distribution)
	require.Equal(t, map[uint]uint{3: 2, 2: 1, 1: 0}, distribution)

	Fair(priorities, 3, distribution)
	require.Equal(t, map[uint]uint{3: 3, 2: 2, 1: 1}, distribution)

	Fair(priorities, 6, distribution)
	require.Equal(t, map[uint]uint{3: 5, 2: 4, 1: 3}, distribution)

	Fair(priorities, 9, distribution)
	require.Equal(t, map[uint]uint{3: 8, 2: 7, 1: 6}, distribution)

	Fair(priorities[1:], 9, distribution)
	require.Equal(t, map[uint]uint{3: 8, 2: 12, 1: 10}, distribution)

	Fair(priorities[1:], 10, distribution)
	require.Equal(t, map[uint]uint{3: 8, 2: 17, 1: 15}, distribution)

	Fair(priorities[2:], 10, distribution)
	require.Equal(t, map[uint]uint{3: 8, 2: 17, 1: 25}, distribution)
}

func TestFairDividerDiscontinuous(t *testing.T) {
	priorities := []uint{3, 1}

	distribution := Fair(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0, 1: 0}, distribution)

	distribution = Fair(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1, 1: 0}, distribution)

	distribution = Fair(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 1, 1: 1}, distribution)

	distribution = Fair(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 2, 1: 1}, distribution)

	distribution = Fair(priorities, 4, nil)
	require.Equal(t, map[uint]uint{3: 2, 1: 2}, distribution)

	distribution = Fair(priorities, 5, nil)
	require.Equal(t, map[uint]uint{3: 3, 1: 2}, distribution)

	distribution = Fair(priorities, 6, nil)
	require.Equal(t, map[uint]uint{3: 3, 1: 3}, distribution)

	distribution = Fair(priorities, 7, nil)
	require.Equal(t, map[uint]uint{3: 4, 1: 3}, distribution)

	distribution = Fair(priorities, 8, nil)
	require.Equal(t, map[uint]uint{3: 4, 1: 4}, distribution)
}

func TestFairDividerError(t *testing.T) {
	distribution := Fair([]uint{6, 5, 4, 3, 2, 1}, 6, nil)
	require.Equal(t, map[uint]uint{6: 1, 5: 1, 4: 1, 3: 1, 2: 1, 1: 1}, distribution)

	distribution = Fair([]uint{5, 4, 3, 2, 1}, 6, nil)
	require.Equal(t, map[uint]uint{5: 2, 4: 1, 3: 1, 2: 1, 1: 1}, distribution)

	// Fatal dividing error - values for one or more priorities are zero
	// They also occurs because of the small value of the dividend
	distribution = Fair([]uint{4, 3, 2, 1}, 6, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 0}, distribution)

	// At large values of the dividend values of the distribution are no longer zero
	distribution = Fair([]uint{4, 3, 2, 1}, 12, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 3, 2: 3, 1: 3}, distribution)

	// At large values of the dividend values of the distribution are no longer zero
	distribution = Fair([]uint{4, 3, 2, 1}, 60, nil)
	require.Equal(t, map[uint]uint{4: 15, 3: 15, 2: 15, 1: 15}, distribution)

	// Non-fatal dividing error - poor proportions
	// Occurs because of the small value of the dividend
	distribution = Fair([]uint{4, 3, 2, 1}, 7, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 1}, distribution)

	// At larger values of the dividend, the proportions differ not so significantly
	distribution = Fair([]uint{4, 3, 2, 1}, 70, nil)
	require.Equal(t, map[uint]uint{4: 18, 3: 18, 2: 18, 1: 16}, distribution)
}

func TestRateDivider(t *testing.T) {
	priorities := []uint{3, 2, 1}

	distribution := Rate(nil, 3, nil)
	require.Equal(t, map[uint]uint(nil), distribution)

	distribution = Rate(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0, 2: 0, 1: 0}, distribution)

	distribution = Rate(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1, 2: 0, 1: 0}, distribution)

	distribution = Rate(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 1, 2: 1, 1: 0}, distribution)

	distribution = Rate(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 2, 2: 1, 1: 0}, distribution)

	distribution = Rate(priorities, 4, nil)
	require.Equal(t, map[uint]uint{3: 2, 2: 1, 1: 1}, distribution)

	// sequence of low priorities is not monotonic due to rounding in highest priorities
	distribution = Rate(priorities, 5, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 2, 1: 0}, distribution)

	distribution = Rate(priorities, 6, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 2, 1: 1}, distribution)

	distribution = Rate(priorities, 7, nil)
	require.Equal(t, map[uint]uint{3: 4, 2: 2, 1: 1}, distribution)

	distribution = Rate(priorities, 8, nil)
	require.Equal(t, map[uint]uint{3: 4, 2: 3, 1: 1}, distribution)

	distribution = Rate(priorities, 9, nil)
	require.Equal(t, map[uint]uint{3: 5, 2: 3, 1: 1}, distribution)

	distribution = Rate(priorities, 10, nil)
	require.Equal(t, map[uint]uint{3: 5, 2: 3, 1: 2}, distribution)

	// sequence of low priorities is not monotonic due to rounding in highest priorities
	distribution = Rate(priorities, 11, nil)
	require.Equal(t, map[uint]uint{3: 6, 2: 4, 1: 1}, distribution)

	distribution = Rate(priorities, 12, nil)
	require.Equal(t, map[uint]uint{3: 6, 2: 4, 1: 2}, distribution)

	distribution = Rate(priorities, 100, nil)
	require.Equal(t, map[uint]uint{3: 50, 2: 33, 1: 17}, distribution)
}

func TestRateDividerEven(t *testing.T) {
	priorities := []uint{4, 3, 2, 1}

	distribution := Rate(priorities, 0, nil)
	require.Equal(t, map[uint]uint{4: 0, 3: 0, 2: 0, 1: 0}, distribution)

	distribution = Rate(priorities, 1, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 0, 2: 0, 1: 0}, distribution)

	distribution = Rate(priorities, 2, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 1, 2: 0, 1: 0}, distribution)

	distribution = Rate(priorities, 3, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 0}, distribution)

	distribution = Rate(priorities, 4, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 1, 2: 1, 1: 0}, distribution)

	distribution = Rate(priorities, 5, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 1, 1: 0}, distribution)

	distribution = Rate(priorities, 6, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 1, 1: 1}, distribution)

	distribution = Rate(priorities, 7, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 2, 2: 1, 1: 1}, distribution)

	distribution = Rate(priorities, 8, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 2, 2: 2, 1: 1}, distribution)

	// sequence of low priorities is not monotonic due to rounding in highest priorities
	distribution = Rate(priorities, 9, nil)
	require.Equal(t, map[uint]uint{4: 4, 3: 3, 2: 2, 1: 0}, distribution)

	distribution = Rate(priorities, 10, nil)
	require.Equal(t, map[uint]uint{4: 4, 3: 3, 2: 2, 1: 1}, distribution)

	distribution = Rate(priorities, 11, nil)
	require.Equal(t, map[uint]uint{4: 5, 3: 3, 2: 2, 1: 1}, distribution)

	distribution = Rate(priorities, 12, nil)
	require.Equal(t, map[uint]uint{4: 5, 3: 4, 2: 2, 1: 1}, distribution)

	distribution = Rate(priorities, 100, nil)
	require.Equal(t, map[uint]uint{4: 40, 3: 30, 2: 20, 1: 10}, distribution)
}

func TestRateDividerSingle(t *testing.T) {
	priorities := []uint{3}

	distribution := Rate(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0}, distribution)

	distribution = Rate(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1}, distribution)

	distribution = Rate(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 2}, distribution)

	distribution = Rate(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 3}, distribution)
}

func TestRateDividerAdd(t *testing.T) {
	priorities := []uint{3, 2, 1}

	distribution := map[uint]uint{3: 0, 1: 0}

	Rate(priorities, 0, distribution)
	require.Equal(t, map[uint]uint{3: 0, 2: 0, 1: 0}, distribution)

	Rate(priorities, 1, distribution)
	require.Equal(t, map[uint]uint{3: 1, 2: 0, 1: 0}, distribution)

	Rate(priorities, 2, distribution)
	require.Equal(t, map[uint]uint{3: 2, 2: 1, 1: 0}, distribution)

	Rate(priorities, 3, distribution)
	require.Equal(t, map[uint]uint{3: 4, 2: 2, 1: 0}, distribution)

	Rate(priorities, 6, distribution)
	require.Equal(t, map[uint]uint{3: 7, 2: 4, 1: 1}, distribution)

	Rate(priorities, 9, distribution)
	require.Equal(t, map[uint]uint{3: 12, 2: 7, 1: 2}, distribution)

	Rate(priorities[1:], 9, distribution)
	require.Equal(t, map[uint]uint{3: 12, 2: 13, 1: 5}, distribution)

	Rate(priorities[1:], 10, distribution)
	require.Equal(t, map[uint]uint{3: 12, 2: 20, 1: 8}, distribution)

	Rate(priorities[2:], 10, distribution)
	require.Equal(t, map[uint]uint{3: 12, 2: 20, 1: 18}, distribution)
}

func TestRateDividerDiscontinuous(t *testing.T) {
	priorities := []uint{3, 1}

	distribution := Rate(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0, 1: 0}, distribution)

	distribution = Rate(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1, 1: 0}, distribution)

	distribution = Rate(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 2, 1: 0}, distribution)

	distribution = Rate(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 2, 1: 1}, distribution)

	distribution = Rate(priorities, 4, nil)
	require.Equal(t, map[uint]uint{3: 3, 1: 1}, distribution)

	distribution = Rate(priorities, 5, nil)
	require.Equal(t, map[uint]uint{3: 4, 1: 1}, distribution)

	distribution = Rate(priorities, 6, nil)
	require.Equal(t, map[uint]uint{3: 5, 1: 1}, distribution)

	distribution = Rate(priorities, 7, nil)
	require.Equal(t, map[uint]uint{3: 5, 1: 2}, distribution)

	distribution = Rate(priorities, 8, nil)
	require.Equal(t, map[uint]uint{3: 6, 1: 2}, distribution)

	distribution = Rate(priorities, 9, nil)
	require.Equal(t, map[uint]uint{3: 7, 1: 2}, distribution)

	distribution = Rate(priorities, 10, nil)
	require.Equal(t, map[uint]uint{3: 8, 1: 2}, distribution)

	distribution = Rate(priorities, 11, nil)
	require.Equal(t, map[uint]uint{3: 8, 1: 3}, distribution)

	distribution = Rate(priorities, 12, nil)
	require.Equal(t, map[uint]uint{3: 9, 1: 3}, distribution)

	distribution = Rate(priorities, 13, nil)
	require.Equal(t, map[uint]uint{3: 10, 1: 3}, distribution)

	distribution = Rate(priorities, 14, nil)
	require.Equal(t, map[uint]uint{3: 11, 1: 3}, distribution)

	distribution = Rate(priorities, 15, nil)
	require.Equal(t, map[uint]uint{3: 11, 1: 4}, distribution)
}

func TestRateDividerError(t *testing.T) {
	distribution := Rate([]uint{3, 2, 1}, 6, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 2, 1: 1}, distribution)

	distribution = Rate([]uint{2, 1}, 6, nil)
	require.Equal(t, map[uint]uint{2: 4, 1: 2}, distribution)

	// Non-fatal dividing error - poor proportions
	// Must be 3:1, but returned 5:1
	// Occurs because of the small value of the dividend
	distribution = Rate([]uint{3, 1}, 6, nil)
	require.Equal(t, map[uint]uint{3: 5, 1: 1}, distribution)

	// At larger values of the dividend, the proportions are restored
	distribution = Rate([]uint{3, 1}, 60, nil)
	require.Equal(t, map[uint]uint{3: 45, 1: 15}, distribution)

	// Fatal dividing error - values for one or more priorities are zero
	// They also occurs because of the small value of the dividend
	distribution = Rate([]uint{3, 1}, 2, nil)
	require.Equal(t, map[uint]uint{3: 2, 1: 0}, distribution)

	// At large values of the dividend values of the distribution are no longer zero
	distribution = Rate([]uint{3, 1}, 20, nil)
	require.Equal(t, map[uint]uint{3: 15, 1: 5}, distribution)
}

func TestRateDividerLifeHack(t *testing.T) {
	priorities := []uint{70, 20, 10}

	distribution := Rate(priorities, 100, nil)
	require.Equal(t, map[uint]uint{70: 70, 20: 20, 10: 10}, distribution)
}
