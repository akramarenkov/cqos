package qot

import (
	"time"
)

// Quantity over time.
type QOT struct {
	Quantity     uint
	RelativeTime time.Duration
}
