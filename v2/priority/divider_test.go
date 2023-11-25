package priority

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFairDivider(t *testing.T) {
	priorities := []uint{3, 2, 1}

	distribution := FairDivider(nil, 3, nil)
	require.Equal(t, map[uint]uint(nil), distribution)

	distribution = FairDivider(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0, 2: 0, 1: 0}, distribution)

	distribution = FairDivider(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1, 2: 0, 1: 0}, distribution)

	distribution = FairDivider(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 1, 2: 1, 1: 0}, distribution)

	distribution = FairDivider(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 1, 2: 1, 1: 1}, distribution)

	distribution = FairDivider(priorities, 4, nil)
	require.Equal(t, map[uint]uint{3: 2, 2: 1, 1: 1}, distribution)

	distribution = FairDivider(priorities, 5, nil)
	require.Equal(t, map[uint]uint{3: 2, 2: 2, 1: 1}, distribution)

	distribution = FairDivider(priorities, 6, nil)
	require.Equal(t, map[uint]uint{3: 2, 2: 2, 1: 2}, distribution)

	distribution = FairDivider(priorities, 7, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 2, 1: 2}, distribution)

	distribution = FairDivider(priorities, 8, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 3, 1: 2}, distribution)

	distribution = FairDivider(priorities, 9, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 3, 1: 3}, distribution)

	distribution = FairDivider(priorities, 10, nil)
	require.Equal(t, map[uint]uint{3: 4, 2: 3, 1: 3}, distribution)

	distribution = FairDivider(priorities, 11, nil)
	require.Equal(t, map[uint]uint{3: 4, 2: 4, 1: 3}, distribution)

	distribution = FairDivider(priorities, 12, nil)
	require.Equal(t, map[uint]uint{3: 4, 2: 4, 1: 4}, distribution)

	distribution = FairDivider(priorities, 100, nil)
	require.Equal(t, map[uint]uint{3: 34, 2: 33, 1: 33}, distribution)
}

func TestFairDividerEven(t *testing.T) {
	priorities := []uint{4, 3, 2, 1}

	distribution := FairDivider(priorities, 0, nil)
	require.Equal(t, map[uint]uint{4: 0, 3: 0, 2: 0, 1: 0}, distribution)

	distribution = FairDivider(priorities, 1, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 0, 2: 0, 1: 0}, distribution)

	distribution = FairDivider(priorities, 2, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 1, 2: 0, 1: 0}, distribution)

	distribution = FairDivider(priorities, 3, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 0}, distribution)

	distribution = FairDivider(priorities, 4, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 1}, distribution)

	distribution = FairDivider(priorities, 5, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 1, 2: 1, 1: 1}, distribution)

	// sequence of low priorities is not monotonic due to rounding in highest priorities
	distribution = FairDivider(priorities, 6, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 0}, distribution)

	distribution = FairDivider(priorities, 7, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 1}, distribution)

	distribution = FairDivider(priorities, 8, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 2}, distribution)

	distribution = FairDivider(priorities, 9, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 2, 2: 2, 1: 2}, distribution)

	// sequence of low priorities is not monotonic due to rounding in highest priorities
	distribution = FairDivider(priorities, 10, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 3, 2: 3, 1: 1}, distribution)

	distribution = FairDivider(priorities, 11, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 3, 2: 3, 1: 2}, distribution)

	distribution = FairDivider(priorities, 12, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 3, 2: 3, 1: 3}, distribution)

	distribution = FairDivider(priorities, 100, nil)
	require.Equal(t, map[uint]uint{4: 25, 3: 25, 2: 25, 1: 25}, distribution)
}

func TestFairDividerSingle(t *testing.T) {
	priorities := []uint{3}

	distribution := FairDivider(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0}, distribution)

	distribution = FairDivider(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1}, distribution)

	distribution = FairDivider(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 2}, distribution)

	distribution = FairDivider(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 3}, distribution)
}

