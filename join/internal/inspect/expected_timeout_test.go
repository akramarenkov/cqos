package inspect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpectedWithTimeoutZeroes(t *testing.T) {
	require.Equal(t, [][]int{}, ExpectedWithTimeout(0, 0, 0, 0))
	require.Equal(t, [][]int{}, ExpectedWithTimeout(0, 5, 0, 0))
	require.Equal(t, [][]int{}, ExpectedWithTimeout(12, 0, 0, 0))
	require.Equal(t, [][]int{}, ExpectedWithTimeout(12, 5, 0, 0))

	require.Equal(t, [][]int{}, ExpectedWithTimeout(0, 0, 4, 0))
	require.Equal(t, [][]int{}, ExpectedWithTimeout(0, 5, 4, 0))
	require.Equal(t, [][]int{}, ExpectedWithTimeout(12, 0, 4, 0))
	require.Equal(t, [][]int{}, ExpectedWithTimeout(12, 5, 4, 0))

	require.Equal(t, [][]int{}, ExpectedWithTimeout(0, 0, 0, 10))
	require.Equal(t, [][]int{}, ExpectedWithTimeout(0, 5, 0, 10))
	require.Equal(t, [][]int{}, ExpectedWithTimeout(12, 0, 0, 10))
	require.Equal(t, [][]int{}, ExpectedWithTimeout(12, 5, 0, 10))

	require.Equal(t, [][]int{}, ExpectedWithTimeout(0, 0, 4, 10))
	require.Equal(t, [][]int{}, ExpectedWithTimeout(0, 5, 4, 10))

	require.Equal(
		t,
		[][]int{
			{1, 2, 3, 4, 5, 6, 7, 8},
			{9, 10, 11, 12},
		},
		ExpectedWithTimeout(12, 0, 4, 10),
	)
}

func TestExpectedWithTimeoutBlockSize1(t *testing.T) {
	require.Equal(t, [][]int{{1}}, ExpectedWithTimeout(1, 1, 1, 10))

	require.Equal(t, [][]int{{1, 2}}, ExpectedWithTimeout(2, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2}}, ExpectedWithTimeout(2, 2, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3}}, ExpectedWithTimeout(3, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3}}, ExpectedWithTimeout(3, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3}}, ExpectedWithTimeout(3, 3, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4}}, ExpectedWithTimeout(4, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3, 4}}, ExpectedWithTimeout(4, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3, 4}}, ExpectedWithTimeout(4, 3, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4}}, ExpectedWithTimeout(4, 4, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5}}, ExpectedWithTimeout(5, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3, 4, 5}}, ExpectedWithTimeout(5, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3, 4, 5}}, ExpectedWithTimeout(5, 3, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5}}, ExpectedWithTimeout(5, 4, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5}}, ExpectedWithTimeout(5, 5, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}}, ExpectedWithTimeout(6, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3, 4, 5, 6}}, ExpectedWithTimeout(6, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3, 4, 5, 6}}, ExpectedWithTimeout(6, 3, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}}, ExpectedWithTimeout(6, 4, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6}}, ExpectedWithTimeout(6, 5, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}, {6}}, ExpectedWithTimeout(6, 6, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}}, ExpectedWithTimeout(7, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3, 4, 5, 6, 7}}, ExpectedWithTimeout(7, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3, 4, 5, 6, 7}}, ExpectedWithTimeout(7, 3, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7}}, ExpectedWithTimeout(7, 4, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7}}, ExpectedWithTimeout(7, 5, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}, {6, 7}}, ExpectedWithTimeout(7, 6, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7}}, ExpectedWithTimeout(7, 7, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3, 4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3, 4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 3, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 4, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}},
		ExpectedWithTimeout(8, 5, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}, {6, 7, 8}},
		ExpectedWithTimeout(8, 6, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8}},
		ExpectedWithTimeout(8, 7, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}, {8}},
		ExpectedWithTimeout(8, 8, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3, 4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3, 4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 3, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 4, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 5, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}, {6, 7, 8, 9}},
		ExpectedWithTimeout(9, 6, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9}},
		ExpectedWithTimeout(9, 7, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}, {8, 9}},
		ExpectedWithTimeout(9, 8, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9}},
		ExpectedWithTimeout(9, 9, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3, 4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3, 4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 3, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 4, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 5, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 6, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10}},
		ExpectedWithTimeout(10, 7, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10}},
		ExpectedWithTimeout(10, 8, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10}},
		ExpectedWithTimeout(10, 9, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10}},
		ExpectedWithTimeout(10, 10, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11}},
		ExpectedWithTimeout(11, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3, 4, 5, 6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3, 4, 5, 6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 3, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 4, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 5, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 6, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 7, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10, 11}},
		ExpectedWithTimeout(11, 8, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11}},
		ExpectedWithTimeout(11, 9, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11}},
		ExpectedWithTimeout(11, 10, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11}},
		ExpectedWithTimeout(11, 11, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11, 12}},
		ExpectedWithTimeout(12, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, {12}},
		ExpectedWithTimeout(12, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 3, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 4, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 5, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 6, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 7, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 8, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}},
		ExpectedWithTimeout(12, 9, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12}},
		ExpectedWithTimeout(12, 10, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11, 12}},
		ExpectedWithTimeout(12, 11, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11}, {12}},
		ExpectedWithTimeout(12, 12, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11, 12, 13}},
		ExpectedWithTimeout(13, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, {12, 13}},
		ExpectedWithTimeout(13, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, {13}},
		ExpectedWithTimeout(13, 3, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 4, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 5, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 6, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 7, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 8, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 9, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}},
		ExpectedWithTimeout(13, 10, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11, 12, 13}},
		ExpectedWithTimeout(13, 11, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11}, {12, 13}},
		ExpectedWithTimeout(13, 12, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11, 12}, {13}},
		ExpectedWithTimeout(13, 13, 1, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11, 12, 13, 14}},
		ExpectedWithTimeout(14, 1, 1, 10))
	require.Equal(t, [][]int{{1}, {2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, {12, 13, 14}},
		ExpectedWithTimeout(14, 2, 1, 10))
	require.Equal(t, [][]int{{1, 2}, {3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, {13, 14}},
		ExpectedWithTimeout(14, 3, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, {14}},
		ExpectedWithTimeout(14, 4, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 5, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 6, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 7, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 8, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 9, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 10, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11, 12, 13, 14}},
		ExpectedWithTimeout(14, 11, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11}, {12, 13, 14}},
		ExpectedWithTimeout(14, 12, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11, 12}, {13, 14}},
		ExpectedWithTimeout(14, 13, 1, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, {11, 12, 13}, {14}},
		ExpectedWithTimeout(14, 14, 1, 10))
}

