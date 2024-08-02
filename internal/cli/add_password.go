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
	"github.com/mritd/bubbles/common"
)

type sessionStateAddPassword uint

const (
	addPasswordFocused sessionStateAddPassword = iota
	addPasswordBackButtonFocused
)

// AddPasswordModel is a model for adding a new password
type AddPasswordModel struct {
	state         sessionStateAddPassword
	focusIndex    int
	inputs        []textinput.Model
	showErr       bool
	err           error
	Cancelled     bool
	Back          bool
	PasswordAdded bool

	store   database.StoreExecutor
	session utils.Session
}

// ExtractPasswordDataFromModel maps data from the model into Password struct
func ExtractPasswordDataFromModel(m tea.Model) model.Password {
	return model.Password{
		Title:    m.(AddPasswordModel).inputs[0].Value(),
		Username: m.(AddPasswordModel).inputs[1].Value(),
		Password: m.(AddPasswordModel).inputs[2].Value(),
		Url:      m.(AddPasswordModel).inputs[3].Value(),
	}
}

// NewAddPasswordModel returns model for adding a new password
func NewAddPasswordModel(store database.StoreExecutor, session utils.Session) AddPasswordModel {
	m := AddPasswordModel{
		state:   addPasswordFocused,
		inputs:  make([]textinput.Model, 4),
		store:   store,
		session: session,
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
		case 3:
			t.Placeholder = "URL (optional)"
			t.PromptStyle = blurredStyle
			t.TextStyle = blurredStyle
			t.CharLimit = 64
		}
		m.inputs[i] = t
	}

	return m
}

// Init initializes AddPasswordModel
func (m AddPasswordModel) Init() tea.Cmd {
	m.inputs[0].Focus()
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}
	return textinput.Blink
}

// Update updates AddPasswordModel
func (m AddPasswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Cancelled = true
			return m, tea.Quit

		case tea.KeyTab, tea.KeyShiftTab:
			if m.state == addPasswordFocused {
				m.state = addPasswordBackButtonFocused
				for i := 0; i <= len(m.inputs)-1; i++ {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}
			} else {
				var cmd tea.Cmd
				m.state = addPasswordFocused
				if m.focusIndex != len(m.inputs) {
					m.inputs[m.focusIndex].PromptStyle = focusedStyle
					m.inputs[m.focusIndex].TextStyle = focusedStyle
					cmd = m.inputs[m.focusIndex].Focus()
				}
				return m, cmd
			}

		case tea.KeyEnter, tea.KeyUp, tea.KeyDown, tea.KeyPgDown:
			key := msg.Type
			if m.state == addPasswordBackButtonFocused {
				if key == tea.KeyEnter {
					m.Back = true
					return m, tea.Quit
				}
			} else {
				if key == tea.KeyEnter && m.focusIndex == len(m.inputs) {
					if m.err == nil {
						_, err := m.store.GetPassword(m.session.GetUserID(), m.inputs[0].Value(), m.inputs[1].Value())
						if err != nil && errors.As(err, &model.PasswordNotFoundError{}) {
							m.PasswordAdded = true
							return m, tea.Quit
						}
						m.err = fmt.Errorf("this password already exists, change inputs or update existing password")
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
				if m.state == addPasswordFocused {
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
	cmd := updateAddPasswordModelInputs(&m, msg)
	m.err = validateAddPasswordModelInputs(m.inputs, m.err)

	return m, cmd
}

func updateAddPasswordModelInputs(m *AddPasswordModel, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// View renders AddPasswordModel
func (m AddPasswordModel) View() string {
	if m.Cancelled {
		return quitTextStyle.Render("Quitting.")
	}
	var b strings.Builder

	b.WriteString(titleStyle.Render("ADD A NEW PASSWORD") + "\n\n")

	var screenMsg string
	if m.err != nil {
		if m.showErr {
			screenMsg = common.FontColor(fmt.Sprintf("%s ERROR: %s\n", validateErrPrefix, m.err.Error()), colorValidateErr)
		}
	}

	if m.PasswordAdded {
		screenMsg = common.FontColor(fmt.Sprintf("\n%s Password added.\n", validateOkPrefix), colorValidateOk)
	}

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	addButton := &blurredAddButton
	backButton := &blurredBackButton
	if m.focusIndex == len(m.inputs) && m.state == addPasswordFocused {
		addButton = &focusedAddButton
		backButton = &blurredBackButton
	}
	if m.state == addPasswordBackButtonFocused {
		addButton = &blurredAddButton
		backButton = &focusedBackButton
	}

	fmt.Fprintf(&b, "\n\n%s\t\t%s\n%s\n", *addButton, *backButton, screenMsg)

	return b.String()
}

func validateAddPasswordModelInputs(input []textinput.Model, err error) error {
	if err != nil {
		return err
	}

	titleIsEmpty := func() bool { return strings.TrimSpace(input[0].Value()) == "" }
	usernameIsEmpty := func() bool { return strings.TrimSpace(input[1].Value()) == "" }
	passwordIsEmpty := func() bool { return strings.TrimSpace(input[2].Value()) == "" }

	if titleIsEmpty() || usernameIsEmpty() || passwordIsEmpty() {
		return fmt.Errorf("only optional fields can be empty")
	}
	return nil
}
