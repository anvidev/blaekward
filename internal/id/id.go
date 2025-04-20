// Package id provides a function to generate random id of any positive length, with the added flexibility of using
// predefined sets or user defined sets.
package id

import (
	"crypto/rand"
	"math/big"
)

const (
	LettersLower string = "abcdefghijklmnopqrstuvxyz"
	LettersUpper        = "ABCDEFGHIJKLMNOPQRSTUVXYZ"
	Numbers             = "1234567890"
)

// String returns an [n] length generated id, or "" if it encounters an error.
func String(n int, sets ...string) string {
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
			return ""
		}
		result[i] = charSet[randIndex.Int64()]
	}

	return string(result)
}
