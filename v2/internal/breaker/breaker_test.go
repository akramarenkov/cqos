package breaker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBreaker(t *testing.T) {
	breaker := New()

	select {
	case <-breaker.Breaked():
		require.FailNow(t, "must not be breaked")
	default:
	}

	go func() {
		defer breaker.Complete()
		<-breaker.Breaked()
	}()

	breaker.Break()

	select {
	case <-breaker.Breaked():
	default:
		require.FailNow(t, "must be breaked")
	}

	breaker.Break()

	select {
	case <-breaker.Breaked():
	default:
		require.FailNow(t, "must be breaked")
	}
}
