package blocks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDivideSequence(t *testing.T) {
	require.Equal(t, [][]int(nil), DivideSequence(0, 0))
	require.Equal(t, [][]int(nil), DivideSequence(1, 0))
	require.Equal(t, [][]int(nil), DivideSequence(0, 1))

	require.Equal(t, [][]int{{1}}, DivideSequence(1, 1))
	require.Equal(t, [][]int{{1}, {2}}, DivideSequence(2, 1))
	require.Equal(t, [][]int{{1}, {2}, {3}}, DivideSequence(3, 1))
	require.Equal(t, [][]int{{1}, {2}, {3}, {4}}, DivideSequence(4, 1))
	require.Equal(t, [][]int{{1}, {2}, {3}, {4}, {5}}, DivideSequence(5, 1))
	require.Equal(t, [][]int{{1}, {2}, {3}, {4}, {5}, {6}}, DivideSequence(6, 1))
	require.Equal(t, [][]int{{1}, {2}, {3}, {4}, {5}, {6}, {7}}, DivideSequence(7, 1))

	require.Equal(t, [][]int{{1}}, DivideSequence(1, 2))
	require.Equal(t, [][]int{{1, 2}}, DivideSequence(2, 2))
	require.Equal(t, [][]int{{1, 2}, {3}}, DivideSequence(3, 2))
	require.Equal(t, [][]int{{1, 2}, {3, 4}}, DivideSequence(4, 2))
	require.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, DivideSequence(5, 2))
	require.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}}, DivideSequence(6, 2))
	require.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}, {7}}, DivideSequence(7, 2))

	require.Equal(t, [][]int{{1}}, DivideSequence(1, 3))
	require.Equal(t, [][]int{{1, 2}}, DivideSequence(2, 3))
	require.Equal(t, [][]int{{1, 2, 3}}, DivideSequence(3, 3))
	require.Equal(t, [][]int{{1, 2, 3}, {4}}, DivideSequence(4, 3))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5}}, DivideSequence(5, 3))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}}, DivideSequence(6, 3))
	require.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7}}, DivideSequence(7, 3))
}

func TestCalcExpectedJoinsZeroes(t *testing.T) {
	require.Equal(t, 0, CalcExpectedJoins(0, 0, 0))
	require.Equal(t, 0, CalcExpectedJoins(0, 4, 0))
	require.Equal(t, 0, CalcExpectedJoins(12, 0, 0))
	require.Equal(t, 0, CalcExpectedJoins(12, 4, 0))

	require.Equal(t, 0, CalcExpectedJoins(0, 0, 10))
	require.Equal(t, 0, CalcExpectedJoins(0, 4, 10))
	require.Equal(t, 0, CalcExpectedJoins(12, 0, 10))
}

