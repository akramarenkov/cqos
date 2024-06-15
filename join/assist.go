package join

import (
	"errors"
	"time"

	"github.com/akramarenkov/cqos/internal/consts"
)

var (
	ErrTimeoutInaccuracyTooBig = errors.New("timeout inaccuracy is too big")
	ErrTimeoutInaccuracyZero   = errors.New("timeout inaccuracy is zero")
	ErrTimeoutTooSmall         = errors.New("timeout value is too small")
)

// Maximum timeout error is calculated as timeout + timeout/divider.
//
// Relative timeout error in percent (inaccuracy) is calculated as 100/divider.
func calcInterruptInterval(
	timeout time.Duration,
	inaccuracy uint,
) (time.Duration, error) {
	if inaccuracy == 0 {
		return 0, ErrTimeoutInaccuracyZero
	}

	divider := consts.HundredPercent / inaccuracy

	if divider == 0 {
		return 0, ErrTimeoutInaccuracyTooBig
	}

	interval := timeout / time.Duration(divider)

	if interval < consts.ReliablyMeasurableDuration {
		return 0, ErrTimeoutTooSmall
	}

	return interval, nil
}

func calcInterruptIntervalNonPositiveAllowed(
	timeout time.Duration,
	inaccuracy uint,
) (time.Duration, error) {
	if timeout <= 0 {
		return 0, nil
	}

	return calcInterruptInterval(timeout, inaccuracy)
}
