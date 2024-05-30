// Internal package with implementation of the Stressor that was used to load the
// system and runtime.
package stressor

import (
	"runtime"
	"sync"

	"github.com/akramarenkov/breaker/breaker"
	"github.com/akramarenkov/cqos/v2/internal/launcher"
)

const (
	defaultCPUFactor     = 64
	defaultDataAmount    = 512
	defaultLaunchedAfter = 10
)

type Stressor struct {
	breaker  *breaker.Breaker
	launcher *launcher.Launcher

	cpuFactor int
	data      string
}

func New(cpuFactor int, dataAmount int) (*Stressor, error) {
	if cpuFactor == 0 {
		cpuFactor = defaultCPUFactor
	}

	if dataAmount == 0 {
		dataAmount = defaultDataAmount
	}

	data, err := getRandom(dataAmount)
	if err != nil {
		return nil, err
	}

	str := &Stressor{
		breaker:  breaker.New(),
		launcher: launcher.New(defaultLaunchedAfter),

		cpuFactor: cpuFactor,
		data:      string(data),
	}

	go str.main()

	str.launcher.Wait()

	return str, nil
}

func (str *Stressor) Stop() {
	str.breaker.Break()
}

func (str *Stressor) main() {
	defer str.breaker.Complete()

	str.loop()
}

func (str *Stressor) loop() {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for actors := 0; actors < str.cpuFactor*runtime.NumCPU(); actors++ {
		strings := make(chan string, 1)
		runes := make(chan []rune)

		strings <- str.data

		wg.Add(1)
		wg.Add(1)

		go str.runer(str.launcher.Add(), wg, strings, runes)
		go str.stringer(str.launcher.Add(), wg, runes, strings)
	}

	str.launcher.Launched()
}

func (str *Stressor) runer(
	id int,
	wg *sync.WaitGroup,
	input chan string,
	output chan []rune,
) {
	defer wg.Done()

	for {
		select {
		case <-str.breaker.IsBreaked():
			return
		case data := <-input:
			select {
			case <-str.breaker.IsBreaked():
				return
			case output <- []rune(data):
				str.launcher.Done(id)
			}
		}
	}
}

func (str *Stressor) stringer(
	id int,
	wg *sync.WaitGroup,
	input chan []rune,
	output chan string,
) {
	defer wg.Done()

	for {
		select {
		case <-str.breaker.IsBreaked():
			return
		case data := <-input:
			select {
			case <-str.breaker.IsBreaked():
				return
			case output <- string(data):
				str.launcher.Done(id)
			}
		}
	}
}
