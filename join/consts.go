package join

import "github.com/akramarenkov/cqos/internal/consts"

const (
	defaultTimeoutInaccuracy = 25

	minDefaultTimeout = (consts.OneHundredPercent *
		consts.ReliablyMeasurableDuration) / defaultTimeoutInaccuracy
)
