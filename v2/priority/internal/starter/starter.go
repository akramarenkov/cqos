// Internal package with implementation of the Starter which is used to run multiple
// goroutines at the same time
package starter

import (
	"sync"
	"time"
)

type Starter struct {
	StartedAt time.Time

	trigger chan struct{}
	wg      *sync.WaitGroup
}

func New() *Starter {
	str := &Starter{
		trigger: make(chan struct{}),
		wg:      &sync.WaitGroup{},
	}

	return str
}

func (str *Starter) Ready(delta int) {
	str.wg.Add(delta)
}

func (str *Starter) Set() {
	str.wg.Done()

	<-str.trigger
}

func (str *Starter) Go() {
	str.wg.Wait()

	str.StartedAt = time.Now()

	close(str.trigger)
}
