//go:build integration

package app

import (
	"fmt"
	"testing"
	"yubigo-pass/internal/app/crypto"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/services"
	"yubigo-pass/internal/app/utils"
	"yubigo-pass/internal/cli"
	"yubigo-pass/internal/database"
	"yubigo-pass/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
)

func TestRunnerLoginShouldSucceed(t *testing.T) {
	// given
	loginModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModel.LoggedIn = true
	mainMenuModel := cli.NewMainMenuModel()
	mainMenuModel.Choice = cli.QuitItem
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:      test.NewTeaModelMock().WillReturnOnce(loginModel),
			CreateUser: test.NewTeaModelMock(),
			MainMenu:   test.NewTeaModelMock().WillReturnOnce(mainMenuModel),
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
	assert.Equal(t, mainMenuAction, runner.currentAction)
}

func TestRunnerLoginShouldCancelExecution(t *testing.T) {
	// given
	loginModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModel.Cancelled = true
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login: test.NewTeaModelMock().WillReturnOnce(loginModel),
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

func TestRunnerLoginShouldEnterCreateUserActionAndCancelIt(t *testing.T) {
	// given
	loginModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModel.CreateUserPicked = true
	createUserModel := cli.NewCreateUserModel(test.NewStoreExecutorMock())
	createUserModel.Cancelled = true
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:      test.NewTeaModelMock().WillReturnOnce(loginModel),
			CreateUser: test.NewTeaModelMock().WillReturnOnce(createUserModel),
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
	assert.Equal(t, createUserAction, runner.currentAction)
}

func TestRunnerLoginShouldEnterCreateUserActionAndBackToLogin(t *testing.T) {
	// given
	loginModelFirst := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModelFirst.CreateUserPicked = true
	loginModelSecond := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModelSecond.LoggedIn = true
	createUserModel := cli.NewCreateUserModel(test.NewStoreExecutorMock())
	createUserModel.UserCreated = true
	mainMenuModel := cli.NewMainMenuModel()
	mainMenuModel.Choice = cli.QuitItem
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:      test.NewTeaModelMock().WillReturnOnce(loginModelFirst).WillReturnOnce(loginModelSecond),
			CreateUser: test.NewTeaModelMock().WillReturnOnce(createUserModel),
			MainMenu:   test.NewTeaModelMock().WillReturnOnce(mainMenuModel),
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
	assert.Equal(t, mainMenuAction, runner.currentAction)
}

func TestRunnerLoginShouldEnterCreateUserActionAndAbortIt(t *testing.T) {
	// given
	loginModelFirst := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModelFirst.CreateUserPicked = true
	loginModelSecond := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModelSecond.LoggedIn = true
	createUserModel := cli.NewCreateUserModel(test.NewStoreExecutorMock())
	createUserModel.UserCreationAborted = true
	mainMenuModel := cli.NewMainMenuModel()
	mainMenuModel.Choice = cli.QuitItem
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:      test.NewTeaModelMock().WillReturnOnce(loginModelFirst).WillReturnOnce(loginModelSecond),
			CreateUser: test.NewTeaModelMock().WillReturnOnce(createUserModel),
			MainMenu:   test.NewTeaModelMock().WillReturnOnce(mainMenuModel),
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
	assert.Equal(t, mainMenuAction, runner.currentAction)
}

func TestRunnerMainMenuShouldChooseGetPassword(t *testing.T) {
	// given
	loginModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModel.LoggedIn = true
	mainMenuModel := cli.NewMainMenuModel()
	mainMenuModel.Choice = cli.GetPasswordItem
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:      test.NewTeaModelMock().WillReturnOnce(loginModel),
			CreateUser: test.NewTeaModelMock(),
			MainMenu:   test.NewTeaModelMock().WillReturnOnce(mainMenuModel),
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
	assert.Equal(t, getPasswordAction, runner.currentAction)
}

func TestRunnerMainMenuShouldChooseViewPasswords(t *testing.T) {
	// given
	loginModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModel.LoggedIn = true
	mainMenuModel := cli.NewMainMenuModel()
	mainMenuModel.Choice = cli.ViewPasswordItem
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:      test.NewTeaModelMock().WillReturnOnce(loginModel),
			CreateUser: test.NewTeaModelMock(),
			MainMenu:   test.NewTeaModelMock().WillReturnOnce(mainMenuModel),
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
	assert.Equal(t, viewPasswordsAction, runner.currentAction)
}

func TestRunnerMainMenuShouldChooseAddNewPassword(t *testing.T) {
	// given
	loginModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModel.LoggedIn = true
	mainMenuModel := cli.NewMainMenuModel()
	mainMenuModel.Choice = cli.AddPasswordItem
	addPasswordModel := cli.NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession())
	addPasswordModel.Cancelled = true
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:            test.NewTeaModelMock().WillReturnOnce(loginModel),
			CreateUser:       test.NewTeaModelMock(),
			MainMenu:         test.NewTeaModelMock().WillReturnOnce(mainMenuModel),
			AddPasswordModel: test.NewTeaModelMock().WillReturnOnce(addPasswordModel),
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
	assert.Equal(t, addPasswordAction, runner.currentAction)
}

