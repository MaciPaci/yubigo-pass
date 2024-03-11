package cli

import (
	"fmt"
	"strings"
	"yubigo-pass/internal/app/crypto"
	"yubigo-pass/internal/database"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mritd/bubbles/common"
)

type sessionStateLogin uint

const (
	loginState sessionStateLogin = iota
	createUserButtonFocused
)

var (
	focusedLoginButton      = focusedStyle.Copy().Render("[ Login ]")
	blurredLoginButton      = fmt.Sprintf("[ %s ]", blurredStyle.Render("Login"))
	focusedCreateUserButton = focusedStyle.Copy().Render("[ Create new user ]")
	blurredCreateUserButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Create new user"))
)

// LoginModel is a model for user login
type LoginModel struct {
	state            sessionStateLogin
	focusIndex       int
	inputs           []textinput.Model
	showErr          bool
	err              error
	LoggedIn         bool
	Cancelled        bool
	CreateUserPicked bool

	store database.StoreExecutor
}

// NewLoginModel returns model for user creation
func NewLoginModel(store database.StoreExecutor) LoginModel {
	m := LoginModel{
		state:  loginState,
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

// Init initializes for LoginModel
func (m LoginModel) Init() tea.Cmd {
	m.inputs[0].Focus()
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}
	return textinput.Blink
}

// Update updates LoginModel
func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Cancelled = true
			return m, tea.Quit

		case tea.KeyTab, tea.KeyShiftTab:
			if m.state == loginState {
				m.state = createUserButtonFocused
				for i := 0; i <= len(m.inputs)-1; i++ {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}
			} else {
				var cmd tea.Cmd
				m.state = loginState
				if m.focusIndex != len(m.inputs) {
					m.inputs[m.focusIndex].PromptStyle = focusedStyle
					m.inputs[m.focusIndex].TextStyle = focusedStyle
					cmd = m.inputs[m.focusIndex].Focus()
				}
				return m, cmd
			}

		case tea.KeyEnter, tea.KeyUp, tea.KeyDown, tea.KeyPgDown:
			key := msg.Type
			if m.state == createUserButtonFocused {
				if key == tea.KeyEnter {
					m.CreateUserPicked = true
					return m, tea.Quit
				}
			} else {
				if key == tea.KeyEnter && m.focusIndex == len(m.inputs) {
					if m.err == nil {
						user, _ := m.store.GetUser(m.inputs[0].Value())
						if user.Username == "" {
							m.err = fmt.Errorf("incorrect credentials")
						} else {
							hashedPassword := crypto.HashPasswordWithSalt(m.inputs[1].Value(), user.Salt)
							if hashedPassword == user.Password {
								m.LoggedIn = true
								return m, tea.Quit
							} else {
								m.err = fmt.Errorf("incorrect credentials")
							}
						}
					}
					m.showErr = true
				}

				if key == tea.KeyUp || key == tea.KeyShiftTab {
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
				if m.state == loginState {
					for i := 0; i <= len(m.inputs)-1; i++ {
						if i == m.focusIndex {
							cmds[i] = m.inputs[i].Focus()
							m.inputs[i].PromptStyle = focusedStyle
							m.inputs[i].TextStyle = focusedStyle
							continue
						}
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
	cmd := updateLoginModelInputs(&m, msg)
	m.err = validateLoginModelInputs(m.inputs, m.err)

	return m, cmd
}

func updateLoginModelInputs(m *LoginModel, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// View renders LoginModel
func (m LoginModel) View() string {
	var b strings.Builder

	b.WriteString("\n---------------LOGIN---------------\n\n")

	var screenMsg string
	if m.err != nil {
		if m.showErr {
			screenMsg = common.FontColor(fmt.Sprintf("%s ERROR: %s\n", validateErrPrefix, m.err.Error()), colorValidateErr)
		}
	}

	if m.LoggedIn {
		screenMsg = common.FontColor(fmt.Sprintf("%s Logged in\n", validateOkPrefix), colorValidateOk)
	}

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	loginButton := &blurredLoginButton
	createUserButton := &blurredCreateUserButton
	if m.focusIndex == len(m.inputs) && m.state == loginState {
		loginButton = &focusedLoginButton
		createUserButton = &blurredCreateUserButton
	}
	if m.state == createUserButtonFocused {
		loginButton = &blurredLoginButton
		createUserButton = &focusedCreateUserButton
	}

	fmt.Fprintf(&b, "\n\n%s\t%s\n%s\n", *loginButton, *createUserButton, screenMsg)

	return b.String()
}

func validateLoginModelInputs(input []textinput.Model, err error) error {
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
