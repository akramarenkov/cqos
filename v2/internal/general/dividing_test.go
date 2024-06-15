package general

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDivideWithMin(t *testing.T) {
	require.Equal(t, 10, DivideWithMin(10, 0, 3))
	require.Equal(t, 10, DivideWithMin(10, 1, 3))
	require.Equal(t, 5, DivideWithMin(10, 2, 3))
	require.Equal(t, 3, DivideWithMin(10, 3, 3))
	require.Equal(t, 3, DivideWithMin(10, 4, 3))
	require.Equal(t, 3, DivideWithMin(10, 5, 3))
	require.Equal(t, 3, DivideWithMin(10, 10, 3))
	require.Equal(t, 3, DivideWithMin(10, 11, 3))
	require.Equal(t, 3, DivideWithMin(10, 12, 3))
}
