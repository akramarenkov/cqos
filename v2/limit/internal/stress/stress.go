package stress

import (
	"runtime"
	"sync"

	"github.com/akramarenkov/cqos/v2/internal/breaker"
)

const (
	defaultDataAmount = 10000
)

type Stress struct {
	breaker *breaker.Breaker

	data string
}

func New(dataAmount uint) (*Stress, error) {
	if dataAmount == 0 {
		dataAmount = defaultDataAmount
	}

	data, err := getRandom(dataAmount)
	if err != nil {
		return nil, err
	}

	str := &Stress{
		breaker: breaker.New(),

		data: string(data),
	}

	go str.main()

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
	const pair = 2

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for actors := 0; actors < runtime.NumCPU(); actors++ {
		strings := make(chan string, 1)
		runes := make(chan []rune)

		strings <- str.data

		wg.Add(pair)

		go str.runer(wg, strings, runes)
		go str.stringer(wg, runes, strings)
	}
}

func (str *Stress) runer(
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
			}
		}
	}
}

func (str *Stress) stringer(
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
			}
		}
	}
}
