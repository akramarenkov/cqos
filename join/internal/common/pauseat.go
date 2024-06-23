package common

import "time"

const (
	pausetAtRelativeDuration = 2.75
)

// Calculates pause duration value for timeout tests with which the tests will
// run reliably.
//
// For the reasons for choosing a value for the pausetAtRelativeDuration constant,
// see doc/test-pauseat-duration-select-problem.svg.
func CalcPauseAtDuration(timeout time.Duration) time.Duration {
	return time.Duration(pausetAtRelativeDuration * float64(timeout))
}
