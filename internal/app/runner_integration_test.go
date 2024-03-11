//go:build integration

package app

import (
	"fmt"
	"testing"
	"yubigo-pass/internal/app/crypto"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/services"
	"yubigo-pass/internal/cli"
	"yubigo-pass/internal/database"
	"yubigo-pass/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
)

func TestRunnerLoginShouldSucceed(t *testing.T) {
	// given
	returnedModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	returnedModel.LoggedIn = true
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:      test.NewTeaModelMock().WillReturnOnce(returnedModel),
			CreateUser: test.NewTeaModelMock(),
		},
	}
	runner := NewRunner(serviceContainer)

	// when
	err := runner.Run()

	// this is a workaround so that gotest to correctly catch PASS message
	// it is awful and I hate it, but it works
	fmt.Println('\n')

	// then
	assert.NoError(t, err)
	assert.Equal(t, loginAction, runner.currentAction)
}

func TestRunnerLoginShouldCancelExecution(t *testing.T) {
	// given
	returnedModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	returnedModel.Cancelled = true
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:      test.NewTeaModelMock().WillReturnOnce(returnedModel),
			CreateUser: test.NewTeaModelMock(),
		},
	}
	runner := NewRunner(serviceContainer)

	// expected
	expectedError := fmt.Errorf("login action failed: login action cancelled")

	// when
	err := runner.Run()

	// this is a workaround so that gotest to correctly catch PASS message
	// it is awful and I hate it, but it works
	fmt.Println('\n')

	// then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, loginAction, runner.currentAction)
}

func TestRunnerLoginShouldEnterCreateUserActionAndCancelIt(t *testing.T) {
	// given
	returnedLoginModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	returnedLoginModel.CreateUserPicked = true
	returnedCreateUserModel := cli.NewCreateUserModel(test.NewStoreExecutorMock())
	returnedCreateUserModel.Cancelled = true
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:      test.NewTeaModelMock().WillReturnOnce(returnedLoginModel),
			CreateUser: test.NewTeaModelMock().WillReturnOnce(returnedCreateUserModel),
		},
	}
	runner := NewRunner(serviceContainer)

	// expected
	expectedError := fmt.Errorf("create user action failed: create user action cancelled")

	// when
	err := runner.Run()

	// this is a workaround so that gotest to correctly catch PASS message
	// it is awful and I hate it, but it works
	fmt.Println('\n')

	// then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, createUserAction, runner.currentAction)
}

func TestRunnerLoginShouldEnterCreateUserActionAndBackToLogin(t *testing.T) {
	// given
	returnedLoginModelFirst := cli.NewLoginModel(test.NewStoreExecutorMock())
	returnedLoginModelFirst.CreateUserPicked = true
	returnedLoginModelSecond := cli.NewLoginModel(test.NewStoreExecutorMock())
	returnedLoginModelSecond.LoggedIn = true
	returnedCreateUserModel := cli.NewCreateUserModel(test.NewStoreExecutorMock())
	returnedCreateUserModel.UserCreated = true
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:      test.NewTeaModelMock().WillReturnOnce(returnedLoginModelFirst).WillReturnOnce(returnedLoginModelSecond),
			CreateUser: test.NewTeaModelMock().WillReturnOnce(returnedCreateUserModel),
		},
	}
	runner := NewRunner(serviceContainer)

	// when
	err := runner.Run()

	// this is a workaround so that gotest to correctly catch PASS message
	// it is awful and I hate it, but it works
	fmt.Println('\n')

	// then
	assert.NoError(t, err)
	assert.Equal(t, loginAction, runner.currentAction)
}

func TestCreateUserFlowShouldCreateNewUserInDB(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)

	tm := teatest.NewTestModel(
		t,
		cli.NewCreateUserModel(database.NewStore(db)),
		teatest.WithInitialTermSize(300, 100),
	)

	// given
	createUserWithTestModel(tm)
	fm := tm.FinalModel(t)
	username, password := cli.ExtractDataFromModel(fm)

	container := services.Container{
		Store: database.NewStore(db),
		Models: services.TeaModels{
			CreateUser: fm,
		},
	}

	// when
	err = createNewUser(container, fm)
	assert.Nil(t, err)

	// then
	createdUser, err := container.Store.GetUser(username)
	assert.Nil(t, err)

	hashedPassword := crypto.HashPasswordWithSalt(password, createdUser.Salt)
	assert.Equal(t, username, createdUser.Username)
	assert.Equal(t, hashedPassword, createdUser.Password)
}

func TestCreateUserFlowShouldNotCreateNewUserInDBIfOneWithTheSameUsernameAlreadyExists(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)

	tm := teatest.NewTestModel(
		t,
		cli.NewCreateUserModel(database.NewStore(db)),
		teatest.WithInitialTermSize(300, 100),
	)

	// given
	createUserWithTestModel(tm)
	fm := tm.FinalModel(t)
	username, _ := cli.ExtractDataFromModel(fm)

	container := services.Container{
		Store: database.NewStore(db),
		Models: services.TeaModels{
			CreateUser: fm,
		},
	}

	// and user with the same username already exists in the database
	test.InsertIntoUsers(t, db, model.User{
		Uuid:     test.RandomString(),
		Username: username,
		Password: test.RandomString(),
		Salt:     test.RandomString(),
	})

	// expected
	expectedError := model.NewUserAlreadyExistsError(username)

	// when
	err = createNewUser(container, fm)

	// then
	assert.EqualError(t, err, expectedError.Error())
}

func createUserWithTestModel(tm *teatest.TestModel) {
	exampleUsername := test.RandomString()
	examplePassword := test.RandomString()

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(exampleUsername),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyDown,
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(examplePassword),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyTab,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})
}
