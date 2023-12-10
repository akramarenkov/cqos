package join

import (
	"time"

	"github.com/akramarenkov/cqos/v2/internal/consts"
)

// Maximum timeout error is calculated as timeout + timeout/divider.
//
// Relative timeout error in percent (inaccuracy) is calculated as 100/divider
func calcTickerDuration(timeout time.Duration, inaccuracy uint) (time.Duration, error) {
	if inaccuracy == 0 {
		return 0, ErrInvalidTimeoutInaccuracy
	}

	divider := consts.OneHundredPercent / inaccuracy

	if divider == 0 {
		return 0, ErrInvalidTimeoutInaccuracy
	}

	timeout /= time.Duration(divider)

	if timeout == 0 {
		return 0, ErrTimeoutTooSmall
	}

	return timeout, nil
}
