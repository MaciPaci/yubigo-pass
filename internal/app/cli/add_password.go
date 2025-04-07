package cli

import (
	"errors"
	"fmt"
	"strings"
	"yubigo-pass/internal/app/common"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/utils"
	"yubigo-pass/internal/database"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nbutton23/zxcvbn-go"
)

// sessionStateAddPassword defines the focus state within the add password view.
type sessionStateAddPassword uint

const (
	addPasswordInputsFocused sessionStateAddPassword = iota
	addPasswordBackFocused
)

// AddPasswordModel is a Bubble Tea model for adding a new password entry.
// It gathers input, allows toggling password visibility, and triggers the
// password addition process via messages.
type AddPasswordModel struct {
	state            sessionStateAddPassword
	focusIndex       int
	inputs           []textinput.Model
	showErr          bool
	err              error
	passwordStrength int
	passwordVisible  bool

	store   database.StoreExecutor
	session utils.Session
}

// ExtractPasswordDataFromModel creates a model.Password struct from the input fields.
func ExtractPasswordDataFromModel(m AddPasswordModel) model.Password {
	return model.Password{
		Title:    m.inputs[0].Value(),
		Username: m.inputs[1].Value(),
		Password: m.inputs[2].Value(),
		Url:      m.inputs[3].Value(),
	}
}

// NewAddPasswordModel creates a new instance of the AddPasswordModel.
func NewAddPasswordModel(store database.StoreExecutor, session utils.Session) AddPasswordModel {
	m := AddPasswordModel{
		state:            addPasswordInputsFocused,
		inputs:           make([]textinput.Model, 4),
		store:            store,
		session:          session,
		passwordStrength: 0,
		passwordVisible:  false,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 0

		switch i {
		case 0:
			t.Placeholder = "Title"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 128
		case 1:
			t.Placeholder = "Username"
			t.PromptStyle = noStyle
			t.TextStyle = noStyle
			t.CharLimit = 128
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.PromptStyle = noStyle
			t.TextStyle = noStyle
		case 3:
			t.Placeholder = "URL (optional)"
			t.PromptStyle = noStyle
			t.TextStyle = noStyle
			t.CharLimit = 512
		}
		m.inputs[i] = t
	}
	m.focusIndex = 0

	return m
}

