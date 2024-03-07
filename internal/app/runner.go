package app

import (
	"fmt"
	"os"
	"yubigo-pass/internal/app/crypto"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/services"
	"yubigo-pass/internal/cli"
	"yubigo-pass/internal/database"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type programAction int

const (
	loginAction programAction = iota
	createUserAction
)

// Runner is a program flow execution controller
type Runner struct {
	currentAction programAction
}

// NewRunner returns new Runner instance
func NewRunner() Runner {
	return Runner{loginAction}
}

// Run runs the application
func (r Runner) Run() {
	serviceContainer := services.Build()
	defer database.CloseDB()

	for {
		switch r.currentAction {
		case loginAction:
			createUserActionPicked, err := runLoginAction(serviceContainer)
			if err != nil {
				logrus.Errorf("login action failed: %s:\n", err)
				os.Exit(1)
			}
			if createUserActionPicked {
				r.currentAction = createUserAction
				continue
			}
			return

		case createUserAction:
			userCreated, err := runCreateUserAction(serviceContainer)
			if err != nil {
				logrus.Errorf("create user action failed: %s:\n", err)
				os.Exit(1)
			}
			if userCreated {
				r.currentAction = loginAction
				continue
			}
			return
		}
	}
}

func runLoginAction(serviceContainer services.Container) (bool, error) {
	m, err := tea.NewProgram(serviceContainer.Models.Login).Run()
	return m.(cli.LoginModel).CreateUserActionPicked(), err
}

func runCreateUserAction(serviceContainer services.Container) (bool, error) {
	m, err := tea.NewProgram(serviceContainer.Models.CreateUser).Run()
	if err != nil {
		return false, fmt.Errorf("could not start program: %w", err)
	}
	if m.(cli.CreateUserModel).WasCancelled() {
		return false, fmt.Errorf("create user action cancelled")
	}

	err = createNewUser(serviceContainer, m)
	if err != nil {
		return false, err
	}

	return m.(cli.CreateUserModel).WasUserCreated(), nil
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