func TestFairDividerAdd(t *testing.T) {
	priorities := []uint{3, 2, 1}

	distribution := map[uint]uint{3: 0, 1: 0}

	FairDivider(priorities, 0, distribution)
	require.Equal(t, map[uint]uint{3: 0, 2: 0, 1: 0}, distribution)

	FairDivider(priorities, 1, distribution)
	require.Equal(t, map[uint]uint{3: 1, 2: 0, 1: 0}, distribution)

	FairDivider(priorities, 2, distribution)
	require.Equal(t, map[uint]uint{3: 2, 2: 1, 1: 0}, distribution)

	FairDivider(priorities, 3, distribution)
	require.Equal(t, map[uint]uint{3: 3, 2: 2, 1: 1}, distribution)

	FairDivider(priorities, 6, distribution)
	require.Equal(t, map[uint]uint{3: 5, 2: 4, 1: 3}, distribution)

	FairDivider(priorities, 9, distribution)
	require.Equal(t, map[uint]uint{3: 8, 2: 7, 1: 6}, distribution)

	FairDivider(priorities[1:], 9, distribution)
	require.Equal(t, map[uint]uint{3: 8, 2: 12, 1: 10}, distribution)

	FairDivider(priorities[1:], 10, distribution)
	require.Equal(t, map[uint]uint{3: 8, 2: 17, 1: 15}, distribution)

	FairDivider(priorities[2:], 10, distribution)
	require.Equal(t, map[uint]uint{3: 8, 2: 17, 1: 25}, distribution)
}

func TestFairDividerDiscontinuous(t *testing.T) {
	priorities := []uint{3, 1}

	distribution := FairDivider(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0, 1: 0}, distribution)

	distribution = FairDivider(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1, 1: 0}, distribution)

	distribution = FairDivider(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 1, 1: 1}, distribution)

	distribution = FairDivider(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 2, 1: 1}, distribution)

	distribution = FairDivider(priorities, 4, nil)
	require.Equal(t, map[uint]uint{3: 2, 1: 2}, distribution)

	distribution = FairDivider(priorities, 5, nil)
	require.Equal(t, map[uint]uint{3: 3, 1: 2}, distribution)

	distribution = FairDivider(priorities, 6, nil)
	require.Equal(t, map[uint]uint{3: 3, 1: 3}, distribution)

	distribution = FairDivider(priorities, 7, nil)
	require.Equal(t, map[uint]uint{3: 4, 1: 3}, distribution)

	distribution = FairDivider(priorities, 8, nil)
	require.Equal(t, map[uint]uint{3: 4, 1: 4}, distribution)
}

func TestFairDividerError(t *testing.T) {
	distribution := FairDivider([]uint{6, 5, 4, 3, 2, 1}, 6, nil)
	require.Equal(t, map[uint]uint{6: 1, 5: 1, 4: 1, 3: 1, 2: 1, 1: 1}, distribution)

	distribution = FairDivider([]uint{5, 4, 3, 2, 1}, 6, nil)
	require.Equal(t, map[uint]uint{5: 2, 4: 1, 3: 1, 2: 1, 1: 1}, distribution)

	// Fatal dividing error - values for one or more priorities are zero
	// They also occurs because of the small value of the dividend
	distribution = FairDivider([]uint{4, 3, 2, 1}, 6, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 0}, distribution)

	// At large values of the dividend values of the distribution are no longer zero
	distribution = FairDivider([]uint{4, 3, 2, 1}, 12, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 3, 2: 3, 1: 3}, distribution)

	// At large values of the dividend values of the distribution are no longer zero
	distribution = FairDivider([]uint{4, 3, 2, 1}, 60, nil)
	require.Equal(t, map[uint]uint{4: 15, 3: 15, 2: 15, 1: 15}, distribution)

	// Non-fatal dividing error - poor proportions
	// Occurs because of the small value of the dividend
	distribution = FairDivider([]uint{4, 3, 2, 1}, 7, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 2, 1: 1}, distribution)

	// At larger values of the dividend, the proportions differ not so significantly
	distribution = FairDivider([]uint{4, 3, 2, 1}, 70, nil)
	require.Equal(t, map[uint]uint{4: 18, 3: 18, 2: 18, 1: 16}, distribution)
}

