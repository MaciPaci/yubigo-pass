//go:build e2e

package cli

import (
	"testing"
	"yubigo-pass/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestShouldQuitMainMenuAction(t *testing.T) {

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
					NewMainMenuModel(),
					teatest.WithInitialTermSize(300, 100),
				)
				test.PressKey(tm, testCase.key)
				// Model finishes because it sent StateQuit command
				err := tm.Quit()
				require.NoError(t, err, "Failed to quit the model")
				fm := tm.FinalModel(t)
				m, ok := fm.(MainMenuModel)
				require.Truef(t, ok, "final model has wrong type: %T", fm)
				assert.True(t, m.quitting)
			},
		)
	}
}

func TestMainMenuShouldChooseGetPassword(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewMainMenuModel(),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyEnter) // Select first item (Get Password)

	// then
	// Model finishes because it sent StateGoToGetPassword command
	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.False(t, m.quitting)
}

func TestMainMenuShouldChooseViewPasswords(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewMainMenuModel(),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyDown) // -> View Passwords
	test.PressKey(tm, tea.KeyEnter)

	// then
	// Model finishes because it sent StateGoToViewPasswords command
	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.False(t, m.quitting)
}

func TestMainMenuShouldChooseAddPasswords(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewMainMenuModel(),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyDown) // -> View Passwords
	test.PressKey(tm, tea.KeyDown) // -> Add Password
	test.PressKey(tm, tea.KeyEnter)

	// then
	// Model finishes because it sent StateGoToAddPassword command
	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.False(t, m.quitting)
}

func TestMainMenuShouldChooseLogout(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewMainMenuModel(),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyDown) // -> View Passwords
	test.PressKey(tm, tea.KeyDown) // -> Add Password
	test.PressKey(tm, tea.KeyDown) // -> Logout
	test.PressKey(tm, tea.KeyEnter)

	// then
	// Model finishes because it sent StateLogout command
	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.False(t, m.quitting)
}

func TestMainMenuShouldChooseQuit(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewMainMenuModel(),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyDown) // -> View Passwords
	test.PressKey(tm, tea.KeyDown) // -> Add Password
	test.PressKey(tm, tea.KeyDown) // -> Logout
	test.PressKey(tm, tea.KeyDown) // -> Quit
	test.PressKey(tm, tea.KeyEnter)

	// then
	// Model finishes because it sent StateQuit command
	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.True(t, m.quitting)
}

func TestMainMenuShouldQuitByPressingQ(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewMainMenuModel(),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.TypeString(tm, "q")

	// then
	// Model finishes because it sent StateQuit command
	err := tm.Quit()
	require.NoError(t, err, "Failed to quit the model")
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	require.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.True(t, m.quitting)
}
