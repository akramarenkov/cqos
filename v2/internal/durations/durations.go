// Internal package used to work with slice of durations
package durations

import (
	"slices"
	"time"
)

func CalcTotalDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	copied := slices.Clone(durations)

	slices.Sort(copied)

	return copied[len(copied)-1]
}
