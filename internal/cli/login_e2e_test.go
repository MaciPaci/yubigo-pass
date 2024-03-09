//go:build e2e

package cli

import (
	"bytes"
	"io"
	"testing"
	"time"
	"yubigo-pass/internal/app/crypto"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/database"
	"yubigo-pass/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestShouldQuitLoginAction(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)

	testCases := []struct {
		name string
		key  tea.KeyType
	}{
		{
			"escape was pressed",
			tea.KeyEsc,
		},
		{
			"ctrl+c was pressed",
			tea.KeyCtrlC,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				tm.Send(tea.KeyMsg{
					Type: testCase.key,
				})
				fm := tm.FinalModel(t)
				m, ok := fm.(LoginModel)
				assert.Truef(t, ok, "final model has wrong type: %T", fm)
				assert.Truef(t, m.WasCancelled(), "final model is not cancelled")
				tm.WaitFinished(t, teatest.WithFinalTimeout(time.Millisecond*100))
			},
		)
	}
}

func TestLoginShouldSucceed(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(database.NewStore(db)),
		teatest.WithInitialTermSize(300, 100),
	)

	// expected
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	if err != nil {
		t.Fatalf("could not create salt")
	}

	existingUser := model.User{
		Uuid:     uuid.New().String(),
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}

	// and user exists in database
	test.InsertIntoUsers(t, db, existingUser)

	// when
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(existingUsername),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyDown,
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(existingPassword),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(LoginModel)
	assert.True(t, ok)
	assert.True(t, m.LoggedInSuccessfully())

	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	assert.True(t, bytes.Contains(out, []byte("Logged in")))

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestLoginShouldNotSucceedWithIncorrectCredentials(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(database.NewStore(db)),
		teatest.WithInitialTermSize(300, 100),
	)
	randomUsername := test.RandomString()
	randomPassword := test.RandomString()

	// expected
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	if err != nil {
		t.Fatalf("could not create salt")
	}

	existingUser := model.User{
		Uuid:     uuid.New().String(),
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}

	// and user exists in database
	test.InsertIntoUsers(t, db, existingUser)

	// when
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(randomUsername),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyDown,
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(randomPassword),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	// then
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("ERROR: incorrect credentials"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*1),
	)
}

func TestLoginShouldNotSucceedWithIncorrectPassword(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(database.NewStore(db)),
		teatest.WithInitialTermSize(300, 100),
	)
	incorrectPassword := test.RandomString()

	// expected
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	if err != nil {
		t.Fatalf("could not create salt")
	}

	existingUser := model.User{
		Uuid:     uuid.New().String(),
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}

	// and user exists in database
	test.InsertIntoUsers(t, db, existingUser)

	// when
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(existingUsername),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyDown,
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(incorrectPassword),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	// then
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("ERROR: incorrect credentials"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*1),
	)
}

func TestLoginShouldNotValidateEmptyUsername(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(database.NewStore(db)),
		teatest.WithInitialTermSize(300, 100),
	)

	emptyUsername := ""
	randomPassword := test.RandomString()

	// when
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(emptyUsername),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyDown,
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(randomPassword),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	// then
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("ERROR: username cannot be empty"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*1),
	)
}

func TestLoginShouldNotValidateEmptyPassword(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(database.NewStore(db)),
		teatest.WithInitialTermSize(300, 100),
	)

	emptyPassword := ""
	randomUsername := test.RandomString()

	// when
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(randomUsername),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyDown,
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(emptyPassword),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	// then
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("ERROR: password cannot be empty"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*1),
	)
}

func TestLoginShouldNotValidateEmptyUsernameAndPassword(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(database.NewStore(db)),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	tm.Send(tea.KeyMsg{
		Type: tea.KeyUp,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	// then
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("ERROR: username and password cannot be empty"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*1),
	)
}

func TestLoginShouldEnterCreateUserAction(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(database.NewStore(db)),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	tm.Send(tea.KeyMsg{
		Type: tea.KeyTab,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(LoginModel)
	assert.True(t, ok)
	assert.True(t, m.CreateUserActionPicked())

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestLoginShouldSucceedAfterEnteringCreateUserFlowAndBack(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	tm := teatest.NewTestModel(
		t,
		NewLoginModel(database.NewStore(db)),
		teatest.WithInitialTermSize(300, 100),
	)

	// expected
	existingUsername := test.RandomString()
	existingPassword := test.RandomString()
	existingSalt, err := crypto.NewSalt()
	if err != nil {
		t.Fatalf("could not create salt")
	}

	existingUser := model.User{
		Uuid:     uuid.New().String(),
		Username: existingUsername,
		Password: crypto.HashPasswordWithSalt(existingPassword, existingSalt),
		Salt:     existingSalt,
	}

	// and user exists in database
	test.InsertIntoUsers(t, db, existingUser)

	// when
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(existingUsername),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyTab,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyTab,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyDown,
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(existingPassword),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyEnter,
	})

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(LoginModel)
	assert.True(t, ok)
	assert.True(t, m.LoggedInSuccessfully())

	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	assert.True(t, bytes.Contains(out, []byte("Logged in")))

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