func TestRateDivider(t *testing.T) {
	priorities := []uint{3, 2, 1}

	distribution := RateDivider(nil, 3, nil)
	require.Equal(t, map[uint]uint(nil), distribution)

	distribution = RateDivider(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0, 2: 0, 1: 0}, distribution)

	distribution = RateDivider(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1, 2: 0, 1: 0}, distribution)

	distribution = RateDivider(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 1, 2: 1, 1: 0}, distribution)

	distribution = RateDivider(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 2, 2: 1, 1: 0}, distribution)

	distribution = RateDivider(priorities, 4, nil)
	require.Equal(t, map[uint]uint{3: 2, 2: 1, 1: 1}, distribution)

	// sequence of low priorities is not monotonic due to rounding in highest priorities
	distribution = RateDivider(priorities, 5, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 2, 1: 0}, distribution)

	distribution = RateDivider(priorities, 6, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 2, 1: 1}, distribution)

	distribution = RateDivider(priorities, 7, nil)
	require.Equal(t, map[uint]uint{3: 4, 2: 2, 1: 1}, distribution)

	distribution = RateDivider(priorities, 8, nil)
	require.Equal(t, map[uint]uint{3: 4, 2: 3, 1: 1}, distribution)

	distribution = RateDivider(priorities, 9, nil)
	require.Equal(t, map[uint]uint{3: 5, 2: 3, 1: 1}, distribution)

	distribution = RateDivider(priorities, 10, nil)
	require.Equal(t, map[uint]uint{3: 5, 2: 3, 1: 2}, distribution)

	// sequence of low priorities is not monotonic due to rounding in highest priorities
	distribution = RateDivider(priorities, 11, nil)
	require.Equal(t, map[uint]uint{3: 6, 2: 4, 1: 1}, distribution)

	distribution = RateDivider(priorities, 12, nil)
	require.Equal(t, map[uint]uint{3: 6, 2: 4, 1: 2}, distribution)

	distribution = RateDivider(priorities, 100, nil)
	require.Equal(t, map[uint]uint{3: 50, 2: 33, 1: 17}, distribution)
}

func TestRateDividerEven(t *testing.T) {
	priorities := []uint{4, 3, 2, 1}

	distribution := RateDivider(priorities, 0, nil)
	require.Equal(t, map[uint]uint{4: 0, 3: 0, 2: 0, 1: 0}, distribution)

	distribution = RateDivider(priorities, 1, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 0, 2: 0, 1: 0}, distribution)

	distribution = RateDivider(priorities, 2, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 1, 2: 0, 1: 0}, distribution)

	distribution = RateDivider(priorities, 3, nil)
	require.Equal(t, map[uint]uint{4: 1, 3: 1, 2: 1, 1: 0}, distribution)

	distribution = RateDivider(priorities, 4, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 1, 2: 1, 1: 0}, distribution)

	distribution = RateDivider(priorities, 5, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 1, 1: 0}, distribution)

	distribution = RateDivider(priorities, 6, nil)
	require.Equal(t, map[uint]uint{4: 2, 3: 2, 2: 1, 1: 1}, distribution)

	distribution = RateDivider(priorities, 7, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 2, 2: 1, 1: 1}, distribution)

	distribution = RateDivider(priorities, 8, nil)
	require.Equal(t, map[uint]uint{4: 3, 3: 2, 2: 2, 1: 1}, distribution)

	// sequence of low priorities is not monotonic due to rounding in highest priorities
	distribution = RateDivider(priorities, 9, nil)
	require.Equal(t, map[uint]uint{4: 4, 3: 3, 2: 2, 1: 0}, distribution)

	distribution = RateDivider(priorities, 10, nil)
	require.Equal(t, map[uint]uint{4: 4, 3: 3, 2: 2, 1: 1}, distribution)

	distribution = RateDivider(priorities, 11, nil)
	require.Equal(t, map[uint]uint{4: 5, 3: 3, 2: 2, 1: 1}, distribution)

	distribution = RateDivider(priorities, 12, nil)
	require.Equal(t, map[uint]uint{4: 5, 3: 4, 2: 2, 1: 1}, distribution)

	distribution = RateDivider(priorities, 100, nil)
	require.Equal(t, map[uint]uint{4: 40, 3: 30, 2: 20, 1: 10}, distribution)
}

func TestRateDividerSingle(t *testing.T) {
	priorities := []uint{3}

	distribution := RateDivider(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0}, distribution)

	distribution = RateDivider(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1}, distribution)

	distribution = RateDivider(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 2}, distribution)

	distribution = RateDivider(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 3}, distribution)
}

func TestRateDividerAdd(t *testing.T) {
	priorities := []uint{3, 2, 1}

	distribution := map[uint]uint{3: 0, 1: 0}

	RateDivider(priorities, 0, distribution)
	require.Equal(t, map[uint]uint{3: 0, 2: 0, 1: 0}, distribution)

	RateDivider(priorities, 1, distribution)
	require.Equal(t, map[uint]uint{3: 1, 2: 0, 1: 0}, distribution)

	RateDivider(priorities, 2, distribution)
	require.Equal(t, map[uint]uint{3: 2, 2: 1, 1: 0}, distribution)

	RateDivider(priorities, 3, distribution)
	require.Equal(t, map[uint]uint{3: 4, 2: 2, 1: 0}, distribution)

	RateDivider(priorities, 6, distribution)
	require.Equal(t, map[uint]uint{3: 7, 2: 4, 1: 1}, distribution)

	RateDivider(priorities, 9, distribution)
	require.Equal(t, map[uint]uint{3: 12, 2: 7, 1: 2}, distribution)

	RateDivider(priorities[1:], 9, distribution)
	require.Equal(t, map[uint]uint{3: 12, 2: 13, 1: 5}, distribution)

	RateDivider(priorities[1:], 10, distribution)
	require.Equal(t, map[uint]uint{3: 12, 2: 20, 1: 8}, distribution)

	RateDivider(priorities[2:], 10, distribution)
	require.Equal(t, map[uint]uint{3: 12, 2: 20, 1: 18}, distribution)
}

