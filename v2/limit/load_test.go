package limit

import (
	"crypto/rand"
	"runtime"
	"sync"

	"github.com/akramarenkov/cqos/v2/internal/breaker"
)

const defaultDataAmount = 10000

func getRandom(amount int) ([]byte, error) {
	random := make([]byte, amount)

	if _, err := rand.Read(random); err != nil {
		return nil, err
	}

	return random, nil
}

type loadSystem struct {
	breaker *breaker.Breaker
	data    string
}

func newLoadSystem(amount int) (*loadSystem, error) {
	if amount <= 0 {
		amount = defaultDataAmount
	}

	random, err := getRandom(amount)
	if err != nil {
		return nil, err
	}

	lds := &loadSystem{
		breaker: breaker.New(),
		data:    string(random),
	}

	go lds.loop()

	return lds, nil
}

func (lds *loadSystem) Stop() {
	lds.breaker.Break()
}

func (lds *loadSystem) loop() {
	defer lds.breaker.Complete()

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for pair := 0; pair < runtime.NumCPU(); pair++ {
		input := make(chan string, 1)
		output := make(chan []rune)

		input <- lds.data

		wg.Add(2)

		go lds.sender(input, output, wg)
		go lds.receiver(output, input, wg)
	}
}

func (lds *loadSystem) sender(
	input chan string,
	output chan []rune,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-lds.breaker.Breaked():
			return
		case data := <-input:
			select {
			case <-lds.breaker.Breaked():
				return
			case output <- []rune(data):
			}
		}
	}
}

func (lds *loadSystem) receiver(
	input chan []rune,
	output chan string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-lds.breaker.Breaked():
			return
		case data := <-input:
			select {
			case <-lds.breaker.Breaked():
				return
			case output <- string(data):
			}
		}
	}
}
