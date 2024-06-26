package inspect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInput(t *testing.T) {
	require.Equal(t, [][]int(nil), Input(0, 0))
	require.Equal(t, [][]int(nil), Input(1, 0))
	require.Equal(t, [][]int(nil), Input(0, 1))

	require.Equal(t, [][]int{{1}}, Input(1, 1))
	require.Equal(t, [][]int{{1}, {2}}, Input(2, 1))
	require.Equal(t, [][]int{{1}, {2}, {3}}, Input(3, 1))
	require.Equal(t, [][]int{{1}, {2}, {3}, {4}}, Input(4, 1))
	require.Equal(t, [][]int{{1}, {2}, {3}, {4}, {5}}, Input(5, 1))
	require.Equal(t, [][]int{{1}, {2}, {3}, {4}, {5}, {6}}, Input(6, 1))
	require.Equal(t, [][]int{{1}, {2}, {3}, {4}, {5}, {6}, {7}}, Input(7, 1))

	require.Equal(t, [][]int{{1}}, Input(1, 2))
	require.Equal(t, [][]int{{1, 2}}, Input(2, 2))
	require.Equal(t, [][]int{{1, 2}, {3}}, Input(3, 2))
	require.Equal(t, [][]int{{1, 2}, {3, 4}}, Input(4, 2))
	require.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, Input(5, 2))
	require.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}}, Input(6, 2))
	require.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}, {7}}, Input(7, 2))

	require.Equal(t, [][]int{{1}}, Input(1, 3))
	require.Equal(t, [][]int{{1, 2}}, Input(2, 3))
	require.Equal(t, [][]int{{1, 2, 3}}, Input(3, 3))
	require.Equal(t, [][]int{{1, 2, 3}, {4}}, Input(4, 3))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5}}, Input(5, 3))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}}, Input(6, 3))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7}}, Input(7, 3))
}
