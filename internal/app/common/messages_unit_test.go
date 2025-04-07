//go:build unit

package common

import (
	"errors"
	"fmt"
	"testing"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLogSuccessCmd verifies that LogSuccessCmd creates the correct LoginSuccessMsg.
func TestLogSuccessCmd(t *testing.T) {
	expectedUserID := "user123"
	expectedPassphrase := "passphrase"
	expectedSalt := "salt123"
	session := utils.NewSession(expectedUserID, expectedPassphrase, expectedSalt)

	cmd := LogSuccessCmd(session)
	require.NotNil(t, cmd, "Command should not be nil")

	msg := cmd()
	resultMsg, ok := msg.(LoginSuccessMsg)
	require.True(t, ok, "Message should be of type LoginSuccessMsg")

	assert.Equal(t, session, resultMsg.Session)
	assert.Equal(t, expectedUserID, resultMsg.Session.GetUserID())
	assert.Equal(t, expectedPassphrase, resultMsg.Session.GetPassphrase())
	assert.Equal(t, expectedSalt, resultMsg.Session.GetSalt())
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
	expectedData := model.Password{
		UserID:   "uid",
		Title:    "Test Title",
		Username: "pwduser",
		Password: "pwd",
		Url:      "http://example.com",
		Nonce:    []byte("nonce"),
	}

	cmd := AddPasswordCmd(expectedData)
	require.NotNil(t, cmd, "Command should not be nil")

	msg := cmd()
	resultMsg, ok := msg.(PasswordToAddMsg)
	require.True(t, ok, "Message should be of type PasswordToAddMsg")

	assert.Equal(t, expectedData, resultMsg.Data)
}

// TestChangeStateCmd verifies that ChangeStateCmd creates the correct StateMsg.
func TestChangeStateCmd(t *testing.T) {
	testCases := []MsgState{
		StateGoToMainMenu,
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
