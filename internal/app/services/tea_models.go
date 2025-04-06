package services

import (
	"yubigo-pass/internal/app/utils"
	"yubigo-pass/internal/cli"
	"yubigo-pass/internal/database"

	tea "github.com/charmbracelet/bubbletea"
)

// TeaModels is a struct holding all Bubbletea models
type TeaModels struct {
	Login            tea.Model
	CreateUser       tea.Model
	MainMenu         tea.Model
	AddPasswordModel tea.Model
}

// InitTeaModels initializes all Bubbletea models
func InitTeaModels(store database.StoreExecutor, session utils.Session) TeaModels {
	return TeaModels{
		Login:            cli.NewLoginModel(store),
		CreateUser:       cli.NewCreateUserModel(store),
		MainMenu:         cli.NewMainMenuModel(),
		AddPasswordModel: cli.NewAddPasswordModel(store, session),
	}
}
