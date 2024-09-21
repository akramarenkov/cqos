// Internal default constants for all join packages.
package defs

import (
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
	"github.com/akramarenkov/cqos/v2/join/defaults"
)

const (
	// Minimum timeout, specifying which will not lead to an error when creating
	// discipline with default value for TimeoutInaccuracy option.
	MinTimeout = consts.HundredPercent / defaults.TimeoutInaccuracy

	// Default timeout used in tests.
	TestTimeout = 10 * time.Second
)