func TestCalcExpectedJoinsBlockSize1(t *testing.T) {
	require.Equal(t, 1, CalcExpectedJoins(1, 1, 10))
	require.Equal(t, 1, CalcExpectedJoins(2, 1, 10))
	require.Equal(t, 1, CalcExpectedJoins(3, 1, 10))
	require.Equal(t, 1, CalcExpectedJoins(4, 1, 10))
	require.Equal(t, 1, CalcExpectedJoins(5, 1, 10))
	require.Equal(t, 1, CalcExpectedJoins(6, 1, 10))
	require.Equal(t, 1, CalcExpectedJoins(7, 1, 10))
	require.Equal(t, 1, CalcExpectedJoins(8, 1, 10))
	require.Equal(t, 1, CalcExpectedJoins(9, 1, 10))
	require.Equal(t, 1, CalcExpectedJoins(10, 1, 10))
	require.Equal(t, 2, CalcExpectedJoins(11, 1, 10))
	require.Equal(t, 2, CalcExpectedJoins(12, 1, 10))

	require.Equal(t, 10, CalcExpectedJoins(100, 1, 10))
	require.Equal(t, 11, CalcExpectedJoins(101, 1, 10))
	require.Equal(t, 11, CalcExpectedJoins(102, 1, 10))
	require.Equal(t, 11, CalcExpectedJoins(103, 1, 10))
	require.Equal(t, 11, CalcExpectedJoins(104, 1, 10))
	require.Equal(t, 11, CalcExpectedJoins(105, 1, 10))
	require.Equal(t, 11, CalcExpectedJoins(106, 1, 10))
	require.Equal(t, 11, CalcExpectedJoins(107, 1, 10))
	require.Equal(t, 11, CalcExpectedJoins(108, 1, 10))
	require.Equal(t, 11, CalcExpectedJoins(109, 1, 10))
	require.Equal(t, 11, CalcExpectedJoins(110, 1, 10))
	require.Equal(t, 12, CalcExpectedJoins(111, 1, 10))
	require.Equal(t, 12, CalcExpectedJoins(112, 1, 10))
	require.Equal(t, 12, CalcExpectedJoins(113, 1, 10))
	require.Equal(t, 12, CalcExpectedJoins(114, 1, 10))
	require.Equal(t, 12, CalcExpectedJoins(115, 1, 10))
	require.Equal(t, 12, CalcExpectedJoins(116, 1, 10))
	require.Equal(t, 12, CalcExpectedJoins(117, 1, 10))
	require.Equal(t, 12, CalcExpectedJoins(118, 1, 10))
	require.Equal(t, 12, CalcExpectedJoins(119, 1, 10))
	require.Equal(t, 12, CalcExpectedJoins(120, 1, 10))
	require.Equal(t, 13, CalcExpectedJoins(121, 1, 10))
	require.Equal(t, 13, CalcExpectedJoins(122, 1, 10))
	require.Equal(t, 13, CalcExpectedJoins(123, 1, 10))
	require.Equal(t, 13, CalcExpectedJoins(124, 1, 10))
	require.Equal(t, 13, CalcExpectedJoins(125, 1, 10))
}

func TestCalcExpectedJoinsBlockSize3(t *testing.T) {
	require.Equal(t, 1, CalcExpectedJoins(1, 3, 10))
	require.Equal(t, 1, CalcExpectedJoins(2, 3, 10))
	require.Equal(t, 1, CalcExpectedJoins(3, 3, 10))
	require.Equal(t, 1, CalcExpectedJoins(3, 3, 10))
	require.Equal(t, 1, CalcExpectedJoins(5, 3, 10))
	require.Equal(t, 1, CalcExpectedJoins(6, 3, 10))
	require.Equal(t, 1, CalcExpectedJoins(7, 3, 10))
	require.Equal(t, 1, CalcExpectedJoins(8, 3, 10))
	require.Equal(t, 1, CalcExpectedJoins(9, 3, 10))
	require.Equal(t, 1, CalcExpectedJoins(10, 3, 10))
	require.Equal(t, 2, CalcExpectedJoins(11, 3, 10))
	require.Equal(t, 2, CalcExpectedJoins(12, 3, 10))

	require.Equal(t, 11, CalcExpectedJoins(100, 3, 10))
	require.Equal(t, 12, CalcExpectedJoins(101, 3, 10))
	require.Equal(t, 12, CalcExpectedJoins(102, 3, 10))
	require.Equal(t, 12, CalcExpectedJoins(103, 3, 10))
	require.Equal(t, 12, CalcExpectedJoins(103, 3, 10))
	require.Equal(t, 12, CalcExpectedJoins(105, 3, 10))
	require.Equal(t, 12, CalcExpectedJoins(106, 3, 10))
	require.Equal(t, 12, CalcExpectedJoins(107, 3, 10))
	require.Equal(t, 12, CalcExpectedJoins(108, 3, 10))
	require.Equal(t, 12, CalcExpectedJoins(109, 3, 10))
	require.Equal(t, 13, CalcExpectedJoins(110, 3, 10))
	require.Equal(t, 13, CalcExpectedJoins(111, 3, 10))
	require.Equal(t, 13, CalcExpectedJoins(112, 3, 10))
	require.Equal(t, 13, CalcExpectedJoins(113, 3, 10))
	require.Equal(t, 13, CalcExpectedJoins(113, 3, 10))
	require.Equal(t, 13, CalcExpectedJoins(115, 3, 10))
	require.Equal(t, 13, CalcExpectedJoins(116, 3, 10))
	require.Equal(t, 13, CalcExpectedJoins(117, 3, 10))
	require.Equal(t, 13, CalcExpectedJoins(118, 3, 10))
	require.Equal(t, 14, CalcExpectedJoins(119, 3, 10))
	require.Equal(t, 14, CalcExpectedJoins(120, 3, 10))
	require.Equal(t, 14, CalcExpectedJoins(121, 3, 10))
	require.Equal(t, 14, CalcExpectedJoins(122, 3, 10))
	require.Equal(t, 14, CalcExpectedJoins(123, 3, 10))
	require.Equal(t, 14, CalcExpectedJoins(123, 3, 10))
	require.Equal(t, 14, CalcExpectedJoins(125, 3, 10))
}

