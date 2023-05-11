package types

type Prioritized[Type any] struct {
	Item     Type
	Priority uint
}
