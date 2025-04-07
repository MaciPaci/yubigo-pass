//go:build e2e

package cli

import (
	"bytes"
	"testing"
	"time"
	"yubigo-pass/internal/app/crypto"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/app/services"
	"yubigo-pass/internal/database"
	"yubigo-pass/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppModel_LoginSuccessFlow(t *testing.T) {
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed setup")
	defer test.TeardownTestDB(db)
	store := database.NewStore(db)
	container := services.Container{Store: store}

	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	require.NoError(t, err)
	existingUser := model.User{
		UserID:   uuid.New().String(),
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}
	test.InsertIntoUsers(t, db, existingUser)

	tm := teatest.NewTestModel(t, NewAppModel(container), teatest.WithInitialTermSize(300, 100))

	// Wait for Login screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("LOGIN"))
	}, teatest.WithDuration(2*time.Second))

	test.TypeString(tm, existingUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, existingPassword)
	test.PressKey(tm, tea.KeyDown)  // -> Login Button
	test.PressKey(tm, tea.KeyEnter) // Submit Login

	// Wait for Main Menu screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("MAIN MENU")) && // Check for Main Menu title
			!bytes.Contains(bts, []byte("LOGIN")) // Ensure Login view is gone
	}, teatest.WithDuration(2*time.Second))

	tm.Quit() // Manually quit as the app is now in main menu
}

func TestAppModel_CreateUserFlow_Success(t *testing.T) {
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed setup")
	defer test.TeardownTestDB(db)
	store := database.NewStore(db)
	container := services.Container{Store: store}

	tm := teatest.NewTestModel(t, NewAppModel(container), teatest.WithInitialTermSize(300, 100))

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("LOGIN"))
	}, teatest.WithDuration(2*time.Second))

	// Navigate to Create User button
	test.PressKey(tm, tea.KeyTab)   // -> Create User Btn
	test.PressKey(tm, tea.KeyEnter) // Activate Create User

	// Wait for Create User screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("CREATE NEW USER")) &&
			!bytes.Contains(bts, []byte("LOGIN"))
	}, teatest.WithDuration(2*time.Second))

	newUsername := test.RandomString()
	newPassword := test.RandomString()

	test.TypeString(tm, newUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, newPassword)
	test.PressKey(tm, tea.KeyDown)  // -> Submit Button
	test.PressKey(tm, tea.KeyEnter) // Submit Create User

	// Wait to return to Login screen (AppModel logic after StateUserCreated)
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("LOGIN")) &&
			!bytes.Contains(bts, []byte("CREATE NEW USER"))
	}, teatest.WithDuration(3*time.Second))

	// Verify user exists in DB
	_, dbErr := store.GetUser(newUsername)
	assert.NoError(t, dbErr, "User should exist in database after creation")

	tm.Quit()
}

func TestAppModel_CreateUserFlow_UserAlreadyExists(t *testing.T) {
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed setup")
	defer test.TeardownTestDB(db)
	store := database.NewStore(db)
	container := services.Container{Store: store}

	tm := teatest.NewTestModel(t, NewAppModel(container), teatest.WithInitialTermSize(300, 100))

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("LOGIN"))
	}, teatest.WithDuration(2*time.Second))

	// Navigate to Create User button
	test.PressKey(tm, tea.KeyTab)   // -> Create User Btn
	test.PressKey(tm, tea.KeyEnter) // Activate Create User

	// Wait for Create User screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("CREATE NEW USER")) &&
			!bytes.Contains(bts, []byte("LOGIN"))
	}, teatest.WithDuration(2*time.Second))

	newUsername := test.RandomString()
	newPassword := test.RandomString()

	test.InsertIntoUsers(t, db, model.User{UserID: test.RandomString(), Username: newUsername, Password: newPassword, Salt: test.RandomString()})

	test.TypeString(tm, newUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, newPassword)
	test.PressKey(tm, tea.KeyDown)  // -> Submit Button
	test.PressKey(tm, tea.KeyEnter) // Submit Create User

	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
}

