package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const listHeight = 14

const (
	GetPasswordItem  = "Get password"
	ViewPasswordItem = "View your passwords"
	AddPasswordItem  = "Add a new password" // #nosec G101
	LogoutItem       = "Logout"
	QuitItem         = "Quit"
)

type item string

// FilterValue for item
func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

// Height for itemDelegate
func (d itemDelegate) Height() int { return 1 }

// Spacing for itemDelegate
func (d itemDelegate) Spacing() int { return 0 }

// Update for itemDelegate
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render for itemDelegate
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// MainMenuModel is a main menu model
type MainMenuModel struct {
	list     list.Model
	Choice   string
	quitting bool
}

// NewMainMenuModel returns new MainMenuModel
func NewMainMenuModel() MainMenuModel {
	items := []list.Item{
		item(GetPasswordItem),
		item(ViewPasswordItem),
		item(AddPasswordItem),
		item(LogoutItem),
		item(QuitItem),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want to do?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return MainMenuModel{
		list: l,
	}
}

// Init initializes for MainMenuModel
func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

// Update updates MainMenuModel
func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.Choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders MainMenuModel
func (m MainMenuModel) View() string {
	if m.quitting || m.Choice == QuitItem {
		return quitTextStyle.Render("Quitting.")
	}
	switch m.Choice {
	case GetPasswordItem:
		return quitTextStyle.Render(GetPasswordItem)
	case ViewPasswordItem:
		return quitTextStyle.Render(ViewPasswordItem)
	case AddPasswordItem:
		return quitTextStyle.Render(AddPasswordItem)
	case LogoutItem:
		return quitTextStyle.Render(LogoutItem)
	}
	return "\n" + m.list.View()
}
