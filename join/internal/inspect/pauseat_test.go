package inspect

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPickUpPauseAtZeroes(t *testing.T) {
	require.Equal(t, 0, PickUpPauseAt(0, 0, 0, 0))
	require.Equal(t, 0, PickUpPauseAt(0, 5, 0, 0))
	require.Equal(t, 0, PickUpPauseAt(12, 0, 0, 0))
	require.Equal(t, 0, PickUpPauseAt(12, 5, 0, 0))

	require.Equal(t, 0, PickUpPauseAt(0, 0, 4, 0))
	require.Equal(t, 0, PickUpPauseAt(0, 5, 4, 0))
	require.Equal(t, 0, PickUpPauseAt(12, 0, 4, 0))
	require.Equal(t, 0, PickUpPauseAt(12, 5, 4, 0))

	require.Equal(t, 0, PickUpPauseAt(0, 0, 0, 10))
	require.Equal(t, 0, PickUpPauseAt(0, 5, 0, 10))
	require.Equal(t, 0, PickUpPauseAt(12, 0, 0, 10))
	require.Equal(t, 0, PickUpPauseAt(12, 5, 0, 10))

	require.Equal(t, 0, PickUpPauseAt(0, 0, 4, 10))
	require.Equal(t, 0, PickUpPauseAt(0, 5, 4, 10))
	require.Equal(t, 0, PickUpPauseAt(12, 0, 4, 10))
}

func TestPickUpPauseAtEstimatedAboveQuantity(t *testing.T) {
	require.Equal(t, 13, PickUpPauseAt(13, 13, 4, 10))
	require.Equal(t, 0, PickUpPauseAt(13, 14, 4, 10))
	require.Equal(t, 0, PickUpPauseAt(13, 15, 4, 10))
	require.Equal(t, 0, PickUpPauseAt(13, 16, 4, 10))
}

func TestPickUpPauseAtPickedAboveQuantity(t *testing.T) {
	require.Equal(t, 13, PickUpPauseAt(13, 11, 4, 10))
	require.Equal(t, 0, PickUpPauseAt(12, 11, 4, 10))
	require.Equal(t, 0, PickUpPauseAt(11, 11, 4, 10))
}

func TestPickUpPauseAtBlockSize1(t *testing.T) {
	require.Equal(t, 2, PickUpPauseAt(100, 1, 1, 10))
	require.Equal(t, 2, PickUpPauseAt(100, 2, 1, 10))
	require.Equal(t, 3, PickUpPauseAt(100, 3, 1, 10))
	require.Equal(t, 4, PickUpPauseAt(100, 4, 1, 10))
	require.Equal(t, 5, PickUpPauseAt(100, 5, 1, 10))
	require.Equal(t, 6, PickUpPauseAt(100, 6, 1, 10))
	require.Equal(t, 7, PickUpPauseAt(100, 7, 1, 10))
	require.Equal(t, 8, PickUpPauseAt(100, 8, 1, 10))
	require.Equal(t, 9, PickUpPauseAt(100, 9, 1, 10))
	require.Equal(t, 10, PickUpPauseAt(100, 10, 1, 10))
	require.Equal(t, 12, PickUpPauseAt(100, 11, 1, 10))
	require.Equal(t, 12, PickUpPauseAt(100, 12, 1, 10))

	require.Equal(t, 47, PickUpPauseAt(100, 47, 1, 10))
	require.Equal(t, 48, PickUpPauseAt(100, 48, 1, 10))
	require.Equal(t, 49, PickUpPauseAt(100, 49, 1, 10))
	require.Equal(t, 50, PickUpPauseAt(100, 50, 1, 10))
	require.Equal(t, 52, PickUpPauseAt(100, 51, 1, 10))
	require.Equal(t, 52, PickUpPauseAt(100, 52, 1, 10))
	require.Equal(t, 53, PickUpPauseAt(100, 53, 1, 10))
	require.Equal(t, 54, PickUpPauseAt(100, 54, 1, 10))
	require.Equal(t, 55, PickUpPauseAt(100, 55, 1, 10))
	require.Equal(t, 56, PickUpPauseAt(100, 56, 1, 10))
	require.Equal(t, 57, PickUpPauseAt(100, 57, 1, 10))
	require.Equal(t, 58, PickUpPauseAt(100, 58, 1, 10))
	require.Equal(t, 59, PickUpPauseAt(100, 59, 1, 10))
	require.Equal(t, 60, PickUpPauseAt(100, 60, 1, 10))
	require.Equal(t, 62, PickUpPauseAt(100, 61, 1, 10))
	require.Equal(t, 62, PickUpPauseAt(100, 62, 1, 10))
}

