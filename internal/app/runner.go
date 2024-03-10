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
			createUserActionPicked, err := runLoginAction(r.serviceContainer)
			if err != nil {
				return fmt.Errorf("login action failed: %s", err)
			}
			if createUserActionPicked {
				r.currentAction = createUserAction
				continue
			}
			return nil

		case createUserAction:
			userCreated, err := runCreateUserAction(r.serviceContainer)
			if err != nil {
				return fmt.Errorf("create user action failed: %s", err)
			}
			if userCreated {
				r.currentAction = loginAction
				continue
			}
			return nil
		}
	}
}

func runLoginAction(serviceContainer services.Container) (bool, error) {
	m, err := tea.NewProgram(serviceContainer.Models.Login).Run()
	if err != nil {
		return false, fmt.Errorf("could not start login action: %w", err)
	}
	if m.(cli.LoginModel).Cancelled {
		return false, fmt.Errorf("login action cancelled")
	}
	return m.(cli.LoginModel).CreateUserPicked, err
}

func runCreateUserAction(serviceContainer services.Container) (bool, error) {
	m, err := tea.NewProgram(serviceContainer.Models.CreateUser).Run()
	if err != nil {
		return false, fmt.Errorf("could not start create user action: %w", err)
	}
	if m.(cli.CreateUserModel).Cancelled {
		return false, fmt.Errorf("create user action cancelled")
	}

	err = createNewUser(serviceContainer, m)
	if err != nil {
		return false, err
	}

	return m.(cli.CreateUserModel).UserCreated, nil
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
