package limit

import "time"

const (
	// the value was chosen based on studies of the graphical tests results and benchmarks
	minimumMeasuredDuration = 10 * time.Millisecond
)