func TestPickUpPauseAtBlockSize3(t *testing.T) {
	require.Equal(t, 4, PickUpPauseAt(100, 1, 3, 10))
	require.Equal(t, 4, PickUpPauseAt(100, 2, 3, 10))
	require.Equal(t, 4, PickUpPauseAt(100, 3, 3, 10))
	require.Equal(t, 4, PickUpPauseAt(100, 4, 3, 10))
	require.Equal(t, 5, PickUpPauseAt(100, 5, 3, 10))
	require.Equal(t, 6, PickUpPauseAt(100, 6, 3, 10))
	require.Equal(t, 7, PickUpPauseAt(100, 7, 3, 10))
	require.Equal(t, 8, PickUpPauseAt(100, 8, 3, 10))
	require.Equal(t, 9, PickUpPauseAt(100, 9, 3, 10))
	require.Equal(t, 13, PickUpPauseAt(100, 10, 3, 10))
	require.Equal(t, 13, PickUpPauseAt(100, 11, 3, 10))
	require.Equal(t, 13, PickUpPauseAt(100, 12, 3, 10))
	require.Equal(t, 13, PickUpPauseAt(100, 13, 3, 10))
	require.Equal(t, 14, PickUpPauseAt(100, 14, 3, 10))

	require.Equal(t, 49, PickUpPauseAt(100, 47, 3, 10))
	require.Equal(t, 49, PickUpPauseAt(100, 48, 3, 10))
	require.Equal(t, 49, PickUpPauseAt(100, 49, 3, 10))
	require.Equal(t, 50, PickUpPauseAt(100, 50, 3, 10))
	require.Equal(t, 51, PickUpPauseAt(100, 51, 3, 10))
	require.Equal(t, 52, PickUpPauseAt(100, 52, 3, 10))
	require.Equal(t, 53, PickUpPauseAt(100, 53, 3, 10))
	require.Equal(t, 54, PickUpPauseAt(100, 54, 3, 10))
	require.Equal(t, 58, PickUpPauseAt(100, 55, 3, 10))
	require.Equal(t, 58, PickUpPauseAt(100, 56, 3, 10))
	require.Equal(t, 58, PickUpPauseAt(100, 57, 3, 10))
}

func TestPickUpPauseAtBlockSize4(t *testing.T) {
	require.Equal(t, 5, PickUpPauseAt(100, 1, 4, 10))
	require.Equal(t, 5, PickUpPauseAt(100, 2, 4, 10))
	require.Equal(t, 5, PickUpPauseAt(100, 3, 4, 10))
	require.Equal(t, 5, PickUpPauseAt(100, 4, 4, 10))
	require.Equal(t, 5, PickUpPauseAt(100, 5, 4, 10))
	require.Equal(t, 6, PickUpPauseAt(100, 6, 4, 10))
	require.Equal(t, 7, PickUpPauseAt(100, 7, 4, 10))
	require.Equal(t, 8, PickUpPauseAt(100, 8, 4, 10))
	require.Equal(t, 13, PickUpPauseAt(100, 9, 4, 10))
	require.Equal(t, 13, PickUpPauseAt(100, 10, 4, 10))
	require.Equal(t, 13, PickUpPauseAt(100, 11, 4, 10))
	require.Equal(t, 13, PickUpPauseAt(100, 12, 4, 10))
	require.Equal(t, 13, PickUpPauseAt(100, 13, 4, 10))
	require.Equal(t, 14, PickUpPauseAt(100, 14, 4, 10))

	require.Equal(t, 47, PickUpPauseAt(100, 47, 4, 10))
	require.Equal(t, 48, PickUpPauseAt(100, 48, 4, 10))
	require.Equal(t, 53, PickUpPauseAt(100, 49, 4, 10))
	require.Equal(t, 53, PickUpPauseAt(100, 50, 4, 10))
	require.Equal(t, 53, PickUpPauseAt(100, 51, 4, 10))
	require.Equal(t, 53, PickUpPauseAt(100, 52, 4, 10))
	require.Equal(t, 53, PickUpPauseAt(100, 53, 4, 10))
	require.Equal(t, 54, PickUpPauseAt(100, 54, 4, 10))
	require.Equal(t, 55, PickUpPauseAt(100, 55, 4, 10))
	require.Equal(t, 56, PickUpPauseAt(100, 56, 4, 10))
	require.Equal(t, 61, PickUpPauseAt(100, 57, 4, 10))
}

func TestPickUpPauseAtBlockSize10(t *testing.T) {
	require.Equal(t, 0, PickUpPauseAt(100, 1, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 2, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 3, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 4, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 5, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 6, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 7, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 8, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 9, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 10, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 11, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 12, 10, 10))

	require.Equal(t, 0, PickUpPauseAt(100, 47, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 48, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 49, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 50, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 51, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 52, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 53, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 54, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 55, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 56, 10, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 57, 10, 10))
}

func TestPickUpPauseAtBlockSize11(t *testing.T) {
	require.Equal(t, 0, PickUpPauseAt(100, 1, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 2, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 3, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 4, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 5, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 6, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 7, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 8, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 9, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 10, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 11, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 12, 11, 10))

	require.Equal(t, 0, PickUpPauseAt(100, 47, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 48, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 49, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 50, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 51, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 52, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 53, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 54, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 55, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 56, 11, 10))
	require.Equal(t, 0, PickUpPauseAt(100, 57, 11, 10))
}

func TestCalcPauseAtDuration(t *testing.T) {
	require.Equal(t, 275*time.Millisecond, CalcPauseAtDuration(100*time.Millisecond))
}
