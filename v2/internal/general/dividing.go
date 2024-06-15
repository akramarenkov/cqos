package general

import (
	"golang.org/x/exp/constraints"
)

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
