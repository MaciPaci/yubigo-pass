package cli

import (
	"errors"
	"fmt"
	"strings"
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
	addPasswordFocused sessionStateAddPassword = iota
	addPasswordBackButtonFocused
)

// AddPasswordModel is a model for adding a new password.
type AddPasswordModel struct {
	state            sessionStateAddPassword
	focusIndex       int
	inputs           []textinput.Model
	showErr          bool
	err              error
	Cancelled        bool
	Back             bool
	PasswordAdded    bool
	passwordStrength int // Stores the zxcvbn score (0-4)

	store   database.StoreExecutor
	session utils.Session
}

// ExtractPasswordDataFromModel maps data from the model into Password struct.
func ExtractPasswordDataFromModel(m tea.Model) model.Password {
	addModel, ok := m.(AddPasswordModel)
	if !ok {
		return model.Password{}
	}
	return model.Password{
		Title:    addModel.inputs[0].Value(),
		Username: addModel.inputs[1].Value(),
		Password: addModel.inputs[2].Value(),
		Url:      addModel.inputs[3].Value(),
	}
}

// NewAddPasswordModel returns model for adding a new password.
func NewAddPasswordModel(store database.StoreExecutor, session utils.Session) AddPasswordModel {
	m := AddPasswordModel{
		state:            addPasswordFocused,
		inputs:           make([]textinput.Model, 4),
		store:            store,
		session:          session,
		passwordStrength: 0,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Title"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 64
		case 1:
			t.Placeholder = "Username"
			t.PromptStyle = blurredStyle
			t.TextStyle = blurredStyle
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.PromptStyle = blurredStyle
			t.TextStyle = blurredStyle
			t.CharLimit = 0
		case 3:
			t.Placeholder = "URL (optional)"
			t.PromptStyle = blurredStyle
			t.TextStyle = blurredStyle
			t.CharLimit = 256
		}
		m.inputs[i] = t
	}

	return m
}

// Init initializes AddPasswordModel.
func (m AddPasswordModel) Init() tea.Cmd {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = blurredStyle
		m.inputs[i].TextStyle = blurredStyle
	}
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = focusedStyle
	m.inputs[0].TextStyle = focusedStyle

	return textinput.Blink
}

// Update updates the AddPasswordModel based on user input.
func (m AddPasswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == addPasswordFocused && msg.Type == tea.KeyCtrlG {
			generatedPassword, err := utils.GeneratePassword(utils.DefaultLength, true, true, true, true)
			if err != nil {
				m.err = fmt.Errorf("password generation failed: %w", err)
				m.showErr = true
				return m, nil
			}
			passwordInputIndex := 2
			m.inputs[passwordInputIndex].SetValue(generatedPassword)
			m.inputs[passwordInputIndex].CursorEnd()
			m.focusIndex = passwordInputIndex
			updateFocus(&m)
			m.passwordStrength = calculateStrength(&m)
			cmds = append(cmds, textinput.Blink)
			return m, tea.Batch(cmds...)
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Cancelled = true
			return m, tea.Quit
		}

		switch m.state {
		case addPasswordFocused:
			switch msg.Type {
			case tea.KeyTab:
				m.state = addPasswordBackButtonFocused
				if m.focusIndex >= 0 && m.focusIndex < len(m.inputs) {
					m.inputs[m.focusIndex].Blur()
					m.inputs[m.focusIndex].PromptStyle = blurredStyle
					m.inputs[m.focusIndex].TextStyle = blurredStyle
				}

			case tea.KeyShiftTab:
				m.state = addPasswordBackButtonFocused
				if m.focusIndex >= 0 && m.focusIndex < len(m.inputs) {
					m.inputs[m.focusIndex].Blur()
					m.inputs[m.focusIndex].PromptStyle = blurredStyle
					m.inputs[m.focusIndex].TextStyle = blurredStyle
				}

			case tea.KeyEnter:
				if m.focusIndex == len(m.inputs) {
					m.err = validateAddPasswordModelInputs(m.inputs, m.err)
					if m.err != nil {
						m.showErr = true
					} else {
						_, err := m.store.GetPassword(m.session.GetUserID(), m.inputs[0].Value(), m.inputs[1].Value())
						if err != nil && errors.As(err, &model.PasswordNotFoundError{}) {
							m.PasswordAdded = true
							return m, tea.Quit
						} else if err == nil {
							m.err = fmt.Errorf("password entry with this title/username already exists")
							m.showErr = true
						} else {
							m.err = fmt.Errorf("failed to check for existing password: %w", err)
							m.showErr = true
						}
					}
					return m, nil
				} else {
					m.focusIndex++
					if m.focusIndex > len(m.inputs) {
						m.focusIndex = len(m.inputs)
					}
					updateFocus(&m)
					cmds = append(cmds, textinput.Blink)
				}

			case tea.KeyUp, tea.KeyDown:
				key := msg.Type
				if key == tea.KeyUp {
					m.focusIndex--
				} else {
					m.focusIndex++
				}
				if m.focusIndex > len(m.inputs) {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs)
				}
				updateFocus(&m)
				cmds = append(cmds, textinput.Blink)

			case tea.KeyRunes, tea.KeySpace, tea.KeyBackspace:
				m.showErr = false
				m.err = nil
			}

		case addPasswordBackButtonFocused:
			switch msg.Type {
			case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown:
				m.state = addPasswordFocused
				updateFocus(&m)
				cmds = append(cmds, textinput.Blink)

			case tea.KeyEnter:
				m.Back = true
				return m, tea.Quit
			}
		}
	}

	if m.focusIndex >= 0 && m.focusIndex < len(m.inputs) {
		inputCmds := updateAddPasswordModelInputs(&m, msg)
		cmds = append(cmds, inputCmds)
	}

	if msgCouldChangePassword(msg) {
		m.passwordStrength = calculateStrength(&m)
	}

	return m, tea.Batch(cmds...)
}

