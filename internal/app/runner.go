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
		}
	}
}

func runLoginAction(serviceContainer services.Container) (cli.LoginModel, error) {
	m, err := tea.NewProgram(serviceContainer.Models.Login).Run()
	loginModel := m.(cli.LoginModel)
	if err != nil {
		return loginModel, fmt.Errorf("could not start login action: %w", err)
	}
	if loginModel.Cancelled {
		return loginModel, fmt.Errorf("login action cancelled")
	}
	return loginModel, err
}

func runCreateUserAction(serviceContainer services.Container) (cli.CreateUserModel, error) {
	m, err := tea.NewProgram(serviceContainer.Models.CreateUser).Run()
	createUserModel := m.(cli.CreateUserModel)
	if err != nil {
		return createUserModel, fmt.Errorf("could not start create user action: %w", err)
	}
	if m.(cli.CreateUserModel).Cancelled {
		return createUserModel, fmt.Errorf("create user action cancelled")
	}
	if m.(cli.CreateUserModel).UserCreationAborted {
		return createUserModel, nil
	}

	err = createNewUser(serviceContainer, m)
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
