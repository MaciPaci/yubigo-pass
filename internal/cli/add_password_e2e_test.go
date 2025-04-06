//go:build e2e

package cli

import (
	"bytes"
	"io"
	"testing"
	"time"
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

				// Wait a short moment for the quit message to process
				tm.WaitFinished(t, teatest.WithFinalTimeout(time.Millisecond*200))

				fm := tm.FinalModel(t)
				m, ok := fm.(AddPasswordModel)
				require.Truef(t, ok, "final model has wrong type: %T", fm)
				assert.Truef(t, m.Cancelled, "final model is not Cancelled")
				tm.WaitFinished(t, teatest.WithFinalTimeout(time.Millisecond*100))
			},
		)
	}
}

func TestShouldAddPassword(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db) // Ensure teardown

	userID := uuid.New().String()
	session := utils.NewSession(userID, test.RandomString(), test.RandomString())

	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(database.NewStore(db), session),
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
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)

	assert.True(t, m.PasswordAdded, "PasswordAdded flag should be true")
	assert.False(t, m.Cancelled, "Cancelled flag should be false")
	assert.False(t, m.Back, "Back flag should be false")

	result := ExtractPasswordDataFromModel(m)
	assert.Equal(t, exampleUsername, result.Username)
	assert.Equal(t, exampleTitle, result.Title)
	assert.Equal(t, examplePassword, result.Password)
	assert.Equal(t, exampleUrl, result.Url)

	// Check final view output - less critical than flags but can be useful
	out, err := io.ReadAll(tm.FinalOutput(t))
	require.NoError(t, err)
	// The view might show "Saving..." or similar based on PasswordAdded flag
	assert.Contains(t, string(out), "Password validation OK")
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
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte(expectedErrorMsg))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*2),
	)
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
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte(expectedErrorMsg))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*2),
	)
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
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte(expectedErrorMsg))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*2),
	)
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
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte(expectedErrorMsg))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*2),
	)
}

func TestShouldAddPasswordWithEmptyUrl(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	userID := uuid.New().String()
	session := utils.NewSession(userID, test.RandomString(), test.RandomString())

	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(database.NewStore(db), session),
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
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)

	assert.True(t, m.PasswordAdded)
	assert.False(t, m.Cancelled)
	assert.False(t, m.Back)

	result := ExtractPasswordDataFromModel(m)
	assert.Equal(t, exampleUsername, result.Username)
	assert.Equal(t, exampleTitle, result.Title)
	assert.Equal(t, examplePassword, result.Password)
	assert.Equal(t, "", result.Url)

	out, err := io.ReadAll(tm.FinalOutput(t))
	require.NoError(t, err)
	assert.Contains(t, string(out), "Password validation OK")
}

func TestShouldNotAddPasswordIfItAlreadyExists(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed to set up test database")
	defer test.TeardownTestDB(db)

	userID := uuid.New().String()
	session := utils.NewSession(userID, test.RandomString(), test.RandomString())

	store := database.NewStore(db) // Use the real store

	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(store, session),
		teatest.WithInitialTermSize(300, 100),
	)

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
	require.NoError(t, err, "Failed to insert existing password for test")

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
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte(expectedErrorMsg))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*2),
	)
}

func TestShouldAbortAddPasswordActionAfterGoingBackAndForth(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyTab) // Focus Back button
	test.PressKey(tm, tea.KeyTab) // Focus input again
	test.PressKey(tm, tea.KeyEsc) // Quit via Cancelled

	// then
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.True(t, m.Cancelled)
	assert.False(t, m.Back)
	assert.False(t, m.PasswordAdded)
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
	test.PressKey(tm, tea.KeyEnter) // Quit via Back

	// then
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.False(t, m.Cancelled)
	assert.True(t, m.Back)
	assert.False(t, m.PasswordAdded)
}
