package limit

import "time"

const (
	// This value was chosen based on research into the results of the limit discipline
	// graphical tests.
	OptimizationInterval = 10 * time.Millisecond
)
