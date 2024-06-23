package app

import (
	"fmt"
	"yubigo-pass/internal/app/crypto"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/services"
	"yubigo-pass/internal/cli"
	"yubigo-pass/internal/database"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/google/uuid"
)

type programAction int

const (
	loginAction programAction = iota
	createUserAction
	mainMenuAction
	getPasswordAction
	viewPasswordsItem
	addPasswordItem
)

// Runner is a program flow execution controller
type Runner struct {
	currentAction    programAction
	serviceContainer services.Container
}

// NewRunner returns new Runner instance
func NewRunner(serviceContainer services.Container) *Runner {
	return &Runner{
		currentAction:    loginAction,
		serviceContainer: serviceContainer,
	}
}

// Run runs the application
func (r *Runner) Run() error {
	defer database.CloseDB()

	for {
		switch r.currentAction {
		case loginAction:
			m, err := runLoginAction(r.serviceContainer)
			if err != nil {
				return fmt.Errorf("login action failed: %s", err)
			}
			if m.CreateUserPicked {
				r.currentAction = createUserAction
				continue
			}
			if m.LoggedIn {
				r.currentAction = mainMenuAction
				continue
			}
			return nil

		case createUserAction:
			m, err := runCreateUserAction(r.serviceContainer)
			if err != nil {
				return fmt.Errorf("create user action failed: %s", err)
			}
			if m.UserCreated || m.UserCreationAborted {
				r.currentAction = loginAction
				continue
			}
			return nil

		case mainMenuAction:
			m, err := runMainMenuAction(r.serviceContainer)
			if err != nil {
				return fmt.Errorf("main menu action failed: %s", err)
			}
			switch m.Choice {
			case cli.GetPasswordItem:
				r.currentAction = getPasswordAction
				continue
			case cli.ViewPasswordItem:
				r.currentAction = viewPasswordsItem
				continue
			case cli.AddPasswordItem:
				r.currentAction = addPasswordItem
				continue
			case cli.LogoutItem:
				r.currentAction = loginAction
				continue
			}
			return nil
		}
	}
}

func runMainMenuAction(serviceContainer services.Container) (cli.MainMenuModel, error) {
	m, err := tea.NewProgram(serviceContainer.Models.MainMenu).Run()
	mainMenuModel := m.(cli.MainMenuModel)
	if err != nil {
		return mainMenuModel, fmt.Errorf("could not start main menu action: %w", err)
	}
	return mainMenuModel, nil
}

func runLoginAction(serviceContainer services.Container) (cli.LoginModel, error) {
	m, err := tea.NewProgram(serviceContainer.Models.Login).Run()
	loginModel := m.(cli.LoginModel)
	if err != nil {
		return loginModel, fmt.Errorf("could not start login action: %w", err)
	}
	return loginModel, err
}

func runCreateUserAction(serviceContainer services.Container) (cli.CreateUserModel, error) {
	m, err := tea.NewProgram(serviceContainer.Models.CreateUser).Run()
	createUserModel := m.(cli.CreateUserModel)
	if err != nil {
		return createUserModel, fmt.Errorf("could not start create user action: %w", err)
	}
	if createUserModel.UserCreationAborted || createUserModel.Cancelled {
		return createUserModel, nil
	}

	err = createNewUser(serviceContainer, createUserModel)
	if err != nil {
		return createUserModel, err
	}

	return createUserModel, nil
}

func createNewUser(serviceContainer services.Container, m tea.Model) error {
	userUUID := uuid.New().String()
	username, password := cli.ExtractDataFromModel(m)
	salt, err := crypto.NewSalt()
	if err != nil {
		return err
	}
	passwordHash := crypto.HashPasswordWithSalt(password, salt)

	createUserInput := model.NewUser(userUUID, username, passwordHash, salt)
	err = serviceContainer.Store.CreateUser(createUserInput)
	if err != nil {
		return err
	}

	return nil
}
