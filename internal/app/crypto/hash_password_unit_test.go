//go:build unit

package crypto

import (
	"testing"
	"yubigo-pass/test"

	"github.com/stretchr/testify/assert"
)

func TestHashPasswordWithSaltShouldReturnTheSameHashEveryTime(t *testing.T) {
	// given
	password := test.RandomString()
	salt := NewSalt()

	// when
	hashedPassword := HashPasswordWithSalt(password, salt)
	hashedPasswordSecondTime := HashPasswordWithSalt(password, salt)

	// then
	assert.Equal(t, hashedPassword, hashedPasswordSecondTime)
}
