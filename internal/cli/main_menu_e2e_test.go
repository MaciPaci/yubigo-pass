//go:build e2e

package cli

import (
	"bytes"
	"io"
	"testing"
	"time"
	"yubigo-pass/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"
)

func TestShouldQuitMainMenuAction(t *testing.T) {

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
					NewMainMenuModel(),
					teatest.WithInitialTermSize(300, 100),
				)
				test.PressKey(tm, testCase.key)
				fm := tm.FinalModel(t)
				m, ok := fm.(MainMenuModel)
				assert.Truef(t, ok, "final model has wrong type: %T", fm)
				assert.Truef(t, m.quitting, "final model is not quitting")
				out, err := io.ReadAll(tm.FinalOutput(t))
				if err != nil {
					t.Error(err)
				}
				assert.True(t, bytes.Contains(out, []byte("Quitting.")))
				tm.WaitFinished(t, teatest.WithFinalTimeout(time.Millisecond*100))
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
	test.PressKey(tm, tea.KeyEnter)

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	assert.True(t, ok)
	assert.Equal(t, m.Choice, GetPasswordItem)

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestMainMenuShouldChooseViewPasswords(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewMainMenuModel(),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	assert.True(t, ok)
	assert.Equal(t, m.Choice, ViewPasswordItem)

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestMainMenuShouldChooseAddPasswords(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewMainMenuModel(),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	assert.True(t, ok)
	assert.Equal(t, m.Choice, AddPasswordItem)

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestMainMenuShouldChooseLogout(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewMainMenuModel(),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	assert.True(t, ok)
	assert.Equal(t, m.Choice, LogoutItem)

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestMainMenuShouldChooseQuit(t *testing.T) {
	// given
	tm := teatest.NewTestModel(
		t,
		NewMainMenuModel(),
		teatest.WithInitialTermSize(300, 100),
	)

	// when
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyDown)
	test.PressKey(tm, tea.KeyEnter)

	// then
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	assert.True(t, ok)
	assert.Equal(t, m.Choice, QuitItem)

	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	assert.True(t, bytes.Contains(out, []byte("Quitting.")))

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
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
	fm := tm.FinalModel(t)
	m, ok := fm.(MainMenuModel)
	assert.Truef(t, ok, "final model has wrong type: %T", fm)
	assert.Truef(t, m.quitting, "final model is not quitting")

	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	assert.True(t, bytes.Contains(out, []byte("Quitting.")))

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Millisecond*100))
}
