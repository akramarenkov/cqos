package stress

import "crypto/rand"

func getRandom(amount int) ([]byte, error) {
	random := make([]byte, amount)

	if _, err := rand.Read(random); err != nil {
		return nil, err
	}

	return random, nil
}
