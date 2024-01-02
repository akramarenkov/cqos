package closing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClosing(t *testing.T) {
	closing := New()

	select {
	case <-closing.Closed():
		require.FailNow(t, "must not be closed")
	default:
	}

	closing.Close()

	select {
	case <-closing.Closed():
	default:
		require.FailNow(t, "must be closed")
	}

	require.NotPanics(t, func() { closing.Close() })

	select {
	case <-closing.Closed():
	default:
		require.FailNow(t, "must be closed")
	}
}

func BenchmarkRace(b *testing.B) {
	closing := New()

	for run := 0; run < b.N; run++ {
		b.RunParallel(
			func(pb *testing.PB) {
				for pb.Next() {
					closing.Close()
				}
			},
		)
	}
}
