package test

import (
	"crypto/rand"
	"math/big"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const stringLength = 32

// RandomString generates a random string for testing purposes
func RandomString() string {
	b := make([]byte, stringLength)
	for i := range b {
		r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		b[i] = letterBytes[r.Int64()]
	}
	return string(b)
}

// RandomStringWithLength generates a random string with a given length for testing purposes
func RandomStringWithLength(length int) string {
	b := make([]byte, length)
	for i := range b {
		r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		b[i] = letterBytes[r.Int64()]
	}
	return string(b)
}
