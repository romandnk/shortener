package generate

import (
	"crypto/rand"
	"math/big"
)

const (
	stringLength int    = 10
	alphabet     string = "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func Random() (string, error) {
	charsetLength := big.NewInt(int64(len(alphabet)))
	result := make([]byte, stringLength)

	for i := 0; i < stringLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}

		result[i] = alphabet[randomIndex.Int64()]
	}

	return string(result), nil
}
