package cli

import (
	"errors"
	"fmt"
	"strings"
	"yubigo-pass/internal/app/common"
	"yubigo-pass/internal/app/crypto"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/utils"
	"yubigo-pass/internal/database"

	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// sessionStateLogin defines the focus state within the login view.
type sessionStateLogin uint

const (
	loginInputsFocused sessionStateLogin = iota
	createUserButtonFocused
)

var (
	focusedLoginButton      = focusedStyle.Copy().Render("[ Login ]")
	blurredLoginButton      = fmt.Sprintf("[ %s ]", blurredStyle.Render("Login"))
	focusedCreateUserButton = focusedStyle.Copy().Render("[ Create new user ]")
	blurredCreateUserButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Create new user"))
)

// LoginModel is a Bubble Tea model for the user login screen.
// It handles user input for credentials and triggers authentication logic via messages.
type LoginModel struct {
	state      sessionStateLogin
	focusIndex int
	inputs     []textinput.Model
	showErr    bool
	err        error

	store database.StoreExecutor
}

// NewLoginModel creates a new instance of the LoginModel.
func NewLoginModel(store database.StoreExecutor) LoginModel {
	m := LoginModel{
		state:  loginInputsFocused,
		inputs: make([]textinput.Model, 2),
		store:  store,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 64

		switch i {
		case 0:
			t.Placeholder = "Username"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.PromptStyle = noStyle
			t.TextStyle = noStyle
		}
		m.inputs[i] = t
	}
	m.focusIndex = 0

	return m
}

// Init initializes the LoginModel, setting focus and clearing inputs.
func (m LoginModel) Init() tea.Cmd {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = noStyle
		m.inputs[i].TextStyle = noStyle
	}
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = focusedStyle
	m.inputs[0].TextStyle = focusedStyle
	return textinput.Blink
}

// Update handles incoming messages and user input for the login screen.
// It sends messages to the main application model for state transitions or actions.
func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == loginInputsFocused && m.focusIndex < len(m.inputs) {
			switch msg.Type {
			case tea.KeyRunes, tea.KeySpace, tea.KeyBackspace:
				m.showErr = false
				m.err = nil
			}
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, common.ChangeStateCmd(common.StateQuit)

		case tea.KeyTab, tea.KeyShiftTab:
			if m.state == loginInputsFocused {
				m.state = createUserButtonFocused
			} else {
				m.state = loginInputsFocused
			}
			cmds = append(cmds, m.updateFocus())

		case tea.KeyUp, tea.KeyDown:
			if m.state == loginInputsFocused {
				originalFocus := m.focusIndex
				if msg.Type == tea.KeyUp {
					m.focusIndex = (m.focusIndex - 1 + (len(m.inputs) + 1)) % (len(m.inputs) + 1)
				} else {
					m.focusIndex = (m.focusIndex + 1) % (len(m.inputs) + 1)
				}
				if m.focusIndex != originalFocus {
					cmds = append(cmds, m.updateFocus())
				}
			}

		case tea.KeyEnter:
			if m.state == createUserButtonFocused {
				return m, common.ChangeStateCmd(common.StateGoToCreateUser)
			}
			if m.state == loginInputsFocused && m.focusIndex == len(m.inputs) {
				validationErr := validateLoginModelInputs(m.inputs)
				if validationErr != nil {
					m.err = validationErr
					m.showErr = true
					return m, nil
				}

				user, err := m.store.GetUser(m.inputs[0].Value())
				if err != nil {
					if errors.As(err, &model.UserNotFoundError{}) {
						m.err = fmt.Errorf("incorrect username or password")
					} else {
						m.err = fmt.Errorf("login failed: %w", err)
					}
					m.showErr = true
					return m, nil
				}

				hashedPassword := crypto.HashPasswordWithSalt(m.inputs[1].Value(), user.Salt)
				if hashedPassword == user.Password {
					session := utils.NewSession(user.UserID, m.inputs[1].Value(), user.Salt)
					return m, common.LogSuccessCmd(session)
				}

				m.err = fmt.Errorf("incorrect username or password")
				m.showErr = true
				return m, nil

			} else if m.state == loginInputsFocused && m.focusIndex < len(m.inputs) {
				m.focusIndex++
				cmds = append(cmds, m.updateFocus())
			}
		}
	}

	if m.state == loginInputsFocused && m.focusIndex < len(m.inputs) {
		var inputCmd tea.Cmd
		m.inputs[m.focusIndex], inputCmd = m.inputs[m.focusIndex].Update(msg)
		cmds = append(cmds, inputCmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the login screen UI.
func (m LoginModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("LOGIN") + "\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}

	loginButton := blurredLoginButton
	createUserButton := blurredCreateUserButton

	if m.state == loginInputsFocused && m.focusIndex == len(m.inputs) {
		loginButton = focusedLoginButton
	}
	if m.state == createUserButtonFocused {
		createUserButton = focusedCreateUserButton
	}

	fmt.Fprintf(&b, "\n%s\t%s\n", loginButton, createUserButton)

	if m.err != nil && m.showErr {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateErr))
		fmt.Fprintf(&b, "\n%s %s\n", validateErrPrefix, errorStyle.Render(m.err.Error()))
	}

	help := blurredStyle.Render("\n(Tab/Shift+Tab: Navigate, ↑/↓: Cycle Focus, Enter: Select/Login, Esc: Quit)")
	b.WriteString(help)

	return b.String()
}

// updateFocus updates the visual focus styles on inputs and returns the blink command.
func (m *LoginModel) updateFocus() tea.Cmd {
	for i := range m.inputs {
		if m.state == loginInputsFocused && i == m.focusIndex {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
	}
	if m.state == loginInputsFocused && m.focusIndex < len(m.inputs) {
		return textinput.Blink
	}
	return nil
}

// validateLoginModelInputs checks if the required input fields are non-empty.
func validateLoginModelInputs(input []textinput.Model) error {
	usernameIsEmpty := func() bool { return strings.TrimSpace(input[0].Value()) == "" }
	passwordIsEmpty := func() bool { return strings.TrimSpace(input[1].Value()) == "" }

	if usernameIsEmpty() || passwordIsEmpty() {
		return fmt.Errorf("username and password cannot be empty")
	}
	return nil
}
