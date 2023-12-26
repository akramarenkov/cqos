// Internal package with common constants
package consts

import "time"

const (
	OneHundredPercent = 100
)

const (
	EnableGraphsEnv = "CQOS_ENABLE_GRAPHS"
)

const (
	// The value was chosen based on studies of the graphical tests results and
	// benchmarks for the limit discipline
	ReliablyMeasurableDuration = 10 * time.Millisecond
)
