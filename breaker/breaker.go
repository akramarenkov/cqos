// Mostly internal package used to break goroutine and wait it completion
package breaker

import "sync"

// Safely closes the channel
type Closing struct {
	channel chan struct{}
	once    *sync.Once
}

func NewClosing() *Closing {
	cls := &Closing{
		channel: make(chan struct{}),
		once:    &sync.Once{},
	}

	return cls
}

func (cls *Closing) Close() {
	do := func() {
		close(cls.channel)
	}

	cls.once.Do(do)
}

func (cls *Closing) Closed() <-chan struct{} {
	return cls.channel
}

// Used to break goroutine and wait it completion
type Breaker struct {
	completer   *Closing
	interrupter *Closing
}

func New() *Breaker {
	brk := &Breaker{
		completer:   NewClosing(),
		interrupter: NewClosing(),
	}

	return brk
}

func (brk *Breaker) Break() {
	brk.interrupter.Close()
	<-brk.completer.Closed()
}

func (brk *Breaker) Breaked() <-chan struct{} {
	return brk.interrupter.Closed()
}

func (brk *Breaker) Complete() {
	brk.completer.Close()
}
