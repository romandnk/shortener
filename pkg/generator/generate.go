package generator

//go:generate mockgen -source=generate.go -destination=mock/mock.go generate

import (
	"crypto/rand"
	"math/big"
)

const (
	alphabet string = "_0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Generator interface {
	Random() (string, error)
}

type Gen struct {
	Length int
}

func NewGen(length int) *Gen {
	return &Gen{Length: length}
}

func (g *Gen) Random() (string, error) {
	charsetLength := big.NewInt(int64(len(alphabet)))
	result := make([]byte, g.Length)

	for i := 0; i < g.Length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}

		result[i] = alphabet[randomIndex.Int64()]
	}

	return string(result), nil
}