// updateFocus is a helper to update input focus styles and return focus command.
func updateFocus(m *AddPasswordModel) {
	isAddButtonFocused := m.focusIndex == len(m.inputs)

	for i := 0; i < len(m.inputs); i++ {
		shouldFocus := m.state == addPasswordFocused && !isAddButtonFocused && i == m.focusIndex

		if shouldFocus {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = focusedStyle
			m.inputs[i].TextStyle = focusedStyle
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = blurredStyle
			m.inputs[i].TextStyle = blurredStyle
		}
	}
}

// updateAddPasswordModelInputs updates the text input fields.
func updateAddPasswordModelInputs(m *AddPasswordModel, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// msgCouldChangePassword checks if a message could have changed the password input.
func msgCouldChangePassword(msg tea.Msg) bool {
	switch msg.(type) {
	case tea.KeyMsg:
		return true
	default:
		return false
	}
}

// calculateStrength calculates the password strength score using zxcvbn.
func calculateStrength(m *AddPasswordModel) int {
	password := m.inputs[2].Value()
	if password != "" {
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
	return 0
}

// View renders AddPasswordModel.
func (m AddPasswordModel) View() string {
	if m.Cancelled {
		return quitTextStyle.Render("Quitting.")
	}
	if m.Back {
		return quitTextStyle.Render("Going back...")
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("ADD A NEW PASSWORD"))

	var screenMsg string
	if m.err != nil && m.showErr {
		screenMsg = fmt.Sprintf("\n%s %s",
			validateErrPrefix,
			lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateErr)).Render("ERROR: "+m.err.Error()),
		)
	} else if m.PasswordAdded {
		screenMsg = fmt.Sprintf("\n%s %s",
			validateOkPrefix,
			lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateOk)).Render("Password validation OK. Saving..."),
		)
	}

	b.WriteString("\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())

		if i == 2 && m.inputs[i].Value() != "" {
			strengthScore := m.passwordStrength
			strengthText := utils.GetStrengthText(strengthScore)
			strengthStyle := utils.GetStrengthStyle(strengthScore)
			strengthIndicator := strengthStyle.Render(fmt.Sprintf("[%s]", strengthText))
			b.WriteString("  " + strengthIndicator)
		}
		b.WriteRune('\n')
	}

	addBtn := blurredAddButton
	backBtn := blurredBackButton

	if m.state == addPasswordFocused {
		if m.focusIndex == len(m.inputs) {
			addBtn = focusedAddButton
		}
	} else if m.state == addPasswordBackButtonFocused {
		backBtn = focusedBackButton
	}

	buttonRow := lipgloss.JoinHorizontal(lipgloss.Top, addBtn, "    ", backBtn)
	fmt.Fprintf(&b, "\n%s", buttonRow)

	if screenMsg != "" {
		b.WriteString(screenMsg + "\n")
	}

	help := " (Tab/Shift+Tab: Navigate, Enter: Select/Confirm, Ctrl+G: Generate Password, Esc: Cancel)"
	b.WriteString(blurredStyle.Render("\n\n" + help))

	return b.String()
}

// validateAddPasswordModelInputs checks if required fields are empty.
func validateAddPasswordModelInputs(input []textinput.Model, existingErr error) error {
	if existingErr != nil && !strings.Contains(existingErr.Error(), "cannot be empty") {
		return existingErr
	}

	titleIsEmpty := func() bool { return strings.TrimSpace(input[0].Value()) == "" }
	usernameIsEmpty := func() bool { return strings.TrimSpace(input[1].Value()) == "" }
	passwordIsEmpty := func() bool { return strings.TrimSpace(input[2].Value()) == "" }

	if titleIsEmpty() || usernameIsEmpty() || passwordIsEmpty() {
		return fmt.Errorf("title, username, and password fields cannot be empty")
	}

	return nil
}