// Init initializes the AddPasswordModel, setting focus and clearing inputs.
func (m AddPasswordModel) Init() tea.Cmd {
	m.inputs[2].EchoMode = textinput.EchoPassword
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

// Update handles incoming messages and user input for the add password screen.
func (m AddPasswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == addPasswordInputsFocused && m.focusIndex < len(m.inputs) {
			switch msg.Type {
			case tea.KeyRunes, tea.KeySpace, tea.KeyBackspace:
				m.showErr = false
				m.err = nil
			case tea.KeyCtrlG:
				if m.focusIndex == 2 {
					generatedPassword, err := utils.GeneratePassword(utils.DefaultLength, true, true, true, true)
					if err != nil {
						m.err = fmt.Errorf("password generation failed: %w", err)
						m.showErr = true
						return m, nil
					}
					m.inputs[2].SetValue(generatedPassword)
					m.inputs[2].CursorEnd()
					m.passwordStrength = calculateStrength(&m)
					m.showErr = false
					m.err = nil
					return m, textinput.Blink
				}
			case tea.KeyCtrlS:
				if m.focusIndex == 2 {
					m.passwordVisible = !m.passwordVisible
					if m.passwordVisible {
						m.inputs[2].EchoMode = textinput.EchoNormal
					} else {
						m.inputs[2].EchoMode = textinput.EchoPassword
					}
					return m, nil
				}
			}
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, common.ChangeStateCmd(common.StateQuit)

		case tea.KeyTab, tea.KeyShiftTab:
			if m.state == addPasswordInputsFocused {
				m.state = addPasswordBackFocused
			} else {
				m.state = addPasswordInputsFocused
			}
			cmds = append(cmds, m.updateFocus())

		case tea.KeyUp, tea.KeyDown:
			if m.state == addPasswordInputsFocused {
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
			if m.state == addPasswordBackFocused {
				return m, common.ChangeStateCmd(common.StateGoBack)
			}
			if m.state == addPasswordInputsFocused && m.focusIndex == len(m.inputs) {
				validationErr := validateAddPasswordModelInputs(m.inputs)
				if validationErr != nil {
					m.err = validationErr
					m.showErr = true
					return m, nil
				}

				_, checkErr := m.store.GetPassword(m.session.GetUserID(), m.inputs[0].Value(), m.inputs[1].Value())
				if checkErr == nil {
					m.err = model.NewPasswordAlreadyExistsError(m.session.GetUserID(), m.inputs[0].Value(), m.inputs[1].Value())
					m.showErr = true
					return m, nil
				} else if !errors.As(checkErr, &model.PasswordNotFoundError{}) {
					m.err = fmt.Errorf("failed to check existing password: %w", checkErr)
					m.showErr = true
					return m, nil
				}

				passwordData := ExtractPasswordDataFromModel(m)

				return m, common.AddPasswordCmd(passwordData)

			} else if m.state == addPasswordInputsFocused && m.focusIndex < len(m.inputs) {
				m.focusIndex++
				cmds = append(cmds, m.updateFocus())
			}
		}
	}

	if m.state == addPasswordInputsFocused && m.focusIndex < len(m.inputs) {
		var inputCmd tea.Cmd
		originalValue := m.inputs[m.focusIndex].Value()
		m.inputs[m.focusIndex], inputCmd = m.inputs[m.focusIndex].Update(msg)
		cmds = append(cmds, inputCmd)

		if m.focusIndex == 2 && m.inputs[2].Value() != originalValue {
			m.passwordStrength = calculateStrength(&m)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the add password screen UI.
func (m AddPasswordModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("ADD A NEW PASSWORD") + "\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i == 2 && m.inputs[i].Value() != "" {
			strengthScore := m.passwordStrength
			strengthText := utils.GetStrengthText(strengthScore)
			strengthStyle := utils.GetStrengthStyle(strengthScore)
			strengthIndicator := strengthStyle.Render(fmt.Sprintf(" [%s]", strengthText))
			b.WriteString(strengthIndicator)
		}
		b.WriteRune('\n')
	}

	addBtn := blurredAddButton
	backBtn := blurredBackButton

	if m.state == addPasswordInputsFocused && m.focusIndex == len(m.inputs) {
		addBtn = focusedAddButton
	}
	if m.state == addPasswordBackFocused {
		backBtn = focusedBackButton
	}

	buttonRow := lipgloss.JoinHorizontal(lipgloss.Top, addBtn, "    ", backBtn)
	fmt.Fprintf(&b, "\n%s", buttonRow)

	if m.err != nil && m.showErr {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateErr))
		fmt.Fprintf(&b, "\n\n%s %s", validateErrPrefix, errorStyle.Render(m.err.Error()))
	}

	help := blurredStyle.Render("\n\n(Tab/Shift+Tab: Navigate, ↑/↓: Focus, Enter: Select/Add)\n")
	help += blurredStyle.Render("(Ctrl+G on Pwd: Generate, Ctrl+S on Pwd: Show/Hide, Esc: Quit)")
	b.WriteString(help)

	return b.String()
}

// updateFocus updates the visual focus styles on inputs and returns the blink command.
func (m *AddPasswordModel) updateFocus() tea.Cmd {
	for i := 0; i < len(m.inputs); i++ {
		if m.state == addPasswordInputsFocused && i == m.focusIndex {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
	}
	if m.state == addPasswordInputsFocused && m.focusIndex < len(m.inputs) {
		return textinput.Blink
	}
	return nil
}

// calculateStrength calculates the password strength score using zxcvbn.
func calculateStrength(m *AddPasswordModel) int {
	password := m.inputs[2].Value()
	if password == "" {
		return 0
	}
	userInputs := []string{
		m.inputs[0].Value(),
		m.inputs[1].Value(),
	}
	var filteredInputs []string
	for _, input := range userInputs {
		if input != "" {
			filteredInputs = append(filteredInputs, input)
		}
	}
	return zxcvbn.PasswordStrength(password, filteredInputs).Score
}

// validateAddPasswordModelInputs checks if required fields are empty.
func validateAddPasswordModelInputs(input []textinput.Model) error {
	titleIsEmpty := func() bool { return strings.TrimSpace(input[0].Value()) == "" }
	usernameIsEmpty := func() bool { return strings.TrimSpace(input[1].Value()) == "" }
	passwordIsEmpty := func() bool { return strings.TrimSpace(input[2].Value()) == "" }

	if titleIsEmpty() || usernameIsEmpty() || passwordIsEmpty() {
		return fmt.Errorf("title, username, and password fields cannot be empty")
	}
	return nil
}
