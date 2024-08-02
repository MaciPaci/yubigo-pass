package utils

import (
	"testing"
	"yubigo-pass/test"

	"github.com/stretchr/testify/assert"
)

func TestNewEmptySession(t *testing.T) {
	// given
	session := NewEmptySession()

	// then
	assert.Equal(t, "", session.GetUserID())
	assert.Equal(t, "", session.GetPassphrase())
	assert.Equal(t, "", session.GetUserID())
}

func TestNewSession(t *testing.T) {
	// given
	userID := test.RandomString()
	userPassword := test.RandomString()
	userSalt := test.RandomString()

	// when
	session := NewSession(userID, userPassword, userSalt)

	// then
	assert.Equal(t, userID, session.GetUserID())
	assert.Equal(t, userPassword, session.GetPassphrase())
	assert.Equal(t, userSalt, session.GetSalt())
}

func TestGetSessionParameters(t *testing.T) {
	// given
	userID := test.RandomString()
	userPassword := test.RandomString()
	userSalt := test.RandomString()

	session := NewSession(userID, userPassword, userSalt)

	// when
	userID2, userPassword2, userSalt2 := session.GetSessionData()

	// then
	assert.Equal(t, userID, userID2)
	assert.Equal(t, userPassword, userPassword2)
	assert.Equal(t, userSalt, userSalt2)
}

func TestClearSessionParameters(t *testing.T) {
	// given
	userID := test.RandomString()
	userPassword := test.RandomString()
	userSalt := test.RandomString()

	session := NewSession(userID, userPassword, userSalt)

	// when
	session.Clear()

	// then
	assert.Equal(t, "", session.GetUserID())
	assert.Equal(t, "", session.GetPassphrase())
	assert.Equal(t, "", session.GetSalt())
}
