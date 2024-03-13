//go:build e2e

package cli

import (
	"bytes"
	"io"
	"testing"
	"time"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/internal/database"
	"yubigo-pass/test"

	"github.com/google/uuid"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
)

func TestShouldQuitCreateUserAction(t *testing.T) {
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(test.NewStoreExecutorMock()),
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
				test.PressKey(tm, testCase.key)
				fm := tm.FinalModel(t)
				m, ok := fm.(CreateUserModel)
				assert.Truef(t, ok, "final model has wrong type: %T", fm)
				assert.Truef(t, m.Cancelled, "final model is not Cancelled")
				tm.WaitFinished(t, teatest.WithFinalTimeout(time.Millisecond*100))
			},
		)
	}
}

func TestShouldCreateUser(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(database.NewStore(db)),
		teatest.WithInitialTermSize(300, 100),
	)

	// expected
	exampleUsername := test.RandomString()
	examplePassword := test.RandomString()

	// when
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(CreateUserModel)
	assert.True(t, ok)
	assert.True(t, m.UserCreated)
	resultUsername, resultPassword := ExtractDataFromModel(m)
	assert.Equal(t, exampleUsername, resultUsername)
	assert.Equal(t, examplePassword, resultPassword)

	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	assert.True(t, bytes.Contains(out, []byte("User created successfully")))

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
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
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

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
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

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

func TestShouldNotCreateUserWithBothInputFieldsEmpty(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyUp)
	test.PressKey(tm, tea.KeyEnter)

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

func TestShouldNotCreateUserIfUsernameAlreadyExists(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)

	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(database.NewStore(db)),
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
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("ERROR: username already exists"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*1),
	)
}

func TestShouldAbortCreateUserActionAfterGoingBackAndForth(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyTab)
	test.PressKey(tm, tea.KeyTab)
	test.PressKey(tm, tea.KeyTab)
	test.PressKey(tm, tea.KeyEnter)

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(CreateUserModel)
	assert.True(t, ok)
	assert.True(t, m.UserCreationAborted)

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
