package crypto

import (
	"crypto/sha256"
	"encoding/base64"
)

// HashPasswordWithSalt salts and hashes the given password
func HashPasswordWithSalt(password, salt string) string {
	combined := []byte(password + salt)
	hash := sha256.Sum256(combined)
	return base64.URLEncoding.EncodeToString(hash[:])
}

// NewSalt returns new random salt
func NewSalt() (string, error) {
	return RandomStringWithLength(32)
}