func TestCalcExpectedJoinsBlockSize4(t *testing.T) {
	require.Equal(t, 1, CalcExpectedJoins(1, 4, 10))
	require.Equal(t, 1, CalcExpectedJoins(2, 4, 10))
	require.Equal(t, 1, CalcExpectedJoins(3, 4, 10))
	require.Equal(t, 1, CalcExpectedJoins(4, 4, 10))
	require.Equal(t, 1, CalcExpectedJoins(5, 4, 10))
	require.Equal(t, 1, CalcExpectedJoins(6, 4, 10))
	require.Equal(t, 1, CalcExpectedJoins(7, 4, 10))
	require.Equal(t, 1, CalcExpectedJoins(8, 4, 10))
	require.Equal(t, 1, CalcExpectedJoins(9, 4, 10))
	require.Equal(t, 1, CalcExpectedJoins(10, 4, 10))
	require.Equal(t, 2, CalcExpectedJoins(11, 4, 10))
	require.Equal(t, 2, CalcExpectedJoins(12, 4, 10))

	require.Equal(t, 13, CalcExpectedJoins(100, 4, 10))
	require.Equal(t, 13, CalcExpectedJoins(101, 4, 10))
	require.Equal(t, 13, CalcExpectedJoins(102, 4, 10))
	require.Equal(t, 13, CalcExpectedJoins(103, 4, 10))
	require.Equal(t, 13, CalcExpectedJoins(104, 4, 10))
	require.Equal(t, 13, CalcExpectedJoins(105, 4, 10))
	require.Equal(t, 13, CalcExpectedJoins(106, 4, 10))
	require.Equal(t, 14, CalcExpectedJoins(107, 4, 10))
	require.Equal(t, 14, CalcExpectedJoins(108, 4, 10))
	require.Equal(t, 14, CalcExpectedJoins(109, 4, 10))
	require.Equal(t, 14, CalcExpectedJoins(110, 4, 10))
	require.Equal(t, 14, CalcExpectedJoins(111, 4, 10))
	require.Equal(t, 14, CalcExpectedJoins(112, 4, 10))
	require.Equal(t, 14, CalcExpectedJoins(113, 4, 10))
	require.Equal(t, 14, CalcExpectedJoins(114, 4, 10))
	require.Equal(t, 15, CalcExpectedJoins(115, 4, 10))
	require.Equal(t, 15, CalcExpectedJoins(116, 4, 10))
	require.Equal(t, 15, CalcExpectedJoins(117, 4, 10))
	require.Equal(t, 15, CalcExpectedJoins(118, 4, 10))
	require.Equal(t, 15, CalcExpectedJoins(119, 4, 10))
	require.Equal(t, 15, CalcExpectedJoins(120, 4, 10))
	require.Equal(t, 15, CalcExpectedJoins(121, 4, 10))
	require.Equal(t, 15, CalcExpectedJoins(122, 4, 10))
	require.Equal(t, 16, CalcExpectedJoins(123, 4, 10))
	require.Equal(t, 16, CalcExpectedJoins(124, 4, 10))
	require.Equal(t, 16, CalcExpectedJoins(125, 4, 10))
}

