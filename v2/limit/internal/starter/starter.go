// Internal package with implementation of the Starter that used to wait for the start
// of all goroutines after they have declared readiness n times.
package starter

import (
	"sync"

	"github.com/akramarenkov/breaker/closing"
)

const (
	defaultDoneAfter = 1
)

type Starter struct {
	added     []uint
	closing   *closing.Closing
	doneAfter uint
	mutex     *sync.RWMutex
	wg        *sync.WaitGroup
}

func New(doneAfter uint) *Starter {
	if doneAfter == 0 {
		doneAfter = defaultDoneAfter
	}

	str := &Starter{
		closing:   closing.New(),
		doneAfter: doneAfter,
		mutex:     &sync.RWMutex{},
		wg:        &sync.WaitGroup{},
	}

	return str
}

func (str *Starter) add() int {
	str.mutex.Lock()
	defer str.mutex.Unlock()

	str.added = append(str.added, str.doneAfter)

	return len(str.added) - 1
}

func (str *Starter) Add() int {
	str.wg.Add(1)
	return str.add()
}

func (str *Starter) done(id int) bool {
	str.mutex.RLock()
	defer str.mutex.RUnlock()

	if str.added[id] == 0 {
		return false
	}

	str.added[id]--

	return str.added[id] == 0
}

func (str *Starter) Done(id int) {
	if str.done(id) {
		str.wg.Done()
	}
}

func (str *Starter) Started() {
	str.wg.Wait()
	str.closing.Close()
}

func (str *Starter) Wait() {
	<-str.closing.IsClosed()
}
