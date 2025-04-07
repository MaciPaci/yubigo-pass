//go:build e2e

package cli

import (
	"testing"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/utils"
	"yubigo-pass/internal/database"
	"yubigo-pass/test"

	"github.com/google/uuid"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldQuitAddPasswordAction(t *testing.T) {
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
					NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession()),
					teatest.WithInitialTermSize(300, 100),
				)
				test.PressKey(tm, testCase.key)

				err := tm.Quit()
				require.NoError(t, err, "Failed to quit the model")
				fm := tm.FinalModel(t)
				m, ok := fm.(AddPasswordModel)
				require.Truef(t, ok, "final model has wrong type: %T", fm)
				assert.NoError(t, m.err, "Error should be nil on quit")
			},
		)
	}
}

func TestShouldAddPassword(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	userID := uuid.New().String()
	session := utils.NewSession(userID, test.RandomString(), test.RandomString())
	store := database.NewStore(db)

	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(store, session),
		teatest.WithInitialTermSize(300, 100),
	)

	exampleTitle := test.RandomString()
	exampleUsername := test.RandomString()
	examplePassword := test.RandomString()
	exampleUrl := test.RandomString()

	// when
	test.TypeString(tm, exampleTitle)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUrl)
	test.PressKey(tm, tea.KeyDown)  // Focus Add button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)

	result := ExtractPasswordDataFromModel(m)
	assert.Equal(t, exampleUsername, result.Username)
	assert.Equal(t, exampleTitle, result.Title)
	assert.Equal(t, examplePassword, result.Password)
	assert.Equal(t, exampleUrl, result.Url)

	assert.NoError(t, m.err, "Error should be nil on quit")
}

func TestShouldNotAddPasswordWithEmptyTitle(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)
	exampleUsername := test.RandomString()
	examplePassword := test.RandomString()

	// when
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)  // URL
	test.PressKey(tm, tea.KeyDown)  // Add button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	expectedErrorMsg := "ERROR: title, username, and password fields cannot be empty"

	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestShouldNotAddPasswordWithEmptyUsername(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)
	exampleTitle := test.RandomString()
	examplePassword := test.RandomString()

	// when
	test.TypeString(tm, exampleTitle)
	test.PressKey(tm, tea.KeyDown) // Username (empty)
	test.PressKey(tm, tea.KeyDown) // Password
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)  // URL
	test.PressKey(tm, tea.KeyDown)  // Add button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	expectedErrorMsg := "ERROR: title, username, and password fields cannot be empty"

	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestShouldNotAddPasswordWithEmptyPassword(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)
	exampleTitle := test.RandomString()
	exampleUsername := test.RandomString()

	// when
	test.TypeString(tm, exampleTitle)
	test.PressKey(tm, tea.KeyDown) // Username
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown)  // Password (empty)
	test.PressKey(tm, tea.KeyDown)  // URL
	test.PressKey(tm, tea.KeyDown)  // Add button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	expectedErrorMsg := "ERROR: title, username, and password fields cannot be empty"

	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestShouldNotAddPasswordWithEmptyFields(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyDown)  // Username
	test.PressKey(tm, tea.KeyDown)  // Password
	test.PressKey(tm, tea.KeyDown)  // URL
	test.PressKey(tm, tea.KeyDown)  // Add button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	expectedErrorMsg := "ERROR: title, username, and password fields cannot be empty"

	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestShouldAddPasswordWithEmptyUrl(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	userID := uuid.New().String()
	session := utils.NewSession(userID, test.RandomString(), test.RandomString())
	store := database.NewStore(db)

	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(store, session),
		teatest.WithInitialTermSize(300, 100),
	)

	// expected
	exampleTitle := test.RandomString()
	exampleUsername := test.RandomString()
	examplePassword := test.RandomString()

	// when
	test.TypeString(tm, exampleTitle)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)  // Focus URL (empty)
	test.PressKey(tm, tea.KeyDown)  // Focus Add button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)

	result := ExtractPasswordDataFromModel(m)
	assert.Equal(t, exampleUsername, result.Username)
	assert.Equal(t, exampleTitle, result.Title)
	assert.Equal(t, examplePassword, result.Password)
	assert.Equal(t, "", result.Url)

	require.NoError(t, err, "Error should be nil on quit")
}