func TestCalcExpectedJoinsBlockSize10(t *testing.T) {
	require.Equal(t, 1, CalcExpectedJoins(1, 10, 10))
	require.Equal(t, 1, CalcExpectedJoins(2, 10, 10))
	require.Equal(t, 1, CalcExpectedJoins(3, 10, 10))
	require.Equal(t, 1, CalcExpectedJoins(4, 10, 10))
	require.Equal(t, 1, CalcExpectedJoins(5, 10, 10))
	require.Equal(t, 1, CalcExpectedJoins(6, 10, 10))
	require.Equal(t, 1, CalcExpectedJoins(7, 10, 10))
	require.Equal(t, 1, CalcExpectedJoins(8, 10, 10))
	require.Equal(t, 1, CalcExpectedJoins(9, 10, 10))
	require.Equal(t, 1, CalcExpectedJoins(10, 10, 10))
	require.Equal(t, 2, CalcExpectedJoins(11, 10, 10))
	require.Equal(t, 2, CalcExpectedJoins(12, 10, 10))

	require.Equal(t, 10, CalcExpectedJoins(100, 10, 10))
	require.Equal(t, 11, CalcExpectedJoins(101, 10, 10))
	require.Equal(t, 11, CalcExpectedJoins(102, 10, 10))
	require.Equal(t, 11, CalcExpectedJoins(103, 10, 10))
	require.Equal(t, 11, CalcExpectedJoins(104, 10, 10))
	require.Equal(t, 11, CalcExpectedJoins(105, 10, 10))
	require.Equal(t, 11, CalcExpectedJoins(106, 10, 10))
	require.Equal(t, 11, CalcExpectedJoins(107, 10, 10))
	require.Equal(t, 11, CalcExpectedJoins(108, 10, 10))
	require.Equal(t, 11, CalcExpectedJoins(109, 10, 10))
	require.Equal(t, 11, CalcExpectedJoins(110, 10, 10))
	require.Equal(t, 12, CalcExpectedJoins(111, 10, 10))
	require.Equal(t, 12, CalcExpectedJoins(112, 10, 10))
	require.Equal(t, 12, CalcExpectedJoins(113, 10, 10))
	require.Equal(t, 12, CalcExpectedJoins(114, 10, 10))
	require.Equal(t, 12, CalcExpectedJoins(115, 10, 10))
	require.Equal(t, 12, CalcExpectedJoins(116, 10, 10))
	require.Equal(t, 12, CalcExpectedJoins(117, 10, 10))
	require.Equal(t, 12, CalcExpectedJoins(118, 10, 10))
	require.Equal(t, 12, CalcExpectedJoins(119, 10, 10))
	require.Equal(t, 12, CalcExpectedJoins(120, 10, 10))
	require.Equal(t, 13, CalcExpectedJoins(121, 10, 10))
	require.Equal(t, 13, CalcExpectedJoins(122, 10, 10))
	require.Equal(t, 13, CalcExpectedJoins(123, 10, 10))
	require.Equal(t, 13, CalcExpectedJoins(124, 10, 10))
	require.Equal(t, 13, CalcExpectedJoins(125, 10, 10))
}

