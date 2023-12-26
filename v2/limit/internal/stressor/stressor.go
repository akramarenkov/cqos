// Internal package with implementation of theStressor that was used to load the
// system and runtime
package stressor

import (
	"runtime"
	"sync"

	"github.com/akramarenkov/cqos/v2/internal/breaker"
	"github.com/akramarenkov/cqos/v2/limit/internal/starter"
)

const (
	defaultCPUFactor    = 64
	defaultDataAmount   = 512
	defaultStartedAfter = 10
)

type Stressor struct {
	breaker *breaker.Breaker
	starter *starter.Starter

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
		breaker: breaker.New(),
		starter: starter.New(defaultStartedAfter),

		cpuFactor: cpuFactor,
		data:      string(data),
	}

	go str.main()

	str.starter.Wait()

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

		go str.runer(str.starter.Add(), wg, strings, runes)
		go str.stringer(str.starter.Add(), wg, runes, strings)
	}

	str.starter.Started()
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
		case <-str.breaker.Breaked():
			return
		case data := <-input:
			select {
			case <-str.breaker.Breaked():
				return
			case output <- []rune(data):
				str.starter.Done(id)
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
		case <-str.breaker.Breaked():
			return
		case data := <-input:
			select {
			case <-str.breaker.Breaked():
				return
			case output <- string(data):
				str.starter.Done(id)
			}
		}
	}
}
