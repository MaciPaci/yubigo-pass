package app

import (
	"fmt"
	"os"
	"yubigo-pass/internal/app/crypto"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/services"
	"yubigo-pass/internal/cli"
	"yubigo-pass/internal/database"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Run runs the application
func Run() {
	serviceContainer := services.Build()
	defer database.CloseDB()

	err := runCreateUserFlow(serviceContainer)
	if err != nil {
		logrus.Errorf("create user flow failed: %s:\n", err)
		os.Exit(1)
	}
}

func runCreateUserFlow(serviceContainer services.Container) error {
	m, err := serviceContainer.Programs.CreateUserProgram.Run()
	if err != nil {
		return fmt.Errorf("could not start program: %w\n", err)
	}

	userUUID := uuid.New().String()
	username, password := cli.ExtractDataFromModel(m)
	passwordHash, salt, err := crypto.HashPasswordWithSalt(password)
	if err != nil {
		return fmt.Errorf("could not hash password: %w\n", err)
	}

	createUserInput := model.NewUser(userUUID, username, passwordHash, salt)
	err = serviceContainer.Store.CreateUser(createUserInput)
	if err != nil {
		return fmt.Errorf("could not insert new user: %w\n", err)
	}

	return nil
}