func TestExpectedWithTimeoutBlockSize3(t *testing.T) {
	require.Equal(t, [][]int{{1}}, ExpectedWithTimeout(1, 1, 3, 10))

	require.Equal(t, [][]int{{1, 2}}, ExpectedWithTimeout(2, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2}}, ExpectedWithTimeout(2, 2, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3}}, ExpectedWithTimeout(3, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}}, ExpectedWithTimeout(3, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}}, ExpectedWithTimeout(3, 3, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4}}, ExpectedWithTimeout(4, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}}, ExpectedWithTimeout(4, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}}, ExpectedWithTimeout(4, 3, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4}}, ExpectedWithTimeout(4, 4, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5}}, ExpectedWithTimeout(5, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}}, ExpectedWithTimeout(5, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}}, ExpectedWithTimeout(5, 3, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5}}, ExpectedWithTimeout(5, 4, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5}}, ExpectedWithTimeout(5, 5, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}}, ExpectedWithTimeout(6, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}}, ExpectedWithTimeout(6, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}}, ExpectedWithTimeout(6, 3, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}}, ExpectedWithTimeout(6, 4, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}}, ExpectedWithTimeout(6, 5, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}}, ExpectedWithTimeout(6, 6, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}}, ExpectedWithTimeout(7, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}}, ExpectedWithTimeout(7, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}}, ExpectedWithTimeout(7, 3, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7}}, ExpectedWithTimeout(7, 4, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7}}, ExpectedWithTimeout(7, 5, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7}}, ExpectedWithTimeout(7, 6, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7}}, ExpectedWithTimeout(7, 7, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 3, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 4, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 5, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 6, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8}},
		ExpectedWithTimeout(8, 7, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8}},
		ExpectedWithTimeout(8, 8, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 3, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 4, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 5, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 6, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9}},
		ExpectedWithTimeout(9, 7, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9}},
		ExpectedWithTimeout(9, 8, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9}},
		ExpectedWithTimeout(9, 9, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 3, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 4, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 5, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 6, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10}},
		ExpectedWithTimeout(10, 7, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10}},
		ExpectedWithTimeout(10, 8, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10}},
		ExpectedWithTimeout(10, 9, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10}},
		ExpectedWithTimeout(10, 10, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11}},
		ExpectedWithTimeout(11, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11}},
		ExpectedWithTimeout(11, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11}},
		ExpectedWithTimeout(11, 3, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 4, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 5, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 6, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 7, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 8, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 9, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11}},
		ExpectedWithTimeout(11, 10, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11}},
		ExpectedWithTimeout(11, 11, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12}},
		ExpectedWithTimeout(12, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12}},
		ExpectedWithTimeout(12, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12}},
		ExpectedWithTimeout(12, 3, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 4, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 5, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 6, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 7, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 8, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 9, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12}},
		ExpectedWithTimeout(12, 10, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12}},
		ExpectedWithTimeout(12, 11, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12}},
		ExpectedWithTimeout(12, 12, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}},
		ExpectedWithTimeout(13, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}},
		ExpectedWithTimeout(13, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}},
		ExpectedWithTimeout(13, 3, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 4, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 5, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 6, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 7, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 8, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 9, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}},
		ExpectedWithTimeout(13, 10, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}},
		ExpectedWithTimeout(13, 11, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}},
		ExpectedWithTimeout(13, 12, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12}, {13}},
		ExpectedWithTimeout(13, 13, 3, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 1, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 2, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 3, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12}, {13, 14}},
		ExpectedWithTimeout(14, 4, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12}, {13, 14}},
		ExpectedWithTimeout(14, 5, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6, 7, 8, 9, 10, 11, 12}, {13, 14}},
		ExpectedWithTimeout(14, 6, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 7, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 8, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 9, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 10, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 11, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 12, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12}, {13, 14}},
		ExpectedWithTimeout(14, 13, 3, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}, {10, 11, 12}, {13, 14}},
		ExpectedWithTimeout(14, 14, 3, 10))
}

