package measurer

import "time"

type MeasureKind int

const (
	MeasureKindCompleted MeasureKind = iota + 1
	MeasureKindProcessed
	MeasureKindReceived
)

type Measure struct {
	Data         uint
	Kind         MeasureKind
	Priority     uint
	RelativeTime time.Duration
}
