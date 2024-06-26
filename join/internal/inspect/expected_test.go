package inspect

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalcDescriptionZeroes(t *testing.T) {
	require.Equal(t, description{}, calcDescription(0, 0, 0))
	require.Equal(t, description{}, calcDescription(0, 4, 0))
	require.Equal(t, description{}, calcDescription(12, 0, 0))
	require.Equal(t, description{}, calcDescription(12, 4, 0))

	require.Equal(t, description{}, calcDescription(0, 0, 10))
	require.Equal(t, description{}, calcDescription(0, 4, 10))
	require.Equal(t, description{}, calcDescription(12, 0, 10))
}

func TestCalcDescriptionBlockSize1(t *testing.T) {
	for quantity := 1; quantity <= 10; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 10,
				EffectiveQuantity: quantity,
				Joins:             1,
				RemainderQuantity: 0,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 1, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 11; quantity <= 30; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 10,
				EffectiveQuantity: 10 * (quantity / 10),
				Joins:             int(math.Ceil(float64(quantity) / 10)),
				RemainderQuantity: quantity % 10,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 1, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestCalcDescriptionBlockSize3(t *testing.T) {
	for quantity := 1; quantity <= 9; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 9,
				EffectiveQuantity: quantity,
				Joins:             1,
				RemainderQuantity: 0,
				UnusedJoinSize:    1,
			},
			calcDescription(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 9,
			EffectiveQuantity: 9,
			Joins:             1,
			RemainderQuantity: 1,
			UnusedJoinSize:    1,
		},
		calcDescription(10, 3, 10),
	)

	for quantity := 11; quantity <= 18; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 9,
				EffectiveQuantity: 9 * (quantity / 9),
				Joins:             int(math.Ceil(float64(quantity) / 9)),
				RemainderQuantity: quantity % 9,
				UnusedJoinSize:    1,
			},
			calcDescription(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 9,
			EffectiveQuantity: 18,
			Joins:             2,
			RemainderQuantity: 1,
			UnusedJoinSize:    1,
		},
		calcDescription(19, 3, 10),
	)

	for quantity := 20; quantity <= 27; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 9,
				EffectiveQuantity: 9 * (quantity / 9),
				Joins:             int(math.Ceil(float64(quantity) / 9)),
				RemainderQuantity: quantity % 9,
				UnusedJoinSize:    1,
			},
			calcDescription(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 9,
			EffectiveQuantity: 27,
			Joins:             3,
			RemainderQuantity: 1,
			UnusedJoinSize:    1,
		},
		calcDescription(28, 3, 10),
	)
}

func TestCalcDescriptionBlockSize4(t *testing.T) {
	for quantity := 1; quantity <= 8; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 8,
				EffectiveQuantity: quantity,
				Joins:             1,
				RemainderQuantity: 0,
				UnusedJoinSize:    2,
			},
			calcDescription(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 8,
			Joins:             1,
			RemainderQuantity: 1,
			UnusedJoinSize:    2,
		},
		calcDescription(9, 4, 10),
	)

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 8,
			Joins:             1,
			RemainderQuantity: 2,
			UnusedJoinSize:    2,
		},
		calcDescription(10, 4, 10),
	)

	for quantity := 11; quantity <= 16; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 8,
				EffectiveQuantity: 8 * (quantity / 8),
				Joins:             int(math.Ceil(float64(quantity) / 8)),
				RemainderQuantity: quantity % 8,
				UnusedJoinSize:    2,
			},
			calcDescription(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 16,
			Joins:             2,
			RemainderQuantity: 1,
			UnusedJoinSize:    2,
		},
		calcDescription(17, 4, 10),
	)

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 16,
			Joins:             2,
			RemainderQuantity: 2,
			UnusedJoinSize:    2,
		},
		calcDescription(18, 4, 10),
	)

	for quantity := 19; quantity <= 24; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 8,
				EffectiveQuantity: 8 * (quantity / 8),
				Joins:             int(math.Ceil(float64(quantity) / 8)),
				RemainderQuantity: quantity % 8,
				UnusedJoinSize:    2,
			},
			calcDescription(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 24,
			Joins:             3,
			RemainderQuantity: 1,
			UnusedJoinSize:    2,
		},
		calcDescription(25, 4, 10),
	)

	require.Equal(
		t,
		description{
			EffectiveJoinSize: 8,
			EffectiveQuantity: 24,
			Joins:             3,
			RemainderQuantity: 2,
			UnusedJoinSize:    2,
		},
		calcDescription(26, 4, 10),
	)
}

