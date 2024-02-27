package services

import (
	"yubigo-pass/internal/cli"
	"yubigo-pass/internal/database"

	tea "github.com/charmbracelet/bubbletea"
)

// Programs is a struct holding all Bubbletea programs
type Programs struct {
	CreateUserProgram *tea.Program
}

// InitPrograms initializes all programs
func InitPrograms(store database.StoreExecutor) Programs {
	return Programs{
		CreateUserProgram: tea.NewProgram(cli.NewCreateUserModel(store)),
	}
}