func TestCalcExpectedJoinsBlockSize11(t *testing.T) {
	require.Equal(t, 1, CalcExpectedJoins(1, 11, 10))
	require.Equal(t, 1, CalcExpectedJoins(2, 11, 10))
	require.Equal(t, 1, CalcExpectedJoins(3, 11, 10))
	require.Equal(t, 1, CalcExpectedJoins(4, 11, 10))
	require.Equal(t, 1, CalcExpectedJoins(5, 11, 10))
	require.Equal(t, 1, CalcExpectedJoins(6, 11, 10))
	require.Equal(t, 1, CalcExpectedJoins(7, 11, 10))
	require.Equal(t, 1, CalcExpectedJoins(8, 11, 10))
	require.Equal(t, 1, CalcExpectedJoins(9, 11, 10))
	require.Equal(t, 1, CalcExpectedJoins(10, 11, 10))
	require.Equal(t, 1, CalcExpectedJoins(11, 11, 10))
	require.Equal(t, 2, CalcExpectedJoins(12, 11, 10))

	require.Equal(t, 10, CalcExpectedJoins(100, 11, 10))
	require.Equal(t, 10, CalcExpectedJoins(101, 11, 10))
	require.Equal(t, 10, CalcExpectedJoins(102, 11, 10))
	require.Equal(t, 10, CalcExpectedJoins(103, 11, 10))
	require.Equal(t, 10, CalcExpectedJoins(104, 11, 10))
	require.Equal(t, 10, CalcExpectedJoins(105, 11, 10))
	require.Equal(t, 10, CalcExpectedJoins(106, 11, 10))
	require.Equal(t, 10, CalcExpectedJoins(107, 11, 10))
	require.Equal(t, 10, CalcExpectedJoins(108, 11, 10))
	require.Equal(t, 10, CalcExpectedJoins(109, 11, 10))
	require.Equal(t, 10, CalcExpectedJoins(110, 11, 10))
	require.Equal(t, 11, CalcExpectedJoins(111, 11, 10))
	require.Equal(t, 11, CalcExpectedJoins(112, 11, 10))
	require.Equal(t, 11, CalcExpectedJoins(113, 11, 10))
	require.Equal(t, 11, CalcExpectedJoins(114, 11, 10))
	require.Equal(t, 11, CalcExpectedJoins(115, 11, 10))
	require.Equal(t, 11, CalcExpectedJoins(116, 11, 10))
	require.Equal(t, 11, CalcExpectedJoins(117, 11, 10))
	require.Equal(t, 11, CalcExpectedJoins(118, 11, 10))
	require.Equal(t, 11, CalcExpectedJoins(119, 11, 10))
	require.Equal(t, 11, CalcExpectedJoins(120, 11, 10))
	require.Equal(t, 11, CalcExpectedJoins(121, 11, 10))
	require.Equal(t, 12, CalcExpectedJoins(122, 11, 10))
	require.Equal(t, 12, CalcExpectedJoins(123, 11, 10))
	require.Equal(t, 12, CalcExpectedJoins(124, 11, 10))
	require.Equal(t, 12, CalcExpectedJoins(125, 11, 10))
}

func TestCalcExpectedJoinsWithTimeoutZeroes(t *testing.T) {
	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(0, 0, 0, 0))
	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(0, 5, 0, 0))
	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(12, 0, 0, 0))
	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(12, 5, 0, 0))

	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(0, 0, 4, 0))
	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(0, 5, 4, 0))
	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(12, 0, 4, 0))
	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(12, 5, 4, 0))

	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(0, 0, 0, 10))
	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(0, 5, 0, 10))
	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(12, 0, 0, 10))
	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(12, 5, 0, 10))

	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(0, 0, 4, 10))
	require.Equal(t, 0, CalcExpectedJoinsWithTimeout(0, 5, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 0, 4, 10))
}

func TestCalcExpectedJoinsWithTimeoutPauseAtAboveQuantity(t *testing.T) {
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(13, 13, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 14, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 15, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 16, 4, 10))
}

