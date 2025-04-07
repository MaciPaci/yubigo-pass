package cli

import (
	"errors"
	"fmt"
	"strings"
	"yubigo-pass/internal/app/common"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/database"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// sessionStateCreateUser defines the focus state within the create user view.
type sessionStateCreateUser uint

const (
	createUserInputsFocused sessionStateCreateUser = iota
	createUserBackFocused
)

// CreateUserModel is a Bubble Tea model for the new user creation screen.
// It handles user input for credentials and triggers user creation logic via messages.
type CreateUserModel struct {
	state           sessionStateCreateUser
	focusIndex      int
	inputs          []textinput.Model
	showErr         bool
	err             error
	passwordVisible bool

	store database.StoreExecutor
}

// ExtractUserDataFromModel retrieves the username and password from the model's inputs.
// It should be called *after* the model has successfully gathered input, typically
// by the component handling the creation logic (like AppModel).
func ExtractUserDataFromModel(m CreateUserModel) (string, string) {
	return m.inputs[0].Value(), m.inputs[1].Value()
}

// NewCreateUserModel creates a new instance of the CreateUserModel.
func NewCreateUserModel(store database.StoreExecutor) CreateUserModel {
	m := CreateUserModel{
		state:           createUserInputsFocused,
		inputs:          make([]textinput.Model, 2),
		store:           store,
		passwordVisible: false,
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

// Init initializes the CreateUserModel, setting focus and clearing inputs.
func (m CreateUserModel) Init() tea.Cmd {
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

// Update handles incoming messages and user input for the create user screen.
// It sends messages to the main application model for state transitions or actions.
func (m CreateUserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == createUserInputsFocused && m.focusIndex < len(m.inputs) {
			switch msg.Type {
			case tea.KeyRunes, tea.KeySpace, tea.KeyBackspace:
				m.showErr = false
				m.err = nil
			}
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, common.ChangeStateCmd(common.StateQuit)

		case tea.KeyCtrlS:
			if m.focusIndex == 1 {
				m.passwordVisible = !m.passwordVisible
				if m.passwordVisible {
					m.inputs[1].EchoMode = textinput.EchoNormal
				} else {
					m.inputs[1].EchoMode = textinput.EchoPassword
				}
				return m, nil
			}

		case tea.KeyTab, tea.KeyShiftTab:
			if m.state == createUserInputsFocused {
				m.state = createUserBackFocused
			} else {
				m.state = createUserInputsFocused
			}
			cmds = append(cmds, m.updateFocus())

		case tea.KeyUp, tea.KeyDown:
			if m.state == createUserInputsFocused {
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
			if m.state == createUserBackFocused {
				return m, common.ChangeStateCmd(common.StateGoBack)
			}
			if m.state == createUserInputsFocused && m.focusIndex == len(m.inputs) {
				validationErr := validateCreateUserModelInputs(m.inputs)
				if validationErr != nil {
					m.err = validationErr
					m.showErr = true
					return m, nil
				}

				_, checkErr := m.store.GetUser(m.inputs[0].Value())
				if checkErr == nil {
					m.err = model.NewUserAlreadyExistsError(m.inputs[0].Value())
					m.showErr = true
					return m, nil
				} else if !errors.As(checkErr, &model.UserNotFoundError{}) {
					m.err = fmt.Errorf("failed to check username: %w", checkErr)
					m.showErr = true
					return m, nil
				}

				return m, common.CreateUserCmd(m.inputs[0].Value(), m.inputs[1].Value())

			} else if m.state == createUserInputsFocused && m.focusIndex < len(m.inputs) {
				m.focusIndex++
				cmds = append(cmds, m.updateFocus())
			}
		}
	}

	if m.state == createUserInputsFocused && m.focusIndex < len(m.inputs) {
		var inputCmd tea.Cmd
		m.inputs[m.focusIndex], inputCmd = m.inputs[m.focusIndex].Update(msg)
		cmds = append(cmds, inputCmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the create user screen UI.
func (m CreateUserModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("CREATE NEW USER") + "\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}

	submitButton := blurredSubmitButton
	backButton := blurredBackButton

	if m.state == createUserInputsFocused && m.focusIndex == len(m.inputs) {
		submitButton = focusedSubmitButton
	}
	if m.state == createUserBackFocused {
		backButton = focusedBackButton
	}

	fmt.Fprintf(&b, "\n%s\t%s\n", submitButton, backButton)

	if m.err != nil && m.showErr {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateErr))
		fmt.Fprintf(&b, "\n%s %s\n", validateErrPrefix, errorStyle.Render(m.err.Error()))
	}

	help := blurredStyle.Render("\n(Tab/Shift+Tab: Navigate, ↑/↓: Cycle Focus, Ctrl+S on Pwd: Show/Hide, Enter: Select/Submit, Esc: Quit)")
	b.WriteString(help)

	return b.String()
}

// updateFocus updates the visual focus styles on inputs and returns the blink command.
func (m *CreateUserModel) updateFocus() tea.Cmd {
	for i := range m.inputs {
		if m.state == createUserInputsFocused && i == m.focusIndex {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
	}
	if m.state == createUserInputsFocused && m.focusIndex < len(m.inputs) {
		return textinput.Blink
	}
	return nil
}

// validateCreateUserModelInputs checks if the required input fields are non-empty.
func validateCreateUserModelInputs(input []textinput.Model) error {
	usernameIsEmpty := func() bool { return strings.TrimSpace(input[0].Value()) == "" }
	passwordIsEmpty := func() bool { return strings.TrimSpace(input[1].Value()) == "" }

	if usernameIsEmpty() || passwordIsEmpty() {
		return fmt.Errorf("username and password cannot be empty")
	}
	return nil
}