func TestExpectedWithTimeoutBlockSize4(t *testing.T) { //nolint:maintidx
	require.Equal(t, [][]int{{1}}, ExpectedWithTimeout(1, 1, 4, 10))

	require.Equal(t, [][]int{{1, 2}}, ExpectedWithTimeout(2, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2}}, ExpectedWithTimeout(2, 2, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3}}, ExpectedWithTimeout(3, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3}}, ExpectedWithTimeout(3, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3}}, ExpectedWithTimeout(3, 3, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4}}, ExpectedWithTimeout(4, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}}, ExpectedWithTimeout(4, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}}, ExpectedWithTimeout(4, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}}, ExpectedWithTimeout(4, 4, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5}}, ExpectedWithTimeout(5, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}}, ExpectedWithTimeout(5, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}}, ExpectedWithTimeout(5, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5}}, ExpectedWithTimeout(5, 4, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5}}, ExpectedWithTimeout(5, 5, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}}, ExpectedWithTimeout(6, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}}, ExpectedWithTimeout(6, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}}, ExpectedWithTimeout(6, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}}, ExpectedWithTimeout(6, 4, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6}}, ExpectedWithTimeout(6, 5, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6}}, ExpectedWithTimeout(6, 6, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}}, ExpectedWithTimeout(7, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}}, ExpectedWithTimeout(7, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}}, ExpectedWithTimeout(7, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7}}, ExpectedWithTimeout(7, 4, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7}}, ExpectedWithTimeout(7, 5, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7}}, ExpectedWithTimeout(7, 6, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7}}, ExpectedWithTimeout(7, 7, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}},
		ExpectedWithTimeout(8, 4, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}},
		ExpectedWithTimeout(8, 5, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}},
		ExpectedWithTimeout(8, 6, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}},
		ExpectedWithTimeout(8, 7, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}},
		ExpectedWithTimeout(8, 8, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 4, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 5, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 6, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 7, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9}},
		ExpectedWithTimeout(9, 8, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9}},
		ExpectedWithTimeout(9, 9, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 4, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 5, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 6, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 7, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10}},
		ExpectedWithTimeout(10, 8, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10}},
		ExpectedWithTimeout(10, 9, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10}},
		ExpectedWithTimeout(10, 10, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11}},
		ExpectedWithTimeout(11, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11}},
		ExpectedWithTimeout(11, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11}},
		ExpectedWithTimeout(11, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11}},
		ExpectedWithTimeout(11, 4, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 5, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 6, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 7, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11}},
		ExpectedWithTimeout(11, 8, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11}},
		ExpectedWithTimeout(11, 9, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11}},
		ExpectedWithTimeout(11, 10, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11}},
		ExpectedWithTimeout(11, 11, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}},
		ExpectedWithTimeout(12, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}},
		ExpectedWithTimeout(12, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}},
		ExpectedWithTimeout(12, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}},
		ExpectedWithTimeout(12, 4, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 5, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 6, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 7, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12}},
		ExpectedWithTimeout(12, 8, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}},
		ExpectedWithTimeout(12, 9, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}},
		ExpectedWithTimeout(12, 10, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}},
		ExpectedWithTimeout(12, 11, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}},
		ExpectedWithTimeout(12, 12, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 4, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 5, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 6, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 7, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 8, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 9, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 10, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 11, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13}},
		ExpectedWithTimeout(13, 12, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}, {13}},
		ExpectedWithTimeout(13, 13, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 4, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 5, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 6, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 7, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 8, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 9, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 10, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 11, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14}},
		ExpectedWithTimeout(14, 12, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14}},
		ExpectedWithTimeout(14, 13, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14}},
		ExpectedWithTimeout(14, 14, 4, 10))

	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15}},
		ExpectedWithTimeout(15, 1, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15}},
		ExpectedWithTimeout(15, 2, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15}},
		ExpectedWithTimeout(15, 3, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15}},
		ExpectedWithTimeout(15, 4, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12}, {13, 14, 15}},
		ExpectedWithTimeout(15, 5, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12}, {13, 14, 15}},
		ExpectedWithTimeout(15, 6, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12}, {13, 14, 15}},
		ExpectedWithTimeout(15, 7, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10, 11, 12}, {13, 14, 15}},
		ExpectedWithTimeout(15, 8, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15}},
		ExpectedWithTimeout(15, 9, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15}},
		ExpectedWithTimeout(15, 10, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15}},
		ExpectedWithTimeout(15, 11, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15}},
		ExpectedWithTimeout(15, 12, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15}},
		ExpectedWithTimeout(15, 13, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15}},
		ExpectedWithTimeout(15, 14, 4, 10))
	require.Equal(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15}},
		ExpectedWithTimeout(15, 15, 4, 10))
}

