package common

import (
	"yubigo-pass/internal/app/model"

	tea "github.com/charmbracelet/bubbletea"
)

// MsgState represents distinct states or transition signals within the application.
type MsgState int

// Constants defining different application states or transition signals.
const (
	StateError MsgState = iota
	StateLoginSuccess
	StateGoToCreateUser
	StateUserCreated
	StateGoToMainMenu
	StateGoToAddPassword
	StateGoToViewPasswords
	StateGoToGetPassword
	StatePasswordAdded
	StatePasswordCopied
	StateGoBack
	StateLogout
	StateQuit
)

// StateMsg is a generic message used to signal state transitions identified by MsgState.
type StateMsg struct {
	State MsgState
}

// ErrorMsg encapsulates an error to be potentially displayed or handled by the UI.
type ErrorMsg struct {
	Err error
}

// LoginMsg carries the user details upon authentication attempt.
type LoginMsg struct {
	Username string
	Password string
}

// UserToCreateMsg carries the necessary data for initiating the user creation process.
type UserToCreateMsg struct {
	Username string
	Password string
}

// PasswordToAddMsg carries the necessary data for initiating the password creation process.
type PasswordToAddMsg struct {
	Title    string
	Username string
	Password string
	Url      string
}

// DecryptPasswordMsg requests decryption of a specific password entry.
type DecryptPasswordMsg struct {
	PasswordID string
}

// DecryptAndCopyPasswordMsg requests decryption and copying of a specific password.
type DecryptAndCopyPasswordMsg struct {
	PasswordID string
}

// PasswordDecryptedMsg sends back the plaintext of a requested password.
type PasswordDecryptedMsg struct {
	PasswordID string
	Plaintext  string
}

// LoginCmd returns a command that sends a LoginMsg.
func LoginCmd(username, password string) tea.Cmd {
	return func() tea.Msg {
		return LoginMsg{Username: username, Password: password}
	}
}

// CreateUserCmd returns a command that sends a UserToCreateMsg.
func CreateUserCmd(username, password string) tea.Cmd {
	return func() tea.Msg {
		return UserToCreateMsg{Username: username, Password: password}
	}
}

// AddPasswordCmd returns a command that sends a PasswordToAddMsg.
func AddPasswordCmd(data model.Password) tea.Cmd {
	return func() tea.Msg {
		return PasswordToAddMsg{
			Title:    data.Title,
			Username: data.Username,
			Password: data.Password,
			Url:      data.Url,
		}
	}
}

// ChangeStateCmd returns a command that sends a generic StateMsg to trigger a state change.
func ChangeStateCmd(newState MsgState) tea.Cmd {
	return func() tea.Msg {
		return StateMsg{State: newState}
	}
}

// ErrCmd returns a command that sends an ErrorMsg containing the provided error.
func ErrCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{Err: err}
	}
}

// RequestDecryptPasswordCmd returns a command requesting decryption for a password ID.
func RequestDecryptPasswordCmd(passwordID string) tea.Cmd {
	return func() tea.Msg {
		return DecryptPasswordMsg{PasswordID: passwordID}
	}
}

// RequestDecryptAndCopyPasswordCmd returns a command requesting decryption and copy for a password ID.
func RequestDecryptAndCopyPasswordCmd(passwordID string) tea.Cmd {
	return func() tea.Msg {
		return DecryptAndCopyPasswordMsg{PasswordID: passwordID}
	}
}

// PasswordDecryptedCmd returns a command sending back the decrypted password.
func PasswordDecryptedCmd(passwordID, plaintext string) tea.Cmd {
	return func() tea.Msg {
		return PasswordDecryptedMsg{PasswordID: passwordID, Plaintext: plaintext}
	}
}
