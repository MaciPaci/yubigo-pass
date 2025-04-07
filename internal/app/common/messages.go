package common

import (
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/utils"

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

// LoginSuccessMsg carries the user session details upon successful authentication.
type LoginSuccessMsg struct {
	Session utils.Session
}

// UserToCreateMsg carries the necessary data for initiating the user creation process.
type UserToCreateMsg struct {
	Username string
	Password string
}

// PasswordToAddMsg carries the necessary data for initiating the password creation process.
type PasswordToAddMsg struct {
	Data model.Password
}

// LogSuccessCmd returns a command that sends a LoginSuccessMsg.
func LogSuccessCmd(s utils.Session) tea.Cmd {
	return func() tea.Msg {
		return LoginSuccessMsg{Session: s}
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
		return PasswordToAddMsg{Data: data}
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
