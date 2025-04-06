//go:build unit

package crypto

import (
	"testing"
	"yubigo-pass/test"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAESKey(t *testing.T) {
	// when
	key, err := GenerateAESKey()

	// then
	assert.Nil(t, err)
	assert.Len(t, key, 32)
}

func TestDeriveAESKey(t *testing.T) {
	// given
	password := test.RandomString()
	salt := test.RandomString()

	// when
	key := DeriveAESKey(password, salt)

	// then
	assert.Len(t, key, 32)
}

func TestEncryptDecryptAES(t *testing.T) {
	// given
	password := test.RandomString()
	salt := test.RandomString()
	key := DeriveAESKey(password, salt)
	textToEncrypt := test.RandomString()

	// when
	encryptedText, nonce, err := EncryptAES(key, []byte(textToEncrypt))
	if err != nil {
		t.Fatalf("Failed to encrypt text: %v", err)
	}
	ciphertext := append(nonce, encryptedText...)
	decryptedText, err := DecryptAES(key, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt text: %v", err)
	}

	// then
	assert.Equal(t, textToEncrypt, string(decryptedText))
}

func TestEncryptDecryptAESWithWrongKey(t *testing.T) {
	// given
	password := test.RandomString()
	salt := test.RandomString()
	key := DeriveAESKey(password, salt)
	textToEncrypt := test.RandomString()

	//when
	encryptedText, _, err := EncryptAES(key, []byte(textToEncrypt))
	if err != nil {
		t.Fatalf("Failed to encrypt textToEncrypt: %v", err)
	}
	wrongKey := []byte(test.RandomString())
	decryptedText, err := DecryptAES(wrongKey, encryptedText)

	//expected
	expectedError := "cipher: message authentication failed"

	// then
	assert.EqualError(t, err, expectedError)
	assert.Empty(t, decryptedText)
}

func TestEncryptDecryptAESWithEmptyKey(t *testing.T) {
	// given
	password := test.RandomString()
	salt := test.RandomString()
	key := DeriveAESKey(password, salt)
	textToEncrypt := test.RandomString()

	//when
	encryptedText, _, err := EncryptAES(key, []byte(textToEncrypt))
	if err != nil {
		t.Fatalf("Failed to encrypt textToEncrypt: %v", err)
	}
	wrongKey := []byte("")
	decryptedText, err := DecryptAES(wrongKey, encryptedText)

	//expected
	expectedError := "crypto/aes: invalid key size 0"

	// then
	assert.EqualError(t, err, expectedError)
	assert.Empty(t, decryptedText)
}