func TestCalcDescriptionBlockSize10(t *testing.T) {
	for quantity := 1; quantity <= 10; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 10,
				EffectiveQuantity: quantity,
				Joins:             1,
				RemainderQuantity: 0,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 10, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 11; quantity <= 30; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 10,
				EffectiveQuantity: 10 * (quantity / 10),
				Joins:             int(math.Ceil(float64(quantity) / 10)),
				RemainderQuantity: quantity % 10,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 10, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestCalcDescriptionBlockSize11(t *testing.T) {
	for quantity := 1; quantity <= 11; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 11,
				EffectiveQuantity: quantity,
				Joins:             1,
				RemainderQuantity: 0,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 11, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 12; quantity <= 40; quantity++ {
		require.Equal(
			t,
			description{
				EffectiveJoinSize: 11,
				EffectiveQuantity: 11 * (quantity / 11),
				Joins:             int(math.Ceil(float64(quantity) / 11)),
				RemainderQuantity: quantity % 11,
				UnusedJoinSize:    0,
			},
			calcDescription(quantity, 11, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestExpectedZeroes(t *testing.T) {
	require.Equal(t, [][]int{}, Expected(0, 0, 0))
	require.Equal(t, [][]int{}, Expected(0, 4, 0))
	require.Equal(t, [][]int{}, Expected(12, 0, 0))
	require.Equal(t, [][]int{}, Expected(12, 4, 0))

	require.Equal(t, [][]int{}, Expected(0, 0, 10))
	require.Equal(t, [][]int{}, Expected(0, 4, 10))
	require.Equal(t, [][]int{}, Expected(12, 0, 10))
}

func TestExpectedBlockSize1(t *testing.T) {
	for quantity := 1; quantity <= 10; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, quantity)},
			Expected(quantity, 1, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 11; quantity <= 20; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, 10), seq(11, quantity)},
			Expected(quantity, 1, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 21; quantity <= 30; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, 10), seq(11, 20), seq(21, quantity)},
			Expected(quantity, 1, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestExpectedBlockSize3(t *testing.T) {
	for quantity := 1; quantity <= 10; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, quantity)},
			Expected(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 11; quantity <= 19; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, 9), seq(10, quantity)},
			Expected(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 20; quantity <= 28; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, 9), seq(10, 18), seq(19, quantity)},
			Expected(quantity, 3, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestExpectedBlockSize4(t *testing.T) {
	for quantity := 1; quantity <= 10; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, quantity)},
			Expected(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 11; quantity <= 18; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, 8), seq(9, quantity)},
			Expected(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 19; quantity <= 26; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, 8), seq(9, 16), seq(17, quantity)},
			Expected(quantity, 4, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestExpectedBlockSize10(t *testing.T) {
	for quantity := 1; quantity <= 10; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, quantity)},
			Expected(quantity, 10, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 11; quantity <= 20; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, 10), seq(11, quantity)},
			Expected(quantity, 10, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 21; quantity <= 30; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, 10), seq(11, 20), seq(21, quantity)},
			Expected(quantity, 10, 10),
			"quantity: %v",
			quantity,
		)
	}
}

func TestExpectedBlockSize11(t *testing.T) {
	for quantity := 1; quantity <= 11; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, quantity)},
			Expected(quantity, 11, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 12; quantity <= 22; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, 11), seq(12, quantity)},
			Expected(quantity, 11, 10),
			"quantity: %v",
			quantity,
		)
	}

	for quantity := 23; quantity <= 33; quantity++ {
		require.Equal(
			t,
			[][]int{seq(1, 11), seq(12, 22), seq(23, quantity)},
			Expected(quantity, 11, 10),
			"quantity: %v",
			quantity,
		)
	}
}
