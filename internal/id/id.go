package id

import (
	"crypto/rand"
	"math/big"
)

type set string

const (
	LettersLower set = "abcdefghijklmnopqrstuvxyz"
	LetterUpper      = "ABCDEFGHIJKLMNOPQRSTUVXYZ"
	Numbers          = "1234567890"
)

func Generate(n int, sets ...set) string {
	if n <= 0 {
		n = 32
	}

	var charSet string
	for _, set := range sets {
		charSet += string(set)
	}

	setLen := big.NewInt(int64(len(charSet)))
	result := make([]byte, n)

	for i := range result {
		randIndex, err := rand.Int(rand.Reader, setLen)
		if err != nil {
			panic(err)
		}
		result[i] = charSet[randIndex.Int64()]
	}

	return string(result)
}
