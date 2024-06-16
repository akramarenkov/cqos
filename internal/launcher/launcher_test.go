package launcher

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIdle(*testing.T) {
	launcher := New(0)

	launcher.Created()
	launcher.Wait()
}

func TestLaunchedAtOne(*testing.T) {
	launcher := New(0)

	id := launcher.Add()

	launcher.Done(id)

	launcher.Created()
	launcher.Wait()
}

func TestUnlaunchedAtOne(t *testing.T) {
	launcher := New(0)

	_ = launcher.Add()

	isLaunched := make(chan bool)

	go func() {
		launcher.Created()
		launcher.Wait()

		close(isLaunched)
	}()

	// We don't wait for the located above goroutine to start for simplicity
	// We hope that during the timeout all actions to launch the located above goroutine
	// will be completed and the Created method will be called

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-isLaunched:
		require.FailNow(t, "must not be launched")
	case <-timer.C:
	}
}

func TestLaunchedAtTwo(*testing.T) {
	launcher := New(2)

	id := launcher.Add()

	launcher.Done(id)
	launcher.Done(id)

	launcher.Created()
	launcher.Wait()
}

func TestUnlaunchedAtTwo(t *testing.T) {
	launcher := New(2)

	id := launcher.Add()

	launcher.Done(id)

	isLaunched := make(chan bool)

	go func() {
		launcher.Created()
		launcher.Wait()

		close(isLaunched)
	}()

	// We don't wait for the located above goroutine to start for simplicity
	// We hope that during the timeout all actions to launch the located above goroutine
	// will be completed and the Created method will be called

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-isLaunched:
		require.FailNow(t, "must not be launched")
	case <-timer.C:
	}
}

func TestDoneExcess(*testing.T) {
	launcher := New(2)

	id := launcher.Add()

	launcher.Done(id)
	launcher.Done(id)
	launcher.Done(id)
	launcher.Done(id)

	launcher.Created()
	launcher.Wait()
}

func TestLauncher(*testing.T) {
	launchedAt := uint(10)
	goroutinesQuantity := 20
	iterationsNumber := 2 * launchedAt

	launcher := New(launchedAt)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	go func() {
		for range goroutinesQuantity {
			id := launcher.Add()

			wg.Add(1)

			go func() {
				defer wg.Done()

				for range iterationsNumber {
					launcher.Done(id)
				}
			}()
		}

		launcher.Created()
	}()

	launcher.Wait()
}

func TestUnlaunched(t *testing.T) {
	launchedAt := uint(10)
	goroutinesQuantity := 20
	iterationsNumber := 2 * launchedAt

	launcher := New(launchedAt)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	go func() {
		for range goroutinesQuantity {
			id := launcher.Add()

			wg.Add(1)

			go func() {
				defer wg.Done()

				for range iterationsNumber {
					if id != goroutinesQuantity/2 {
						launcher.Done(id)
					}
				}
			}()
		}

		launcher.Created()
	}()

	isLaunched := make(chan bool)

	go func() {
		launcher.Wait()

		close(isLaunched)
	}()

	// We don't wait for the located above goroutine to start for simplicity
	// We hope that during the timeout all actions to launch the located above goroutine
	// will be completed and the Created method will be called

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-isLaunched:
		require.FailNow(t, "must not be launched")
	case <-timer.C:
	}
}
