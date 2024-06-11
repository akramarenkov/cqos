package common

const (
	// Number 2 indicates the number of goroutines involved in processing
	// slices with double buffering. There are 2 of them - accumulation and
	// send goroutines. The accumulation goroutine pass slices to the send goroutine
	// through an interim channel.
	InvolvedInProcessing = 2

	// Defines buffers quantity needed for buffering. Cannot be less than number
	// of goroutines involved in processing slices, but there may be more than
	// this value, although this does not make sense.
	BuffersQuantity = InvolvedInProcessing

	// To prevent from using still unsent slices by the accumulation goroutine,
	// it must be blocked from writing to the interim channel on the last one of
	// the unsent slices.
	InterimCapacity = BuffersQuantity - InvolvedInProcessing
)
