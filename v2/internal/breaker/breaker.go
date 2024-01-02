// Internal package used to break goroutine and wait it completion
package breaker

import "github.com/akramarenkov/cqos/v2/internal/closing"

// Used to break goroutine and wait it completion
type Breaker struct {
	completer   *closing.Closing
	interrupter *closing.Closing
}

func New() *Breaker {
	brk := &Breaker{
		completer:   closing.New(),
		interrupter: closing.New(),
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
