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

func TestRateDividerLifeHack(t *testing.T) {
	priorities := []uint{75, 20, 5}

	distribution := RateDivider(priorities, 100, nil)
	require.Equal(t, map[uint]uint{75: 75, 20: 20, 5: 5}, distribution)
}