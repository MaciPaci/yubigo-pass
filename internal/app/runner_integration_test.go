//go:build integration

package app

import (
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
		Programs: struct {
			CreateUserProgram *tea.Program
		}{
			CreateUserProgram: tea.NewProgram(fm),
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
		Programs: struct {
			CreateUserProgram *tea.Program
		}{
			CreateUserProgram: tea.NewProgram(fm),
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
