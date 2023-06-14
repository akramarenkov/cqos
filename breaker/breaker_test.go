package breaker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClosing(t *testing.T) {
	closing := NewClosing()

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

	closing.Close()

	select {
	case <-closing.Closed():
	default:
		require.FailNow(t, "must be closed")
	}
}

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
