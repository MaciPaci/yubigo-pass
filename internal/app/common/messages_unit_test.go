//go:build unit

package common

import (
	"errors"
	"fmt"
	"testing"
	"yubigo-pass/internal/app/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoginCmd verifies that LoginsCmd creates the correct LoginSuccessMsg.
func TestLoginCmd(t *testing.T) {
	expectedUsername := "user123"
	expectedPassphrase := "passphrase"

	cmd := LoginCmd(expectedUsername, expectedPassphrase)
	require.NotNil(t, cmd, "Command should not be nil")

	msg := cmd()
	resultMsg, ok := msg.(LoginMsg)
	require.True(t, ok, "Message should be of type LoginSuccessMsg")

	assert.Equal(t, expectedUsername, resultMsg.Username)
	assert.Equal(t, expectedPassphrase, resultMsg.Password)
}

// TestCreateUserCmd verifies that CreateUserCmd creates the correct UserToCreateMsg.
func TestCreateUserCmd(t *testing.T) {
	expectedUsername := "newuser"
	expectedPassword := "newpassword"

	cmd := CreateUserCmd(expectedUsername, expectedPassword)
	require.NotNil(t, cmd, "Command should not be nil")

	msg := cmd()
	resultMsg, ok := msg.(UserToCreateMsg)
	require.True(t, ok, "Message should be of type UserToCreateMsg")

	assert.Equal(t, expectedUsername, resultMsg.Username)
	assert.Equal(t, expectedPassword, resultMsg.Password)
}

// TestAddPasswordCmd verifies that AddPasswordCmd creates the correct PasswordToAddMsg.
func TestAddPasswordCmd(t *testing.T) {
	passwordData := model.Password{
		UserID:   "uid",
		Title:    "Test Title",
		Username: "pwduser",
		Password: "pwd",
		Url:      "http://example.com",
		Nonce:    []byte("nonce"),
	}

	expectedData := PasswordToAddMsg{
		Title:    passwordData.Title,
		Username: passwordData.Username,
		Password: passwordData.Password,
		Url:      passwordData.Url,
	}

	cmd := AddPasswordCmd(passwordData)
	require.NotNil(t, cmd, "Command should not be nil")

	msg := cmd()
	resultMsg, ok := msg.(PasswordToAddMsg)
	require.True(t, ok, "Message should be of type PasswordToAddMsg")

	assert.Equal(t, expectedData, resultMsg)
}

// TestChangeStateCmd verifies that ChangeStateCmd creates the correct StateMsg.
func TestChangeStateCmd(t *testing.T) {
	testCases := []MsgState{
		StateLogout,
		StateQuit,
		StateGoBack,
	}

	for _, expectedState := range testCases {
		t.Run(fmt.Sprintf("State_%d", expectedState), func(t *testing.T) {
			cmd := ChangeStateCmd(expectedState)
			require.NotNil(t, cmd, "Command should not be nil")

			msg := cmd()
			resultMsg, ok := msg.(StateMsg)
			require.True(t, ok, "Message should be of type StateMsg")

			assert.Equal(t, expectedState, resultMsg.State)
		})
	}
}

// TestErrCmd verifies that ErrCmd creates the correct ErrorMsg.
func TestErrCmd(t *testing.T) {
	expectedErr := errors.New("this is a test error")

	cmd := ErrCmd(expectedErr)
	require.NotNil(t, cmd, "Command should not be nil")

	msg := cmd()
	resultMsg, ok := msg.(ErrorMsg)
	require.True(t, ok, "Message should be of type ErrorMsg")

	assert.Equal(t, expectedErr, resultMsg.Err)
	assert.ErrorIs(t, resultMsg.Err, expectedErr)
}
