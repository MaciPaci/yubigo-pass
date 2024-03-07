package services

import (
	"yubigo-pass/internal/cli"
	"yubigo-pass/internal/database"

	tea "github.com/charmbracelet/bubbletea"
)

// TeaModels is a struct holding all Bubbletea models
type TeaModels struct {
	Login      tea.Model
	CreateUser tea.Model
}

// InitTeaModels initializes all Bubbletea models
func InitTeaModels(store database.StoreExecutor) TeaModels {
	return TeaModels{
		Login:      cli.NewLoginModel(store),
		CreateUser: cli.NewCreateUserModel(store),
	}
}
