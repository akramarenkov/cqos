// Data types of prioritization discipline.
package types

// Describes the data distributed by the prioritization discipline.
type Prioritized[Type any] struct {
	Item     Type
	Priority uint
}
