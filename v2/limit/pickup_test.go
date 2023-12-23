package limit

import (
	"testing"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/limit/internal/stress"

	"github.com/stretchr/testify/require"
)

func TestPickUp(t *testing.T) {
	testPickUp(t, false)
	testPickUp(t, true)
}

func testPickUp(t *testing.T, stressSystem bool) {
	if stressSystem {
		stress, err := stress.New(0, 0)
		require.NoError(t, err)

		defer stress.Stop()
	}

	t.Log(pickUpMinimumDuration(1e6, 1))
}

func pickUpMinimumDuration(addQuantity int, maxDiff int) time.Duration {
	const min = 1000

	for quantity := min; quantity <= addQuantity+min; quantity++ {
		if picked := pickUpMinimumDurationOne(quantity, maxDiff); picked != 0 {
			return picked
		}
	}

	return 0
}

func pickUpMinimumDurationOne(quantity int, maxDiff int) time.Duration {
	durations := make([]time.Duration, quantity)

	startedAt := time.Now()

	for id := 0; id < quantity; id++ {
		durations[id] = extrapolateDuration(startedAt, quantity, id)
	}

	expected := durations[len(durations)-1]

	for _, duration := range durations {
		diff := int(((duration - expected) * consts.OneHundredPercent) / expected)

		if diff < 0 {
			continue
		}

		if diff > maxDiff {
			continue
		}

		return duration
	}

	return 0
}

func extrapolateDuration(startedAt time.Time, quantity int, id int) time.Duration {
	return time.Duration(float64(time.Since(startedAt)) * float64(quantity) / float64(id+1))
}
