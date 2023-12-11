// Internal package with quantity over time struct used in research
package qot

import (
	"time"
)

type QuantityOverTime struct {
	Quantity     uint
	RelativeTime time.Duration
}
