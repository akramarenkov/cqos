package inspect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeq(t *testing.T) {
	require.Panics(t, func() { seq(0, -2) })
	require.Panics(t, func() { seq(2, 0) })
	require.Panics(t, func() { seq(4, 1) })
	require.Panics(t, func() { seq(6, 2) })

	require.Equal(t, []int{}, seq(0, -1))
	require.Equal(t, []int{}, seq(1, 0))
	require.Equal(t, []int{}, seq(2, 1))
	require.Equal(t, []int{}, seq(3, 2))

	require.Equal(t, []int{0}, seq(0, 0))
	require.Equal(t, []int{0, 1}, seq(0, 1))
	require.Equal(t, []int{0, 1, 2}, seq(0, 2))
	require.Equal(t, []int{0, 1, 2, 3}, seq(0, 3))
	require.Equal(t, []int{0, 1, 2, 3, 4}, seq(0, 4))

	require.Equal(t, []int{7}, seq(7, 7))
	require.Equal(t, []int{7, 8}, seq(7, 8))
	require.Equal(t, []int{7, 8, 9}, seq(7, 9))
	require.Equal(t, []int{7, 8, 9, 10}, seq(7, 10))
	require.Equal(t, []int{7, 8, 9, 10, 11}, seq(7, 11))

	require.Equal(t, []int{8}, seq(8, 8))
	require.Equal(t, []int{8, 9}, seq(8, 9))
	require.Equal(t, []int{8, 9, 10}, seq(8, 10))
	require.Equal(t, []int{8, 9, 10, 11}, seq(8, 11))
	require.Equal(t, []int{8, 9, 10, 11, 12}, seq(8, 12))
}
