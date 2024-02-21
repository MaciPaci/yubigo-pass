package crypto

import "math/rand"

const stringBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomStringWithLength generates random string with a given length
func RandomStringWithLength(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = stringBytes[rand.Intn(len(stringBytes))]
	}
	return string(b)
}
