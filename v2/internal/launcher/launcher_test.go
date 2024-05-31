package launcher

import (
	"sync"
	"testing"
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

func TestDoneAfterTwo(*testing.T) {
	launcher := New(2)

	id := launcher.Add()

	launcher.Done(id)
	launcher.Done(id)

	launcher.Launched()

	launcher.Wait()
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
	launcher := New(10)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	go func() {
		for range 100 {
			id := launcher.Add()

			wg.Add(1)

			go func() {
				defer wg.Done()

				for range 1000 {
					launcher.Done(id)
				}
			}()
		}

		launcher.Launched()
	}()

	launcher.Wait()
}
