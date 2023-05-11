package divider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFair(t *testing.T) {
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

func TestFairEven(t *testing.T) {
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

func TestFairSingle(t *testing.T) {
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

func TestFairAdd(t *testing.T) {
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

func TestFairDiscontinuous(t *testing.T) {
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

func TestRate(t *testing.T) {
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

func TestRateEven(t *testing.T) {
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

func TestRateSingle(t *testing.T) {
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

func TestRateAdd(t *testing.T) {
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

func TestRateDiscontinuous(t *testing.T) {
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

func TestRateLifeHack(t *testing.T) {
	priorities := []uint{75, 20, 5}

	distribution := Rate(priorities, 100, nil)
	require.Equal(t, map[uint]uint{75: 75, 20: 20, 5: 5}, distribution)
}
