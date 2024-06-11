package launcher

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIdle(*testing.T) {
	launcher := New(0)

	launcher.Launched()

	launcher.Wait()
}

func TestDoneAfterOne(*testing.T) {
	launcher := New(0)

	id := launcher.Add()

	launcher.Done(id)

	launcher.Launched()

	launcher.Wait()
}

func TestUndoneAfterOne(t *testing.T) {
	launcher := New(0)

	_ = launcher.Add()

	done := make(chan bool)

	go func() {
		launcher.Launched()
		launcher.Wait()

		close(done)
	}()

	// We don't wait for the goroutine to start for simplicity
	// We hope that during the timeout all actions will be completed

	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()

	select {
	case <-done:
		require.FailNow(t, "must not be launched")
	case <-timeout.C:
	}
}

func TestDoneAfterTwo(*testing.T) {
	launcher := New(2)

	id := launcher.Add()

	launcher.Done(id)
	launcher.Done(id)

	launcher.Launched()

	launcher.Wait()
}

func TestUndoneAfterTwo(t *testing.T) {
	launcher := New(2)

	id := launcher.Add()

	launcher.Done(id)

	done := make(chan bool)

	go func() {
		launcher.Launched()
		launcher.Wait()

		close(done)
	}()

	// We don't wait for the goroutine to start for simplicity
	// We hope that during the timeout all actions will be completed

	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()

	select {
	case <-done:
		require.FailNow(t, "must not be launched")
	case <-timeout.C:
	}
}

func TestDoneExcess(*testing.T) {
	launcher := New(2)

	id := launcher.Add()

	launcher.Done(id)
	launcher.Done(id)
	launcher.Done(id)

	launcher.Launched()

	launcher.Wait()
}

func TestGeneral(*testing.T) {
	doneAfter := uint(10)

	launcher := New(doneAfter)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	handlersQuantity := 20

	go func() {
		for range handlersQuantity {
			id := launcher.Add()

			wg.Add(1)

			go func() {
				defer wg.Done()

				for range 2 * doneAfter {
					launcher.Done(id)
				}
			}()
		}

		launcher.Launched()
	}()

	launcher.Wait()
}

func TestGeneralUndone(t *testing.T) {
	doneAfter := uint(10)

	launcher := New(doneAfter)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	handlersQuantity := 20

	go func() {
		for range handlersQuantity {
			id := launcher.Add()

			wg.Add(1)

			go func() {
				defer wg.Done()

				for range 2 * doneAfter {
					if id != handlersQuantity/2 {
						launcher.Done(id)
					}
				}
			}()
		}

		launcher.Launched()
	}()

	done := make(chan bool)

	go func() {
		launcher.Wait()

		close(done)
	}()

	// We don't wait for the goroutines to start for simplicity
	// We hope that during the timeout all actions will be completed

	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()

	select {
	case <-done:
		require.FailNow(t, "must not be launched")
	case <-timeout.C:
	}
}