func TestAppModel_CreateUserFlow_GoBack(t *testing.T) {
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed setup")
	defer test.TeardownTestDB(db)
	store := database.NewStore(db)
	container := services.Container{Store: store}

	tm := teatest.NewTestModel(t, NewAppModel(container), teatest.WithInitialTermSize(300, 100))

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("LOGIN"))
	}, teatest.WithDuration(2*time.Second))

	// Navigate to Create User button
	test.PressKey(tm, tea.KeyTab)   // -> Create User Btn
	test.PressKey(tm, tea.KeyEnter) // Activate Create User

	// Wait for Create User screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("CREATE NEW USER")) &&
			!bytes.Contains(bts, []byte("LOGIN"))
	}, teatest.WithDuration(2*time.Second))

	test.PressKey(tm, tea.KeyTab)   // -> Back Button
	test.PressKey(tm, tea.KeyEnter) // Go Back

	// Wait to return to Login screen (AppModel logic after StateUserCreated)
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("LOGIN")) &&
			!bytes.Contains(bts, []byte("CREATE NEW USER"))
	}, teatest.WithDuration(3*time.Second))

	tm.Quit()
}

func TestAppModel_AddPasswordFlow_Success(t *testing.T) {
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed setup")
	defer test.TeardownTestDB(db)
	store := database.NewStore(db)
	container := services.Container{Store: store}

	// Setup existing user for login
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	require.NoError(t, err)
	userID := uuid.New().String()
	existingUser := model.User{
		UserID:   userID,
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}
	test.InsertIntoUsers(t, db, existingUser)

	tm := teatest.NewTestModel(t, NewAppModel(container), teatest.WithInitialTermSize(300, 100))

	// Login first
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool { return bytes.Contains(bts, []byte("LOGIN")) })
	test.TypeString(tm, existingUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, existingPassword)
	test.PressKey(tm, tea.KeyDown)  // -> Login Button
	test.PressKey(tm, tea.KeyEnter) // Submit Login

	// Wait for Main Menu
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool { return bytes.Contains(bts, []byte("MAIN MENU")) })

	// Navigate to Add Password and select
	test.PressKey(tm, tea.KeyDown)  // -> View Passwords
	test.PressKey(tm, tea.KeyDown)  // -> Add Password
	test.PressKey(tm, tea.KeyEnter) // Select Add Password

	// Wait for Add Password screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("ADD A NEW PASSWORD")) &&
			!bytes.Contains(bts, []byte("MAIN MENU"))
	}, teatest.WithDuration(2*time.Second))

	newTitle := test.RandomString()
	newPwdUsername := test.RandomString()
	newPassword := test.RandomString()

	// Fill Add Password form
	test.TypeString(tm, newTitle)
	test.PressKey(tm, tea.KeyDown) // -> Username
	test.TypeString(tm, newPwdUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, newPassword)
	test.PressKey(tm, tea.KeyDown)  // -> URL
	test.PressKey(tm, tea.KeyDown)  // -> Add Button
	test.PressKey(tm, tea.KeyEnter) // Submit Add Password

	// Wait to return to Main Menu screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("MAIN MENU")) &&
			!bytes.Contains(bts, []byte("ADD A NEW PASSWORD"))
	}, teatest.WithDuration(3*time.Second))

	// Verify password exists in DB
	_, dbErr := store.GetPassword(userID, newTitle, newPwdUsername)
	assert.NoError(t, dbErr, "Password should exist in database after adding")

	tm.Quit()
}

func TestAppModel_AddPasswordFlow_PasswordAlredyExists(t *testing.T) {
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed setup")
	defer test.TeardownTestDB(db)
	store := database.NewStore(db)
	container := services.Container{Store: store}

	// Setup existing user for login
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	require.NoError(t, err)
	userID := uuid.New().String()
	existingUser := model.User{
		UserID:   userID,
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}
	test.InsertIntoUsers(t, db, existingUser)

	tm := teatest.NewTestModel(t, NewAppModel(container), teatest.WithInitialTermSize(300, 100))

	// Login first
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool { return bytes.Contains(bts, []byte("LOGIN")) })
	test.TypeString(tm, existingUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, existingPassword)
	test.PressKey(tm, tea.KeyDown)  // -> Login Button
	test.PressKey(tm, tea.KeyEnter) // Submit Login

	// Wait for Main Menu
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool { return bytes.Contains(bts, []byte("MAIN MENU")) })

	// Navigate to Add Password and select
	test.PressKey(tm, tea.KeyDown)  // -> View Passwords
	test.PressKey(tm, tea.KeyDown)  // -> Add Password
	test.PressKey(tm, tea.KeyEnter) // Select Add Password

	// Wait for Add Password screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("ADD A NEW PASSWORD")) &&
			!bytes.Contains(bts, []byte("MAIN MENU"))
	}, teatest.WithDuration(2*time.Second))

	newTitle := test.RandomString()
	newPwdUsername := test.RandomString()
	newPassword := test.RandomString()

	test.InsertIntoPasswords(t, db, model.Password{
		UserID:   userID,
		Title:    newTitle,
		Username: newPwdUsername,
		Password: newPassword,
		Url:      "",
		Nonce:    []byte{},
	})

	// Fill Add Password form
	test.TypeString(tm, newTitle)
	test.PressKey(tm, tea.KeyDown) // -> Username
	test.TypeString(tm, newPwdUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, newPassword)
	test.PressKey(tm, tea.KeyDown)  // -> URL
	test.PressKey(tm, tea.KeyDown)  // -> Add Button
	test.PressKey(tm, tea.KeyEnter) // Submit Add Password

	err = tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
}

