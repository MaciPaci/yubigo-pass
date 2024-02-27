package utils

import (
	"crypto/rand"
	"math/big"
)

const stringBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomStringWithLength generates a random string with a given length
func RandomStringWithLength(n int) (string, error) {
	b := make([]byte, n)
	for i := range b {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(stringBytes))))
		if err != nil {
			return "", err
		}
		b[i] = stringBytes[r.Int64()]
	}
	return string(b), nil
}
