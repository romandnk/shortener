package generate

import (
	"crypto/rand"
	"github.com/romandnk/shortener/internal/constant"
	"math/big"
)

const (
	alphabet string = "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func Random() (string, error) {
	charsetLength := big.NewInt(int64(len(alphabet)))
	result := make([]byte, constant.AliasLength)

	for i := 0; i < constant.AliasLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}

		result[i] = alphabet[randomIndex.Int64()]
	}

	return string(result), nil
}
