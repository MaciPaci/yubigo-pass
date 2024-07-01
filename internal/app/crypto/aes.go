package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/argon2"
	"io"
)

func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32) // 32 bytes for AES-256
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

func DeriveAESKey(passphrase, salt string) []byte {
	iterations := 3     // Number of passes
	memory := 32 * 1024 // Memory usage in KB (e.g., 32 MB)
	parallelism := 4    // Number of parallel threads
	keyLength := 32     // Length of the AES-256 key (32 bytes for AES-256 key)

	return argon2.IDKey([]byte(passphrase), []byte(salt), uint32(iterations), uint32(memory), uint8(parallelism), uint32(keyLength))
}

func EncryptAES(key []byte, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// Create a new GCM cipher.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	// Generate a random nonce.
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	// EncryptAES the plaintext using AES-GCM.
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

// DecryptAES decrypts ciphertext with AES-256-GCM.
func DecryptAES(key []byte, ciphertext []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM cipher.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract nonce from the beginning of the ciphertext.
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce := ciphertext[:nonceSize]
	passphrase := ciphertext[nonceSize:]

	// DecryptAES the ciphertext using AES-GCM.
	plaintext, err = gcm.Open(nil, nonce, passphrase, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
