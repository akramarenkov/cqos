package durations

import (
	"sort"
	"time"
)

func createCopy(durations []time.Duration) []time.Duration {
	copied := make([]time.Duration, len(durations))

	copy(copied, durations)

	return copied
}

func Sort(durations []time.Duration) {
	less := func(i int, j int) bool {
		return durations[i] < durations[j]
	}

	sort.SliceStable(durations, less)
}

func IsSorted(durations []time.Duration) bool {
	less := func(i int, j int) bool {
		return durations[i] < durations[j]
	}

	return sort.SliceIsSorted(durations, less)
}

func CalcTotalDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	copied := createCopy(durations)

	Sort(copied)

	return copied[len(copied)-1]
}