func TestRateDividerDiscontinuous(t *testing.T) {
	priorities := []uint{3, 1}

	distribution := RateDivider(priorities, 0, nil)
	require.Equal(t, map[uint]uint{3: 0, 1: 0}, distribution)

	distribution = RateDivider(priorities, 1, nil)
	require.Equal(t, map[uint]uint{3: 1, 1: 0}, distribution)

	distribution = RateDivider(priorities, 2, nil)
	require.Equal(t, map[uint]uint{3: 2, 1: 0}, distribution)

	distribution = RateDivider(priorities, 3, nil)
	require.Equal(t, map[uint]uint{3: 2, 1: 1}, distribution)

	distribution = RateDivider(priorities, 4, nil)
	require.Equal(t, map[uint]uint{3: 3, 1: 1}, distribution)

	distribution = RateDivider(priorities, 5, nil)
	require.Equal(t, map[uint]uint{3: 4, 1: 1}, distribution)

	distribution = RateDivider(priorities, 6, nil)
	require.Equal(t, map[uint]uint{3: 5, 1: 1}, distribution)

	distribution = RateDivider(priorities, 7, nil)
	require.Equal(t, map[uint]uint{3: 5, 1: 2}, distribution)

	distribution = RateDivider(priorities, 8, nil)
	require.Equal(t, map[uint]uint{3: 6, 1: 2}, distribution)

	distribution = RateDivider(priorities, 9, nil)
	require.Equal(t, map[uint]uint{3: 7, 1: 2}, distribution)

	distribution = RateDivider(priorities, 10, nil)
	require.Equal(t, map[uint]uint{3: 8, 1: 2}, distribution)

	distribution = RateDivider(priorities, 11, nil)
	require.Equal(t, map[uint]uint{3: 8, 1: 3}, distribution)

	distribution = RateDivider(priorities, 12, nil)
	require.Equal(t, map[uint]uint{3: 9, 1: 3}, distribution)

	distribution = RateDivider(priorities, 13, nil)
	require.Equal(t, map[uint]uint{3: 10, 1: 3}, distribution)

	distribution = RateDivider(priorities, 14, nil)
	require.Equal(t, map[uint]uint{3: 11, 1: 3}, distribution)

	distribution = RateDivider(priorities, 15, nil)
	require.Equal(t, map[uint]uint{3: 11, 1: 4}, distribution)
}

func TestRateDividerError(t *testing.T) {
	distribution := RateDivider([]uint{3, 2, 1}, 6, nil)
	require.Equal(t, map[uint]uint{3: 3, 2: 2, 1: 1}, distribution)

	distribution = RateDivider([]uint{2, 1}, 6, nil)
	require.Equal(t, map[uint]uint{2: 4, 1: 2}, distribution)

	// Non-fatal dividing error - poor proportions
	// Must be 3:1, but returned 5:1
	// Occurs because of the small value of the dividend
	distribution = RateDivider([]uint{3, 1}, 6, nil)
	require.Equal(t, map[uint]uint{3: 5, 1: 1}, distribution)

	// At larger values of the dividend, the proportions are restored
	distribution = RateDivider([]uint{3, 1}, 60, nil)
	require.Equal(t, map[uint]uint{3: 45, 1: 15}, distribution)

	// Fatal dividing error - values for one or more priorities are zero
	// They also occurs because of the small value of the dividend
	distribution = RateDivider([]uint{3, 1}, 2, nil)
	require.Equal(t, map[uint]uint{3: 2, 1: 0}, distribution)

	// At large values of the dividend values of the distribution are no longer zero
	distribution = RateDivider([]uint{3, 1}, 20, nil)
	require.Equal(t, map[uint]uint{3: 15, 1: 5}, distribution)
}

func TestRateDividerLifeHack(t *testing.T) {
	priorities := []uint{70, 20, 10}

	distribution := RateDivider(priorities, 100, nil)
	require.Equal(t, map[uint]uint{70: 70, 20: 20, 10: 10}, distribution)
}