func TestRunnerMainMenuShouldChooseLogout(t *testing.T) {
	// given
	returnedLoginModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	returnedLoginModel.LoggedIn = true
	returnedLoginModelSecond := cli.NewLoginModel(test.NewStoreExecutorMock())
	returnedMainMenuModel := cli.NewMainMenuModel()
	returnedMainMenuModel.Choice = cli.LogoutItem
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:      test.NewTeaModelMock().WillReturnOnce(returnedLoginModel).WillReturnOnce(returnedLoginModelSecond),
			CreateUser: test.NewTeaModelMock(),
			MainMenu:   test.NewTeaModelMock().WillReturnOnce(returnedMainMenuModel),
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
		UserID:   test.RandomString(),
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

func TestAddPasswordShouldEncryptAndAddNewPasswordToDBAndThenDecryptIt(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)

	tm := teatest.NewTestModel(
		t,
		cli.NewAddPasswordModel(database.NewStore(db), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)

	userID, userPassword, userSalt := test.RandomString(), test.RandomString(), test.RandomString()
	session := utils.NewSession(userID, userPassword, userSalt)

	// given
	title, username, password, url := test.RandomString(), test.RandomString(), test.RandomString(), test.RandomString()
	addNewPasswordInModel(tm, title, username, password, url)

	fm := tm.FinalModel(t)
	m, ok := fm.(cli.AddPasswordModel)
	assert.True(t, ok)

	container := services.Container{
		Store: database.NewStore(db),
		Models: services.TeaModels{
			AddPasswordModel: fm,
		},
	}

	// when
	err = addNewPassword(session, container, fm)

	// then
	assert.NoError(t, err)
	assert.True(t, m.PasswordAdded)

	insertedPassword, err := container.Store.GetPassword(userID, title, username)
	assert.NoError(t, err)

	assert.Equal(t, title, insertedPassword.Title)
	assert.Equal(t, username, insertedPassword.Username)
	assert.Equal(t, url, insertedPassword.Url)
	assert.NotEmpty(t, insertedPassword.Nonce)
	assert.NotEmpty(t, insertedPassword.Password)
	assert.NotEqual(t, password, insertedPassword.Password)

	decryptedPassword, err := crypto.DecryptAES(crypto.DeriveAESKey(userPassword, userSalt), []byte(insertedPassword.Password))
	assert.NoError(t, err)
	assert.Equal(t, password, string(decryptedPassword))
}

func TestAddPasswordShouldNotAddNewPasswordWhenTheSameAlreadyExists(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)

	tm := teatest.NewTestModel(
		t,
		cli.NewAddPasswordModel(database.NewStore(db), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)

	userID, userPassword, userSalt := test.RandomString(), test.RandomString(), test.RandomString()
	session := utils.NewSession(userID, userPassword, userSalt)

	// given
	title, username, password, url := test.RandomString(), test.RandomString(), test.RandomString(), test.RandomString()
	addNewPasswordInModel(tm, title, username, password, url)

	fm := tm.FinalModel(t)
	m, ok := fm.(cli.AddPasswordModel)
	assert.True(t, ok)

	container := services.Container{
		Store: database.NewStore(db),
		Models: services.TeaModels{
			AddPasswordModel: fm,
		},
	}

	// when
	err = addNewPassword(session, container, fm)

	//then
	assert.NoError(t, err)
	assert.True(t, m.PasswordAdded)

	//and add the same password again
	tm2 := teatest.NewTestModel(
		t,
		cli.NewAddPasswordModel(database.NewStore(db), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)
	addNewPasswordInModel(tm2, title, username, password, url)

	fm2 := tm2.FinalModel(t)
	_, ok = fm.(cli.AddPasswordModel)
	assert.True(t, ok)

	container2 := services.Container{
		Store: database.NewStore(db),
		Models: services.TeaModels{
			AddPasswordModel: fm2,
		},
	}

	// when
	err = addNewPassword(session, container2, fm2)

	//expected
	expectedError := model.NewPasswordAlreadyExistsError(userID, title, username)

	// then
	assert.EqualError(t, err, expectedError.Error())
}

func TestRunnerShouldReturnToMainMenuAfterAddingPassword(t *testing.T) {
	// given
	loginModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModel.LoggedIn = true
	loginModelSecond := cli.NewLoginModel(test.NewStoreExecutorMock())
	mainMenuModel := cli.NewMainMenuModel()
	mainMenuModel.Choice = cli.AddPasswordItem
	addPasswordModel := cli.NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession())
	addPasswordModel.PasswordAdded = true
	mainMenuModelSecond := cli.NewMainMenuModel()
	mainMenuModelSecond.Choice = cli.QuitItem
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:            test.NewTeaModelMock().WillReturnOnce(loginModel).WillReturnOnce(loginModelSecond),
			CreateUser:       test.NewTeaModelMock(),
			MainMenu:         test.NewTeaModelMock().WillReturnOnce(mainMenuModel).WillReturnOnce(mainMenuModelSecond),
			AddPasswordModel: test.NewTeaModelMock().WillReturnOnce(addPasswordModel),
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
	assert.Equal(t, mainMenuAction, runner.currentAction)
}

func TestRunnerShouldReturnToMainMenuAfterGoingBack(t *testing.T) {
	// given
	loginModel := cli.NewLoginModel(test.NewStoreExecutorMock())
	loginModel.LoggedIn = true
	loginModelSecond := cli.NewLoginModel(test.NewStoreExecutorMock())
	mainMenuModel := cli.NewMainMenuModel()
	mainMenuModel.Choice = cli.AddPasswordItem
	addPasswordModel := cli.NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession())
	addPasswordModel.Back = true
	mainMenuModelSecond := cli.NewMainMenuModel()
	mainMenuModelSecond.Choice = cli.QuitItem
	serviceContainer := services.Container{
		Store: test.NewStoreExecutorMock(),
		Models: services.TeaModels{
			Login:            test.NewTeaModelMock().WillReturnOnce(loginModel).WillReturnOnce(loginModelSecond),
			CreateUser:       test.NewTeaModelMock(),
			MainMenu:         test.NewTeaModelMock().WillReturnOnce(mainMenuModel).WillReturnOnce(mainMenuModelSecond),
			AddPasswordModel: test.NewTeaModelMock().WillReturnOnce(addPasswordModel),
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
	assert.Equal(t, mainMenuAction, runner.currentAction)
}

func addNewPasswordInModel(tm *teatest.TestModel, title string, username string, password string, url string) {
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(title),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyDown,
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(username),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyDown,
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(password),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyDown,
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(url),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})
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
