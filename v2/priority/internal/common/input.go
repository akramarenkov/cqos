package common

type Input[Type any] struct {
	Channel <-chan Type
	Drained bool
}
