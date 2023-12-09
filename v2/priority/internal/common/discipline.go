package common

import "github.com/akramarenkov/cqos/v2/priority/types"

type Discipline[Type any] interface {
	Output() <-chan types.Prioritized[Type]
	Release(priority uint)
	Err() <-chan error
}
