// Internal package used to closes the channel once.
package closing

import "sync"

// Closes the channel, but only once.
type Closing struct {
	channel chan struct{}
	once    *sync.Once
}

func New() *Closing {
	cls := &Closing{
		channel: make(chan struct{}),
		once:    &sync.Once{},
	}

	return cls
}

func (cls *Closing) Close() {
	cls.once.Do(cls.close)
}

func (cls *Closing) close() {
	close(cls.channel)
}

func (cls *Closing) Closed() <-chan struct{} {
	return cls.channel
}