func TestShouldNotAddPasswordIfItAlreadyExists(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	userID := uuid.New().String()
	session := utils.NewSession(userID, test.RandomString(), test.RandomString())
	store := database.NewStore(db)

	exampleTitle := test.RandomString()
	exampleUsername := test.RandomString()
	examplePassword := test.RandomString()
	exampleExistingPassword := "existing_password_ciphertext" // Placeholder
	exampleExistingNonce := []byte("nonce")                   // Placeholder

	// and password already exists in database
	test.InsertIntoPasswords(t, db, model.NewPassword(
		userID,
		exampleTitle,
		exampleUsername,
		exampleExistingPassword,
		test.RandomString(),
		exampleExistingNonce),
	)

	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(store, session),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.TypeString(tm, exampleTitle)    // Same title
	test.PressKey(tm, tea.KeyDown)       // -> Username
	test.TypeString(tm, exampleUsername) // Same username
	test.PressKey(tm, tea.KeyDown)       // -> Password
	test.TypeString(tm, examplePassword) // New password value
	test.PressKey(tm, tea.KeyDown)       // -> URL
	test.PressKey(tm, tea.KeyDown)       // -> Add button
	test.PressKey(tm, tea.KeyEnter)      // Submit

	// then
	expectedErrorMsg := "ERROR: password entry with this title/username already exists"

	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Errorf(t, m.err, expectedErrorMsg)
}

func TestShouldAbortAddPasswordActionAfterGoingBackAndForth(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyTab)      // Focus Back button
	test.PressKey(tm, tea.KeyShiftTab) // Focus Add button
	test.PressKey(tm, tea.KeyEsc)      // Quit via Cancel (sends StateQuit)

	// then
	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.NoError(t, m.err, "Error should be nil on quit")
}

func TestShouldAbortAddPasswordActionAndGoBack(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyTab)   // Focus Back button
	test.PressKey(tm, tea.KeyEnter) // Activate Back (sends StateGoBack)

	// then
	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.NoError(t, m.err, "Error should be nil on quit")
}

func TestShouldGenerateShowAndAddPassword(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	userID := uuid.New().String()
	session := utils.NewSession(userID, test.RandomString(), test.RandomString())
	store := database.NewStore(db)

	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(store, session),
		teatest.WithInitialTermSize(300, 100),
	)

	exampleTitle := test.RandomString()
	exampleUsername := test.RandomString()
	exampleUrl := test.RandomString()

	// when
	test.TypeString(tm, exampleTitle)
	test.PressKey(tm, tea.KeyEnter)
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyCtrlG) // Generate password
	test.PressKey(tm, tea.KeyCtrlS) // Show password
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUrl)
	test.PressKey(tm, tea.KeyDown)  // Focus Add button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)

	result := ExtractPasswordDataFromModel(m)
	assert.Equal(t, exampleUsername, result.Username)
	assert.Equal(t, exampleTitle, result.Title)
	assert.NotEmpty(t, result.Password)
	assert.Equal(t, exampleUrl, result.Url)
	assert.True(t, m.passwordVisible)

	assert.NoError(t, m.err, "Error should be nil on quit")
}

func TestShouldGenerateShowThenHideAndAddPassword(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	userID := uuid.New().String()
	session := utils.NewSession(userID, test.RandomString(), test.RandomString())
	store := database.NewStore(db)

	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(store, session),
		teatest.WithInitialTermSize(300, 100),
	)

	exampleTitle := test.RandomString()
	exampleUsername := test.RandomString()
	exampleUrl := test.RandomString()

	// when
	test.TypeString(tm, exampleTitle)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyCtrlG) // Generate password
	test.PressKey(tm, tea.KeyCtrlS) // Show password
	test.PressKey(tm, tea.KeyCtrlS) // Hide password
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUrl)
	test.PressKey(tm, tea.KeyDown)  // Focus Add button
	test.PressKey(tm, tea.KeyEnter) // Submit

	// then
	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)

	result := ExtractPasswordDataFromModel(m)
	assert.Equal(t, exampleUsername, result.Username)
	assert.Equal(t, exampleTitle, result.Title)
	assert.NotEmpty(t, result.Password)
	assert.Equal(t, exampleUrl, result.Url)
	assert.False(t, m.passwordVisible)

	assert.NoError(t, m.err, "Error should be nil on quit")
}
