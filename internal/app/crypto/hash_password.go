package crypto

import (
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// HashPasswordWithSalt salts and hashes the given password using bcrypt
func HashPasswordWithSalt(password string) (string, string, error) {
	salt := RandomStringWithLength(32)
	combined := []byte(password + salt)

	hash := sha256.Sum256(combined)
	hashString := base64.URLEncoding.EncodeToString(hash[:])

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(hashString), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return string(hashedPassword), salt, nil
}
