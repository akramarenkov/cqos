package join

import (
	"errors"
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
)

var (
	ErrInvalidTimeoutInaccuracy = errors.New("invalid timeout inaccuracy")
	ErrTimeoutTooSmall          = errors.New("timeout value is too small")
)

// Maximum timeout error is calculated as timeout + timeout/divider.
//
// Relative timeout error in percent (inaccuracy) is calculated as 100/divider
func calcInterruptInterval(
	timeout time.Duration,
	inaccuracy uint,
) (time.Duration, error) {
	if inaccuracy == 0 {
		return 0, ErrInvalidTimeoutInaccuracy
	}

	divider := consts.OneHundredPercent / inaccuracy

	if divider == 0 {
		return 0, ErrInvalidTimeoutInaccuracy
	}

	timeout /= time.Duration(divider)

	if timeout < consts.ReliablyMeasurableDuration {
		return 0, ErrTimeoutTooSmall
	}

	return timeout, nil
}

func calcInterruptIntervalZeroAllowed(
	timeout time.Duration,
	inaccuracy uint,
) (time.Duration, error) {
	if timeout <= 0 {
		return 0, nil
	}

	return calcInterruptInterval(timeout, inaccuracy)
}
