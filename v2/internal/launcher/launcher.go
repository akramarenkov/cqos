// Internal package with implementation of the Launcher that used to wait for the launch
// of all goroutines after they have declared readiness n times.
package launcher

import (
	"sync"

	"github.com/akramarenkov/breaker/closing"
)

const (
	defaultDoneAfter = 1
)

type Launcher struct {
	added     []uint
	closing   *closing.Closing
	doneAfter uint
	mutex     *sync.RWMutex
	wg        *sync.WaitGroup
}

func New(doneAfter uint) *Launcher {
	if doneAfter == 0 {
		doneAfter = defaultDoneAfter
	}

	lnc := &Launcher{
		closing:   closing.New(),
		doneAfter: doneAfter,
		mutex:     &sync.RWMutex{},
		wg:        &sync.WaitGroup{},
	}

	return lnc
}

func (lnc *Launcher) add() int {
	lnc.mutex.Lock()
	defer lnc.mutex.Unlock()

	lnc.added = append(lnc.added, lnc.doneAfter)

	return len(lnc.added) - 1
}

func (lnc *Launcher) Add() int {
	lnc.wg.Add(1)
	return lnc.add()
}

func (lnc *Launcher) done(id int) bool {
	lnc.mutex.RLock()
	defer lnc.mutex.RUnlock()

	if lnc.added[id] == 0 {
		return false
	}

	lnc.added[id]--

	return lnc.added[id] == 0
}

func (lnc *Launcher) Done(id int) {
	if lnc.done(id) {
		lnc.wg.Done()
	}
}

func (lnc *Launcher) Launched() {
	lnc.wg.Wait()
	lnc.closing.Close()
}

func (lnc *Launcher) Wait() {
	<-lnc.closing.IsClosed()
}
