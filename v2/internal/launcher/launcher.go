// Internal package with implementation of the Launcher that is used to wait for
// the launch of several goroutines. They will be considered launched if they report
// this several times each.
package launcher

import (
	"sync"

	"github.com/akramarenkov/breaker/closing"
)

const (
	defaultLaunchedAt = 1
)

// Launcher is used to wait for the launch of several goroutines. They will be
// considered launched if they report this launchedAt times each.
type Launcher struct {
	launchedAt uint

	awaited []uint
	mutex   *sync.RWMutex

	closing *closing.Closing
	wg      *sync.WaitGroup
}

// Creates Launcher instance. If the launchedAt value is zero, the default value
// of 1 will be used.
func New(launchedAt uint) *Launcher {
	if launchedAt == 0 {
		launchedAt = defaultLaunchedAt
	}

	lnc := &Launcher{
		launchedAt: launchedAt,

		mutex: &sync.RWMutex{},

		closing: closing.New(),
		wg:      &sync.WaitGroup{},
	}

	return lnc
}

func (lnc *Launcher) add() int {
	lnc.mutex.Lock()
	defer lnc.mutex.Unlock()

	lnc.awaited = append(lnc.awaited, lnc.launchedAt)

	return len(lnc.awaited) - 1
}

// Adds a goroutine to the list of those whose launch is await. Returns the
// goroutine ID that must be used in the Done method.
func (lnc *Launcher) Add() int {
	lnc.wg.Add(1)
	return lnc.add()
}

func (lnc *Launcher) done(id int) bool {
	lnc.mutex.RLock()
	defer lnc.mutex.RUnlock()

	if lnc.awaited[id] == 0 {
		return false
	}

	lnc.awaited[id]--

	return lnc.awaited[id] == 0
}

// Marks the goroutine with the ID as launched. For a goroutine to be marked as finally
// launched it must call the Done method as many times as specified when creating the
// Launcher instance.
//
// This method with the same ID must not be run in parallel.
func (lnc *Launcher) Done(id int) {
	if lnc.done(id) {
		lnc.wg.Done()
	}
}

// Indicates to the Launcher that the goroutines creation process is complete. This
// method must be called after the goroutines creation process is completed.
func (lnc *Launcher) Created() {
	// Because it is impossible to call the WaitGroup.Wait method in parallel with the
	// WaitGroup.Add methods, then an additional channel is used to notify the "main"
	// goroutine that "child" goroutines have finished launched
	lnc.wg.Wait()
	lnc.closing.Close()
}

// Waits for goroutines to be launched.
func (lnc *Launcher) Wait() {
	<-lnc.closing.IsClosed()
}
