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
)

func TestShouldQuitAddPasswordAction(t *testing.T) {
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
				tm := teatest.NewTestModel(
					t,
					NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession()),
					teatest.WithInitialTermSize(300, 100),
				)
				test.PressKey(tm, testCase.key)
				fm := tm.FinalModel(t)
				m, ok := fm.(AddPasswordModel)
				assert.Truef(t, ok, "final model has wrong type: %T", fm)
				assert.Truef(t, m.Cancelled, "final model is not Cancelled")
				tm.WaitFinished(t, teatest.WithFinalTimeout(time.Millisecond*100))
			},
		)
	}
}

func TestShouldAddPassword(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(database.NewStore(db), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)

	// expected
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
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	assert.True(t, ok)
	assert.True(t, m.PasswordAdded)
	result := ExtractPasswordDataFromModel(m)
	assert.Equal(t, exampleUsername, result.Username)
	assert.Equal(t, exampleTitle, result.Title)
	assert.Equal(t, examplePassword, result.Password)
	assert.Equal(t, exampleUrl, result.Url)

	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	assert.True(t, bytes.Contains(out, []byte("Password added")))

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
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
	exampleUrl := test.RandomString()

	// when
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUrl)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("ERROR: only optional fields can be empty"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*1),
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
	exampleUrl := test.RandomString()

	// when
	test.TypeString(tm, exampleTitle)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUrl)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("ERROR: only optional fields can be empty"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*1),
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
	exampleUrl := test.RandomString()

	// when
	test.TypeString(tm, exampleTitle)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUrl)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("ERROR: only optional fields can be empty"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*1),
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
	test.PressKey(tm, tea.KeyUp)
	test.PressKey(tm, tea.KeyEnter)

	// then
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("ERROR: only optional fields can be empty"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*1),
	)
}

func TestShouldAddPasswordWithEmptyUrl(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(database.NewStore(db), utils.NewEmptySession()),
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
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	assert.True(t, ok)
	assert.True(t, m.PasswordAdded)
	result := ExtractPasswordDataFromModel(m)
	assert.Equal(t, exampleUsername, result.Username)
	assert.Equal(t, exampleTitle, result.Title)
	assert.Equal(t, examplePassword, result.Password)
	assert.Equal(t, "", result.Url)

	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	assert.True(t, bytes.Contains(out, []byte("Password added")))

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestShouldNotAddPasswordIfItAlreadyExists(t *testing.T) {
	// given
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)

	userID := uuid.New().String()

	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(database.NewStore(db), utils.NewSession(userID, test.RandomString(), test.RandomString())),
		teatest.WithInitialTermSize(300, 100),
	)

	exampleTitle := test.RandomString()
	exampleUsername := test.RandomString()
	examplePassword := test.RandomString()

	// and password already exists in database
	test.InsertIntoPasswords(t, db, model.NewPassword(
		userID,
		exampleTitle,
		exampleUsername,
		examplePassword,
		test.RandomString(),
		[]byte(test.RandomString())),
	)

	// when
	test.TypeString(tm, exampleTitle)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, exampleUsername)
	test.PressKey(tm, tea.KeyDown)
	test.TypeString(tm, examplePassword)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	teatest.WaitFor(
		t,
		tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("ERROR: this password already exists, change inputs or update existing password"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*1),
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
	test.PressKey(tm, tea.KeyTab)
	test.PressKey(tm, tea.KeyTab)
	test.PressKey(tm, tea.KeyEsc)

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	assert.True(t, ok)
	assert.True(t, m.Cancelled)

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestShouldAbortAddPasswordActionAndGoBack(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewAddPasswordModel(test.NewStoreExecutorMock(), utils.NewEmptySession()),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyTab)
	test.PressKey(tm, tea.KeyEnter)

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(AddPasswordModel)
	assert.True(t, ok)
	assert.True(t, m.Back)

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