func TestAppModel_AddPasswordFlow_GoBack(t *testing.T) {
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed setup")
	defer test.TeardownTestDB(db)
	store := database.NewStore(db)
	container := services.Container{Store: store}

	// Setup existing user for login
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	require.NoError(t, err)
	userID := uuid.New().String()
	existingUser := model.User{
		UserID:   userID,
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}
	test.InsertIntoUsers(t, db, existingUser)

	tm := teatest.NewTestModel(t, NewAppModel(container), teatest.WithInitialTermSize(300, 100))

	// Login first
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool { return bytes.Contains(bts, []byte("LOGIN")) })
	test.TypeString(tm, existingUsername)
	test.PressKey(tm, tea.KeyDown) // -> Password
	test.TypeString(tm, existingPassword)
	test.PressKey(tm, tea.KeyDown)  // -> Login Button
	test.PressKey(tm, tea.KeyEnter) // Submit Login

	// Wait for Main Menu
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool { return bytes.Contains(bts, []byte("MAIN MENU")) })

	// Navigate to Add Password and select
	test.PressKey(tm, tea.KeyDown)  // -> View Passwords
	test.PressKey(tm, tea.KeyDown)  // -> Add Password
	test.PressKey(tm, tea.KeyEnter) // Select Add Password

	// Wait for Add Password screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("ADD A NEW PASSWORD")) &&
			!bytes.Contains(bts, []byte("MAIN MENU"))
	}, teatest.WithDuration(2*time.Second))

	// Fill Add Password form
	test.PressKey(tm, tea.KeyTab)   // -> Back Button
	test.PressKey(tm, tea.KeyEnter) // Go Back

	// Wait to return to Main Menu screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("MAIN MENU")) &&
			!bytes.Contains(bts, []byte("ADD A NEW PASSWORD"))
	}, teatest.WithDuration(3*time.Second))

	tm.Quit()
}

func TestAppModel_LogoutFlow(t *testing.T) {
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed setup")
	defer test.TeardownTestDB(db)
	store := database.NewStore(db)
	container := services.Container{Store: store}

	// Setup existing user for login
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	require.NoError(t, err)
	existingUser := model.User{
		UserID:   uuid.New().String(),
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}
	test.InsertIntoUsers(t, db, existingUser)

	tm := teatest.NewTestModel(t, NewAppModel(container), teatest.WithInitialTermSize(300, 100))

	// Login
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool { return bytes.Contains(bts, []byte("LOGIN")) })
	test.TypeString(tm, existingUsername)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, existingPassword)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// Wait for Main Menu
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool { return bytes.Contains(bts, []byte("MAIN MENU")) })

	// Navigate to Logout and select
	test.PressKey(tm, tea.KeyDown)  // -> View
	test.PressKey(tm, tea.KeyDown)  // -> Add
	test.PressKey(tm, tea.KeyDown)  // -> Logout
	test.PressKey(tm, tea.KeyEnter) // Select Logout

	// Wait to return to Login screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("LOGIN")) &&
			!bytes.Contains(bts, []byte("MAIN MENU"))
	}, teatest.WithDuration(2*time.Second))

	tm.Quit()
}

func TestAppModel_QuitFlow(t *testing.T) {
	db, err := test.SetupTestDB()
	require.NoError(t, err, "Failed setup")
	defer test.TeardownTestDB(db)
	store := database.NewStore(db)
	container := services.Container{Store: store}

	tm := teatest.NewTestModel(t, NewAppModel(container), teatest.WithInitialTermSize(300, 100))

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("LOGIN"))
	}, teatest.WithDuration(2*time.Second))

	test.PressKey(tm, tea.KeyEsc)

	// Wait for program to finish
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
