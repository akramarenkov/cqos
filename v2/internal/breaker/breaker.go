// Internal package used to break goroutine and wait it completion
package breaker

import "sync"

// Safely closes the channel
type Closing struct {
	closable chan struct{}
	closed   bool
	mutex    *sync.Mutex
}

func NewClosing() *Closing {
	cls := &Closing{
		closable: make(chan struct{}),
		mutex:    &sync.Mutex{},
	}

	return cls
}

func (cls *Closing) Close() {
	cls.mutex.Lock()
	defer cls.mutex.Unlock()

	if cls.closed {
		return
	}

	close(cls.closable)

	cls.closed = true
}

func (cls *Closing) Closed() <-chan struct{} {
	return cls.closable
}

// Used to break goroutine and wait it completion
type Breaker struct {
	breaker   *Closing
	completer *Closing
}

func New() *Breaker {
	brk := &Breaker{
		breaker:   NewClosing(),
		completer: NewClosing(),
	}

	return brk
}

func (brk *Breaker) Break() {
	brk.breaker.Close()
	<-brk.completer.Closed()
}

func (brk *Breaker) Breaked() <-chan struct{} {
	return brk.breaker.Closed()
}

func (brk *Breaker) Complete() {
	brk.completer.Close()
}
