package general

import (
	"golang.org/x/exp/constraints"
)

// Divides base value on to divider. Returns base value if divider is zero.
// Returns min value if result of dividing is less than min.
func DivideWithMin[Type constraints.Integer](base Type, divider Type, min Type) Type {
	if divider == 0 {
		return base
	}

	base /= divider

	if base < min {
		return min
	}

	return base
}
