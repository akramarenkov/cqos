package spinner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpinner(t *testing.T) {
	id := New(0, 2)

	require.Equal(t, 0, id.Actual())

	id.Spin()
	require.Equal(t, 1, id.Actual())

	id.Spin()
	require.Equal(t, 2, id.Actual())

	id.Spin()
	require.Equal(t, 0, id.Actual())

	id.Spin()
	require.Equal(t, 1, id.Actual())

	id.Spin()
	require.Equal(t, 2, id.Actual())

	id.Spin()
	require.Equal(t, 0, id.Actual())
}

func TestSpinnerEndIsEqualBegin(t *testing.T) {
	id := New(0, 0)

	require.Equal(t, 0, id.Actual())

	id.Spin()
	require.Equal(t, 0, id.Actual())

	id.Spin()
	require.Equal(t, 0, id.Actual())

	id.Spin()
	require.Equal(t, 0, id.Actual())

	id = New(1, 1)

	require.Equal(t, 1, id.Actual())

	id.Spin()
	require.Equal(t, 1, id.Actual())

	id.Spin()
	require.Equal(t, 1, id.Actual())

	id.Spin()
	require.Equal(t, 1, id.Actual())
}

func TestSpinnerEndIsLessBegin(t *testing.T) {
	id := New(0, -1)

	require.Equal(t, 0, id.Actual())

	id.Spin()
	require.Equal(t, 0, id.Actual())

	id.Spin()
	require.Equal(t, 0, id.Actual())

	id.Spin()
	require.Equal(t, 0, id.Actual())

	id = New(1, 0)

	require.Equal(t, 1, id.Actual())

	id.Spin()
	require.Equal(t, 1, id.Actual())

	id.Spin()
	require.Equal(t, 1, id.Actual())

	id.Spin()
	require.Equal(t, 1, id.Actual())
}

func BenchmarkSpinner(b *testing.B) {
	id := New(0, 10)

	b.ResetTimer()

	for run := 0; run < b.N; run++ {
		_ = id.Actual()
		id.Spin()
	}
}
