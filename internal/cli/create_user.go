package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mritd/bubbles/common"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle   = focusedStyle.Copy()
	noStyle       = lipgloss.NewStyle()
	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

const (
	validateOkPrefix  = "✔"
	validateErrPrefix = "✘"
	colorValidateOk   = "2"
	colorValidateErr  = "1"
)

// CreateUserModel is a model for user creation
type CreateUserModel struct {
	focusIndex int
	inputs     []textinput.Model
	//cursorMode cursor.Mode
	showErr   bool
	err       error
	finished  bool
	cancelled bool
}

func printFields(m CreateUserModel) {
	for i := range m.inputs {
		fmt.Println(m.inputs[i].Value())
	}
}

// NewCreateUserModel returns model for user creation
func NewCreateUserModel() CreateUserModel {
	m := CreateUserModel{
		inputs: make([]textinput.Model, 2),
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
	return textinput.Blink
}

// Update updates CreateUserModel
func (m CreateUserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.finished {
			printFields(m)
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelled = true
			return m, tea.Quit

		case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown, tea.KeyPgDown:
			key := msg.Type

			if key == tea.KeyEnter && m.focusIndex == len(m.inputs) {
				if m.err == nil {
					m.finished = true
					return m, common.Done
				}
				m.showErr = true
			}

			// Cycle indexes
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

			return m, tea.Batch(cmds...)

		case tea.KeyRunes:
			m.showErr = false
			m.err = nil
		}
	}
	cmd := updateInputs(&m, msg)
	m.err = validateInputs(m.inputs)

	return m, cmd
}

func updateInputs(m *CreateUserModel, msg tea.Msg) tea.Cmd {
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

	var screenMsg string
	if m.err != nil {
		if m.showErr {
			screenMsg = common.FontColor(fmt.Sprintf("%s ERROR: %s\n", validateErrPrefix, m.err.Error()), colorValidateErr)
		}
	}

	if m.finished {
		screenMsg = common.FontColor(fmt.Sprintf("%s User created successfully\n", validateOkPrefix), colorValidateOk)
	}

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n%s\n", *button, screenMsg)

	return b.String()
}

func validateInputs(input []textinput.Model) error {
	if strings.TrimSpace(input[0].Value()) == "" && strings.TrimSpace(input[1].Value()) == "" {
		return errors.New("username and password cannot be empty")
	}
	if strings.TrimSpace(input[0].Value()) == "" {
		return errors.New("username cannot be empty")
	}
	if strings.TrimSpace(input[1].Value()) == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}
