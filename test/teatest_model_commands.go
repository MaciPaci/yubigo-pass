package test

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

// PressKey sends press key message to Bubbletea test model
func PressKey(tm *teatest.TestModel, key tea.KeyType) {
	tm.Send(tea.KeyMsg{
		Type: key,
	})
}

// TypeString sends string message to Bubbletea test model
func TypeString(tm *teatest.TestModel, message string) {
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(message),
	})
}
