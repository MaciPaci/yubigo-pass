//go:build e2e

package cli

import (
	"testing"
	"yubigo-pass/internal/app/crypto"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/database"
	"yubigo-pass/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldQuitLoginAction(t *testing.T) {
	testCases := []struct {
		name string
		key  tea.KeyType
	}{
		{
			name: "escape was pressed",
			key:  tea.KeyEsc,
		},
		{
			name: "ctrl+c was pressed",
			key:  tea.KeyCtrlC,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				tm := teatest.NewTestModel(
					t,
					NewLoginModel(test.NewStoreExecutorMock()),
					teatest.WithInitialTermSize(300, 100),
				)
				test.PressKey(tm, testCase.key)

				err := tm.Quit()
				require.NoError(t, err, "Failed to quit the model")
				fm := tm.FinalModel(t)
				m, ok := fm.(LoginModel)
				require.Truef(t, ok, "final model has wrong type: %T", fm)
				assert.NoError(t, m.err, "Error should be nil on quit")
			},
		)
	}
}

func TestLoginShouldSucceed(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	store := database.NewStore(db)
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(store),
		teatest.WithInitialTermSize(300, 100),
	)

	// expected
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	require.NoError(t, err, "could not create salt")

	existingUser := model.User{
		UserID:   uuid.New().String(),
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}

	// and user exists in database
	test.InsertIntoUsers(t, db, existingUser)

	// when
	test.TypeString(tm, existingUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, existingPassword)
	test.PressKey(tm, tea.KeyDown)  // -> Login Button
	test.PressKey(tm, tea.KeyEnter) // Submit (sends LogSuccessCmd)

	// then
	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(LoginModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Equal(t, existingUsername, m.inputs[0].Value())
	assert.Equal(t, existingPassword, m.inputs[1].Value())
	assert.NoError(t, m.err, "Error should be nil on successful login")
}

func TestLoginShouldNotSucceedWithIncorrectCredentials(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	store := database.NewStore(db)
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(store),
		teatest.WithInitialTermSize(300, 100),
	)
	randomUsername := test.RandomString()
	randomPassword := test.RandomString()

	// expected (user that actually exists)
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	require.NoError(t, err, "could not create salt")

	existingUser := model.User{
		UserID:   uuid.New().String(),
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}

	// and user exists in database
	test.InsertIntoUsers(t, db, existingUser)

	// when
	test.TypeString(tm, randomUsername) // Use incorrect username
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, randomPassword) // Use incorrect password
	test.PressKey(tm, tea.KeyDown)      // -> Login Button
	test.PressKey(tm, tea.KeyEnter)     // Submit

	// then
	expectedErrorMsg := "ERROR: incorrect username or password"

	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(LoginModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestLoginShouldNotSucceedWithIncorrectPassword(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	store := database.NewStore(db)
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(store),
		teatest.WithInitialTermSize(300, 100),
	)
	incorrectPassword := test.RandomString()

	// expected
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	require.NoError(t, err, "could not create salt")

	existingUser := model.User{
		UserID:   uuid.New().String(),
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}

	// and user exists in database
	test.InsertIntoUsers(t, db, existingUser)

	// when
	test.TypeString(tm, existingUsername) // Correct username
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, incorrectPassword) // Incorrect password
	test.PressKey(tm, tea.KeyDown)         // -> Login Button
	test.PressKey(tm, tea.KeyEnter)        // Submit

	// then
	expectedErrorMsg := "ERROR: incorrect username or password"

	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(LoginModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestLoginShouldNotValidateEmptyUsername(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)

	emptyUsername := ""
	randomPassword := test.RandomString()

	// when
	test.TypeString(tm, emptyUsername) // Still need to "type" empty string to potentially clear default values if any
	test.PressKey(tm, tea.KeyDown)     // -> Password
	test.TypeString(tm, randomPassword)
	test.PressKey(tm, tea.KeyDown)  // -> Login Button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	expectedErrorMsg := "ERROR: username and password cannot be empty"

	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(LoginModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestLoginShouldNotValidateEmptyPassword(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)

	emptyPassword := ""
	randomUsername := test.RandomString()

	// when
	test.TypeString(tm, randomUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, emptyPassword)
	test.PressKey(tm, tea.KeyDown)  // -> Login Button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	expectedErrorMsg := "ERROR: username and password cannot be empty"

	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(LoginModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestLoginShouldNotValidateEmptyUsernameAndPassword(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyDown)  // -> Password
	test.PressKey(tm, tea.KeyDown)  // -> Login Button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	expectedErrorMsg := "ERROR: username and password cannot be empty"

	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(LoginModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestLoginShouldEnterCreateUserAction(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyTab)   // -> Create User Button
	test.PressKey(tm, tea.KeyEnter) // Activate Create User (sends StateGoToCreateUser)

	// then
	// Model finishes because it sent StateGoToCreateUser command
	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(LoginModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.NoError(t, m.err, "Error should be nil on quit")
}

func TestLoginShouldSucceedAfterEnteringCreateUserFlowAndBack(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	store := database.NewStore(db)
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(store),
		teatest.WithInitialTermSize(300, 100),
	)

	// expected
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	require.NoError(t, err, "could not create salt")

	existingUser := model.User{
		UserID:   uuid.New().String(),
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}

	// and user exists in database
	test.InsertIntoUsers(t, db, existingUser)

	// when
	// Simulate typing username, tabbing to Create User, tabbing back, typing password, submitting
	test.TypeString(tm, existingUsername)
	test.PressKey(tm, tea.KeyTab)      // -> Create User Btn
	test.PressKey(tm, tea.KeyShiftTab) // <- Username
	test.PressKey(tm, tea.KeyEnter)    // <- Password
	test.TypeString(tm, existingPassword)
	test.PressKey(tm, tea.KeyDown)  // -> Login Button
	test.PressKey(tm, tea.KeyEnter) // Submit (sends LogSuccessCmd)

	// then
	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(LoginModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Equal(t, existingUsername, m.inputs[0].Value())
	assert.Equal(t, existingPassword, m.inputs[1].Value())
	assert.NoError(t, m.err, "Error should be nil on successful login")
}
