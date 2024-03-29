//go:build e2e

package cli

import (
	"bytes"
	"io"
	"testing"
	"time"
	"yubigo-pass/internal/database"
	"yubigo-pass/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
)

func TestShouldQuitCreateUserProgram(t *testing.T) {
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
				tm.Send(tea.KeyMsg{
					Type: testCase.key,
				})
				fm := tm.FinalModel(t)
				m, ok := fm.(CreateUserModel)
				assert.Truef(t, ok, "final model has wrong type: %T", fm)
				assert.Truef(t, m.cancelled, "final model is not cancelled")
				tm.WaitFinished(t, teatest.WithFinalTimeout(time.Millisecond*100))
			},
		)
	}
}

func TestShouldCreateUser(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewCreateUserModel(test.NewStoreExecutorMock()),
		teatest.WithInitialTermSize(300, 100),
	)

	// expected
	exampleUsername := test.RandomString()
	examplePassword := test.RandomString()

	// when
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

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(CreateUserModel)
	assert.True(t, ok)
	assert.True(t, m.finished)
	assert.Equal(t, exampleUsername, m.inputs[0].Value())
	assert.Equal(t, examplePassword, m.inputs[1].Value())

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
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(exampleUsername),
	})

	tm.Send(tea.KeyMsg{
		Type: tea.KeyDown,
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
	tm.Send(tea.KeyMsg{
		Type: tea.KeyUp,
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

	examplePassword := test.RandomString()

	// and username already exists in database
	_, err = db.Exec("INSERT INTO users (id, username, password, salt) VALUES (?, ?, ?, ?)", "1", test.ExistingUsername, test.RandomString(), test.RandomString())
	if err != nil {
		t.Fatalf("failed to insert record to db: %v", err)
	}

	// when
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(test.ExistingUsername),
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
