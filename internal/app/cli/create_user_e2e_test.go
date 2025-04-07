//go:build e2e

package cli

import (
	"testing"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/database"
	"yubigo-pass/test"

	"github.com/google/uuid"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldQuitCreateUserAction(t *testing.T) {
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
					NewCreateUserModel(test.NewStoreExecutorMock()),
					teatest.WithInitialTermSize(300, 100),
				)
				test.PressKey(tm, testCase.key)

				err := tm.Quit()
				require.NoError(t, err, "Failed to quit the model")
				fm := tm.FinalModel(t)
				m, ok := fm.(CreateUserModel)
				require.Truef(t, ok, "final model has wrong type: %T", fm)
				assert.NoError(t, m.err, "Error should be nil on quit")
			},
		)
	}
}

func TestShouldCreateUser(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	store := database.NewStore(db)
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(store),
		teatest.WithInitialTermSize(300, 100),
	)

	// expected
	exampleUsername := test.RandomString()
	examplePassword := test.RandomString()

	// when
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)  // -> Submit Button
	test.PressKey(tm, tea.KeyEnter) // Submit (sends CreateUserCmd)

	// then
	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(CreateUserModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)

	// assert state before finishing
	assert.Equal(t, exampleUsername, m.inputs[0].Value())
	assert.NoError(t, m.err, "Error should be nil on successful submission")
}

func TestShouldShowPasswordDuringCreateUser(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	store := database.NewStore(db)
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(store),
		teatest.WithInitialTermSize(300, 100),
	)

	// expected
	exampleUsername := test.RandomString()
	examplePassword := test.RandomString()

	// when
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyCtrlS) // Show password
	test.PressKey(tm, tea.KeyDown)  // -> Submit Button
	test.PressKey(tm, tea.KeyEnter) // Submit (sends CreateUserCmd)

	// then
	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(CreateUserModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)

	// assert state before finishing
	assert.Equal(t, exampleUsername, m.inputs[0].Value())
	assert.True(t, m.passwordVisible, "Password should be visible")
	assert.NoError(t, m.err, "Error should be nil on successful submission")
}

func TestShouldShowAndHidePasswordDuringCreateUser(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	store := database.NewStore(db)
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(store),
		teatest.WithInitialTermSize(300, 100),
	)

	// expected
	exampleUsername := test.RandomString()
	examplePassword := test.RandomString()

	// when
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyCtrlS) // Show password
	test.PressKey(tm, tea.KeyCtrlS) // Hide password
	test.PressKey(tm, tea.KeyDown)  // -> Submit Button
	test.PressKey(tm, tea.KeyEnter) // Submit (sends CreateUserCmd)

	// then
	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(CreateUserModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)

	// assert state before finishing
	assert.Equal(t, exampleUsername, m.inputs[0].Value())
	assert.False(t, m.passwordVisible, "Password should be visible")
	assert.NoError(t, m.err, "Error should be nil on successful submission")
}

func TestShouldNotCreateUserWithEmptyUsername(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)
	examplePassword := test.RandomString()

	// when
	test.PressKey(tm, tea.KeyEnter) // -> Password
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)  // -> Submit Button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	expectedErrorMsg := "ERROR: username and password cannot be empty"

	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(CreateUserModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestShouldNotCreateUserWithEmptyPassword(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)
	exampleUsername := test.RandomString()

	// when
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyUp)    // -> Submit Button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	expectedErrorMsg := "ERROR: username and password cannot be empty"

	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(CreateUserModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestShouldNotCreateUserWithBothInputFieldsEmpty(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyDown)  // -> Password
	test.PressKey(tm, tea.KeyDown)  // -> Submit button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	expectedErrorMsg := "ERROR: username and password cannot be empty"

	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(CreateUserModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestShouldNotCreateUserIfUsernameAlreadyExists(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	store := database.NewStore(db)

	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(store), // Use real store
		teatest.WithInitialTermSize(300, 100),
	)

	existingUsername := test.RandomString()
	examplePassword := test.RandomString()

	// and username already exists in database
	test.InsertIntoUsers(t, db, model.NewUser(uuid.New().String(), existingUsername, test.RandomString(), test.RandomString()))

	// when
	test.TypeString(tm, existingUsername)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	expectedErrorMsg := "ERROR: username '" + existingUsername + "' already exists"

	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(CreateUserModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestShouldAbortCreateUserActionUsingBack(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyTab)   // -> Back button
	test.PressKey(tm, tea.KeyTab)   // -> Username
	test.PressKey(tm, tea.KeyTab)   // -> Back button
	test.PressKey(tm, tea.KeyEnter) // Activate Back (sends StateGoBack)

	// then
	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(CreateUserModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.NoError(t, m.err, "Error should be nil on quit")
}
