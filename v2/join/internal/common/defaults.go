package common

import "github.com/akramarenkov/cqos/v2/internal/consts"

const (
	// Default value of TimeoutInaccuracy option if it is not specified.
	DefaultTimeoutInaccuracy = 25

	// Minimum timeout, specifying which will not lead to an error when creating
	// disciplines with using the value from DefaultTimeoutInaccuracy as
	// TimeoutInaccuracy option.
	DefaultMinTimeout = (consts.OneHundredPercent *
		consts.ReliablyMeasurableDuration) / DefaultTimeoutInaccuracy

	// Default timeout used in tests.
	DefaultTestTimeout = 10 * DefaultMinTimeout
)
