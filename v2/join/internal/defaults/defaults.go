package defaults

import "github.com/akramarenkov/cqos/v2/internal/consts"

const (
	// Default value of TimeoutInaccuracy option if it is not specified.
	TimeoutInaccuracy = 25

	// Minimum timeout, specifying which will not lead to an error when creating
	// discipline with default value for TimeoutInaccuracy option.
	MinTimeout = (consts.HundredPercent * consts.ReliablyMeasurableDuration) / TimeoutInaccuracy

	// Default timeout used in tests.
	TestTimeout = 100 * MinTimeout
)