func TestCalcExpectedJoinsWithTimeoutBlockSize1(t *testing.T) {
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(1, 0, 1, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(1, 1, 1, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(2, 0, 1, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(2, 1, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(2, 2, 1, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(3, 0, 1, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(3, 1, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(3, 2, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(3, 3, 1, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(4, 0, 1, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(4, 1, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(4, 2, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(4, 3, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(4, 4, 1, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(5, 0, 1, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(5, 1, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(5, 2, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(5, 3, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(5, 4, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(5, 5, 1, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(6, 0, 1, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(6, 1, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(6, 2, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(6, 3, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(6, 4, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(6, 5, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(6, 6, 1, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(7, 0, 1, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(7, 1, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 2, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 3, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 4, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 5, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 6, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 7, 1, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(8, 0, 1, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(8, 1, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 2, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 3, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 4, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 5, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 6, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 7, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 8, 1, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(9, 0, 1, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(9, 1, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 2, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 3, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 4, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 5, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 6, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 7, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 8, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 9, 1, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(10, 0, 1, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(10, 1, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 2, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 3, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 4, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 5, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 6, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 7, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 8, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 9, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 10, 1, 10))

	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 0, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 1, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 2, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 3, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 4, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 5, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 6, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 7, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 8, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 9, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 10, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 11, 1, 10))

	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 0, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 1, 1, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(12, 2, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 3, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 4, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 5, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 6, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 7, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 8, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 9, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 10, 1, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 11, 1, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(12, 12, 1, 10))
}

func TestCalcExpectedJoinsWithTimeoutBlockSize3(t *testing.T) {
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(1, 0, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(1, 1, 3, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(2, 0, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(2, 1, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(2, 2, 3, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(3, 0, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(3, 1, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(3, 2, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(3, 3, 3, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(4, 0, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(4, 1, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(4, 2, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(4, 3, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(4, 4, 3, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(5, 0, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(5, 1, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(5, 2, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(5, 3, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(5, 4, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(5, 5, 3, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(6, 0, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(6, 1, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(6, 2, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(6, 3, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(6, 4, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(6, 5, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(6, 6, 3, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(7, 0, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(7, 1, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(7, 2, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(7, 3, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 4, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 5, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 6, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 7, 3, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(8, 0, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(8, 1, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(8, 2, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(8, 3, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 4, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 5, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 6, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 7, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 8, 3, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(9, 0, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(9, 1, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(9, 2, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(9, 3, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 4, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 5, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 6, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 7, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 8, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 9, 3, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(10, 0, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(10, 1, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(10, 2, 3, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(10, 3, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 4, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 5, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 6, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 7, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 8, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 9, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 10, 3, 10))

	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 0, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 1, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 2, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 3, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 4, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 5, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 6, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 7, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 8, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 9, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 10, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 11, 3, 10))

	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 0, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 1, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 2, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 3, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 4, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 5, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 6, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 7, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 8, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 9, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 10, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 11, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 12, 3, 10))

	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 0, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 1, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 2, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 3, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 4, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 5, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 6, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 7, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 8, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 9, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 10, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 11, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 12, 3, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(13, 13, 3, 10))

	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 0, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 1, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 2, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 3, 3, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(14, 4, 3, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(14, 5, 3, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(14, 6, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 7, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 8, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 9, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 10, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 11, 3, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 12, 3, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(14, 13, 3, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(14, 14, 3, 10))
}

func TestCalcExpectedJoinsWithTimeoutBlockSize4(t *testing.T) {
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(1, 0, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(1, 1, 4, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(2, 0, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(2, 1, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(2, 2, 4, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(3, 0, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(3, 1, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(3, 2, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(3, 3, 4, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(4, 0, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(4, 1, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(4, 2, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(4, 3, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(4, 4, 4, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(5, 0, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(5, 1, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(5, 2, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(5, 3, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(5, 4, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(5, 5, 4, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(6, 0, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(6, 1, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(6, 2, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(6, 3, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(6, 4, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(6, 5, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(6, 6, 4, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(7, 0, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(7, 1, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(7, 2, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(7, 3, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(7, 4, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 5, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 6, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(7, 7, 4, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(8, 0, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(8, 1, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(8, 2, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(8, 3, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(8, 4, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 5, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 6, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 7, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(8, 8, 4, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(9, 0, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(9, 1, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(9, 2, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(9, 3, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(9, 4, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 5, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 6, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 7, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 8, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(9, 9, 4, 10))

	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(10, 0, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(10, 1, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(10, 2, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(10, 3, 4, 10))
	require.Equal(t, 1, CalcExpectedJoinsWithTimeout(10, 4, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 5, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 6, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 7, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 8, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 9, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(10, 10, 4, 10))

	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 0, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 1, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 2, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 3, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 4, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 5, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 6, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 7, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 8, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 9, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 10, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(11, 11, 4, 10))

	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 0, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 1, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 2, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 3, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 4, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 5, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 6, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 7, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 8, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 9, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 10, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 11, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(12, 12, 4, 10))

	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 0, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 1, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 2, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 3, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 4, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 5, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 6, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 7, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 8, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 9, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 10, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 11, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(13, 12, 4, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(13, 13, 4, 10))

	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 0, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 1, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 2, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 3, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 4, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 5, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 6, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 7, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 8, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 9, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 10, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 11, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(14, 12, 4, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(14, 13, 4, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(14, 14, 4, 10))

	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(15, 0, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(15, 1, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(15, 2, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(15, 3, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(15, 4, 4, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(15, 5, 4, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(15, 6, 4, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(15, 7, 4, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(15, 8, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(15, 9, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(15, 10, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(15, 11, 4, 10))
	require.Equal(t, 2, CalcExpectedJoinsWithTimeout(15, 12, 4, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(15, 13, 4, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(15, 14, 4, 10))
	require.Equal(t, 3, CalcExpectedJoinsWithTimeout(15, 15, 4, 10))
}

func TestCalcExpectedJoinsWithTimeoutBlockSize10(t *testing.T) {
	for quantity := 1; quantity <= 100; quantity++ {
		for pauseAt := 1; pauseAt <= quantity; pauseAt++ {
			require.Equal(
				t,
				CalcExpectedJoins(quantity, 10, 10),
				CalcExpectedJoinsWithTimeout(quantity, pauseAt, 10, 10),
				"quantity: %v , pauseAt: %v",
				quantity,
				pauseAt,
			)
		}
	}
}

func TestCalcExpectedJoinsWithTimeoutBlockSize11(t *testing.T) {
	for quantity := 1; quantity <= 100; quantity++ {
		for pauseAt := 1; pauseAt <= quantity; pauseAt++ {
			require.Equal(
				t,
				CalcExpectedJoins(quantity, 11, 10),
				CalcExpectedJoinsWithTimeout(quantity, pauseAt, 11, 10),
				"quantity: %v , pauseAt: %v",
				quantity,
				pauseAt,
			)
		}
	}
}

func TestCalcExpectedJoinsWithTimeout5356(t *testing.T) {
	dataSet := []struct {
		quantity int
		expected int
	}{
		// 6 join before timeout, 1 join by timeout, 6 after timeout and so on
		{100, 6 + 1 + 6},
		// One item of last block fits in unused space of join equal to
		// JoinSize - BlockSize*(whole number of blocks in join)
		// In this case unused space of join equal to 2
		{101, 6 + 1 + 6},
		{102, 6 + 1 + 6}, // It also fits into unused space.

		{103, 6 + 1 + 7}, // It no longer fits into unused space
		{104, 6 + 1 + 7},
		{105, 6 + 1 + 7},
		{106, 6 + 1 + 7},
		{107, 6 + 1 + 7},
		{108, 6 + 1 + 7},
		{109, 6 + 1 + 7}, // It also fits into unused space.
		{110, 6 + 1 + 7}, // It also fits into unused space.

		{111, 6 + 1 + 8}, // It no longer fits into unused space
		{112, 6 + 1 + 8},
		{113, 6 + 1 + 8},
		{114, 6 + 1 + 8},
		{115, 6 + 1 + 8},
		{116, 6 + 1 + 8},
		{117, 6 + 1 + 8}, // It also fits into unused space.
		{118, 6 + 1 + 8}, // It also fits into unused space.

		{119, 6 + 1 + 9}, // It no longer fits into unused space
		{120, 6 + 1 + 9},
		{121, 6 + 1 + 9},
		{122, 6 + 1 + 9},
		{123, 6 + 1 + 9},
		{124, 6 + 1 + 9},
		{125, 6 + 1 + 9}, // It also fits into unused space.
	}

	for pauseAt := 53; pauseAt <= 56; pauseAt++ {
		for _, item := range dataSet {
			require.Equal(
				t,
				item.expected,
				CalcExpectedJoinsWithTimeout(item.quantity, pauseAt, 4, 10),
				"quantity: %v , pauseAt: %v",
				item.quantity,
				pauseAt,
			)
		}
	}
}

func TestCalcExpectedJoinsWithTimeout5760(t *testing.T) {
	// Pausing at these positions will not cause an additional output slice to
	// appear and therefore the values obtained using CalcExpectedJoins and
	// CalcExpectedJoinsWithTimeout must match
	for pauseAt := 57; pauseAt <= 60; pauseAt++ {
		for quantity := 100; quantity <= 125; quantity++ {
			require.Equal(
				t,
				CalcExpectedJoins(quantity, 4, 10),
				CalcExpectedJoinsWithTimeout(quantity, pauseAt, 4, 10),
				"quantity: %v , pauseAt: %v",
				quantity,
				pauseAt,
			)
		}
	}
}

func TestCalcExpectedJoinsWithTimeout6164(t *testing.T) {
	dataSet := []struct {
		quantity int
		expected int
	}{
		// 7 join before timeout, 1 join by timeout, 5 after timeout and so on
		{100, 7 + 1 + 5},
		// One item of last block fits in unused space of join equal to
		// JoinSize - BlockSize*(whole number of blocks in join)
		// In this case unused space of join equal to 2
		{101, 7 + 1 + 5},
		{102, 7 + 1 + 5}, // It also fits into unused space.

		{103, 7 + 1 + 6}, // It no longer fits into unused space
		{104, 7 + 1 + 6},
		{105, 7 + 1 + 6},
		{106, 7 + 1 + 6},
		{107, 7 + 1 + 6},
		{108, 7 + 1 + 6},
		{109, 7 + 1 + 6}, // It also fits into unused space.
		{110, 7 + 1 + 6}, // It also fits into unused space.

		{111, 7 + 1 + 7}, // It no longer fits into unused space
		{112, 7 + 1 + 7},
		{113, 7 + 1 + 7},
		{114, 7 + 1 + 7},
		{115, 7 + 1 + 7},
		{116, 7 + 1 + 7},
		{117, 7 + 1 + 7}, // It also fits into unused space.
		{118, 7 + 1 + 7}, // It also fits into unused space.

		{119, 7 + 1 + 8}, // It no longer fits into unused space
		{120, 7 + 1 + 8},
		{121, 7 + 1 + 8},
		{122, 7 + 1 + 8},
		{123, 7 + 1 + 8},
		{124, 7 + 1 + 8},
		{125, 7 + 1 + 8}, // It also fits into unused space.
	}

	for pauseAt := 61; pauseAt <= 64; pauseAt++ {
		for _, item := range dataSet {
			require.Equal(
				t,
				item.expected,
				CalcExpectedJoinsWithTimeout(item.quantity, pauseAt, 4, 10),
				"quantity: %v , pauseAt: %v",
				item.quantity,
				pauseAt,
			)
		}
	}
}

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
