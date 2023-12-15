package stress

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

type Stress struct {
	breaker *breaker.Breaker
	starter *starter.Starter

	cpuFactor int
	data      string
}

func New(cpuFactor int, dataAmount int) (*Stress, error) {
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

	str := &Stress{
		breaker: breaker.New(),
		starter: starter.New(defaultStartedAfter),

		cpuFactor: cpuFactor,
		data:      string(data),
	}

	go str.main()

	str.starter.Wait()

	return str, nil
}

func (str *Stress) Stop() {
	str.breaker.Break()
}

func (str *Stress) main() {
	defer str.breaker.Complete()

	str.loop()
}

func (str *Stress) loop() {
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

func (str *Stress) runer(
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

func (str *Stress) stringer(
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