func TestExpectedWithTimeoutBlockSize10(t *testing.T) {
	for quantity := 1; quantity <= 10; quantity++ {
		for pauseAt := 1; pauseAt <= quantity; pauseAt++ {
			require.Equal(
				t,
				[][]int{
					seq(1, quantity),
				},
				ExpectedWithTimeout(quantity, pauseAt, 10, 10),
			)
		}
	}

	for quantity := 11; quantity <= 20; quantity++ {
		for pauseAt := 1; pauseAt <= quantity; pauseAt++ {
			require.Equal(
				t,
				[][]int{
					seq(1, 10),
					seq(11, quantity),
				},
				ExpectedWithTimeout(quantity, pauseAt, 10, 10),
			)
		}
	}

	for quantity := 21; quantity <= 30; quantity++ {
		for pauseAt := 1; pauseAt <= quantity; pauseAt++ {
			require.Equal(
				t,
				[][]int{
					seq(1, 10),
					seq(11, 20),
					seq(21, quantity),
				},
				ExpectedWithTimeout(quantity, pauseAt, 10, 10),
			)
		}
	}
}

func TestExpectedWithTimeoutBlockSize11(t *testing.T) {
	for quantity := 1; quantity <= 11; quantity++ {
		for pauseAt := 1; pauseAt <= quantity; pauseAt++ {
			require.Equal(
				t,
				[][]int{
					seq(1, quantity),
				},
				ExpectedWithTimeout(quantity, pauseAt, 11, 10),
			)
		}
	}

	for quantity := 12; quantity <= 22; quantity++ {
		for pauseAt := 1; pauseAt <= quantity; pauseAt++ {
			require.Equal(
				t,
				[][]int{
					seq(1, 11),
					seq(12, quantity),
				},
				ExpectedWithTimeout(quantity, pauseAt, 11, 10),
			)
		}
	}

	for quantity := 23; quantity <= 33; quantity++ {
		for pauseAt := 1; pauseAt <= quantity; pauseAt++ {
			require.Equal(
				t,
				[][]int{
					seq(1, 11),
					seq(12, 22),
					seq(23, quantity),
				},
				ExpectedWithTimeout(quantity, pauseAt, 11, 10),
			)
		}
	}
}

func BenchmarkExpectedWithTimeout(b *testing.B) {
	for range b.N {
		_ = ExpectedWithTimeout(b.N, 13, 4, 10)
	}
}
