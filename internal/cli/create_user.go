package cli

import (
	"fmt"
	"strings"
	"yubigo-pass/internal/database"

	"github.com/mritd/bubbles/common"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionStateCreateUser uint

const (
	createUserState sessionStateCreateUser = iota
	backButtonFocused
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	focusedSubmitButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredSubmitButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	focusedBackButton   = focusedStyle.Copy().Render("[ Back ]")
	blurredBackButton   = fmt.Sprintf("[ %s ]", blurredStyle.Render("Back"))
)

const (
	validateOkPrefix  = "✔"
	validateErrPrefix = "✘"
	colorValidateOk   = "2"
	colorValidateErr  = "1"
)

// CreateUserModel is a model for user creation
type CreateUserModel struct {
	state               sessionStateCreateUser
	focusIndex          int
	inputs              []textinput.Model
	showErr             bool
	err                 error
	Cancelled           bool
	UserCreated         bool
	UserCreationAborted bool

	store database.StoreExecutor
}

// ExtractDataFromModel maps data from the model into strings
func ExtractDataFromModel(m tea.Model) (string, string) {
	return m.(CreateUserModel).inputs[0].Value(), m.(CreateUserModel).inputs[1].Value()
}

// NewCreateUserModel returns model for user creation
func NewCreateUserModel(store database.StoreExecutor) CreateUserModel {
	m := CreateUserModel{
		state:  createUserState,
		inputs: make([]textinput.Model, 2),
		store:  store,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Username"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 64
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}
		m.inputs[i] = t
	}

	return m
}

// Init initializes for CreateUserModel
func (m CreateUserModel) Init() tea.Cmd {
	m.inputs[0].Focus()
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}
	return textinput.Blink
}

// Update updates CreateUserModel
func (m CreateUserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Cancelled = true
			return m, tea.Quit

		case tea.KeyTab, tea.KeyShiftTab:
			if m.state == createUserState {
				m.state = backButtonFocused
				for i := 0; i <= len(m.inputs)-1; i++ {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}
			} else {
				var cmd tea.Cmd
				m.state = createUserState
				if m.focusIndex != len(m.inputs) {
					m.inputs[m.focusIndex].PromptStyle = focusedStyle
					m.inputs[m.focusIndex].TextStyle = focusedStyle
					cmd = m.inputs[m.focusIndex].Focus()
				}
				return m, cmd
			}

		case tea.KeyEnter, tea.KeyUp, tea.KeyDown, tea.KeyPgDown:
			key := msg.Type

			if m.state == backButtonFocused {
				if key == tea.KeyEnter {
					m.UserCreationAborted = true
					return m, tea.Quit
				}
			} else {
				if key == tea.KeyEnter && m.focusIndex == len(m.inputs) {
					if m.err == nil {
						user, _ := m.store.GetUser(m.inputs[0].Value())
						if user.Username == "" {
							m.UserCreated = true
							return m, tea.Quit
						}
						m.err = fmt.Errorf("username already exists")
					}
					m.showErr = true
				}

				// Cycle indexes
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

				cmds := make([]tea.Cmd, len(m.inputs))
				if m.state == createUserState {
					for i := 0; i <= len(m.inputs)-1; i++ {
						if i == m.focusIndex {
							// Set focused state
							cmds[i] = m.inputs[i].Focus()
							m.inputs[i].PromptStyle = focusedStyle
							m.inputs[i].TextStyle = focusedStyle
							continue
						}
						// Remove focused state
						m.inputs[i].Blur()
						m.inputs[i].PromptStyle = noStyle
						m.inputs[i].TextStyle = noStyle
					}
				}

				return m, tea.Batch(cmds...)
			}
		case tea.KeyRunes:
			m.showErr = false
			m.err = nil
		}
	}
	cmd := updateCreateUserModelInputs(&m, msg)
	m.err = validateCreateUserModelInputs(m.inputs, m.err)

	return m, cmd
}

func updateCreateUserModelInputs(m *CreateUserModel, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// View renders CreateUserModel
func (m CreateUserModel) View() string {
	var b strings.Builder

	b.WriteString("\n----------CREATE NEW USER----------\n\n")

	var screenMsg string
	if m.err != nil {
		if m.showErr {
			screenMsg = common.FontColor(fmt.Sprintf("%s ERROR: %s\n", validateErrPrefix, m.err.Error()), colorValidateErr)
		}
	}

	if m.UserCreated {
		screenMsg = common.FontColor(fmt.Sprintf("%s User created successfully\n", validateOkPrefix), colorValidateOk)
	}

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	createUserButton := &blurredSubmitButton
	backButton := &blurredBackButton
	if m.focusIndex == len(m.inputs) && m.state == createUserState {
		createUserButton = &focusedSubmitButton
		backButton = &blurredBackButton
	}
	if m.state == backButtonFocused {
		createUserButton = &blurredSubmitButton
		backButton = &focusedBackButton
	}

	fmt.Fprintf(&b, "\n\n%s\t%s\n%s\n", *createUserButton, *backButton, screenMsg)

	return b.String()
}

func validateCreateUserModelInputs(input []textinput.Model, err error) error {
	if err != nil {
		return err
	}

	usernameIsEmpty := func() bool { return strings.TrimSpace(input[0].Value()) == "" }
	passwordIsEmpty := func() bool { return strings.TrimSpace(input[1].Value()) == "" }

	if usernameIsEmpty() && passwordIsEmpty() {
		return fmt.Errorf("username and password cannot be empty")
	}
	if usernameIsEmpty() {
		return fmt.Errorf("username cannot be empty")
	}
	if passwordIsEmpty() {
		return fmt.Errorf("password cannot be empty")
	}
	return nil
}
