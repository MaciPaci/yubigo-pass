package cli

import (
	"fmt"
	"io"
	"strings"
	"yubigo-pass/internal/app/common"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

// Constants defining main menu item labels.
const (
	ViewPasswordItem = "View your passwords"
	AddPasswordItem  = "Add a new password" // #nosec G101
	LogoutItem       = "Logout"
	QuitItem         = "Quit"
)

// item represents a selectable item in the main menu list.
type item string

// FilterValue implements list.Item interface.
func (i item) FilterValue() string { return "" }

// itemDelegate handles rendering list items.
type itemDelegate struct{}

// Height implements list.ItemDelegate interface.
func (d itemDelegate) Height() int { return 1 }

// Spacing implements list.ItemDelegate interface.
func (d itemDelegate) Spacing() int { return 0 }

// Update implements list.ItemDelegate interface.
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render implements list.ItemDelegate interface.
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

// MainMenuModel is a Bubble Tea model for the main application menu.
// It displays options and sends state change messages based on user selection.
type MainMenuModel struct {
	list     list.Model
	quitting bool
	// lastStateSent common.MsgState
}

// NewMainMenuModel creates a new instance of the MainMenuModel.
func NewMainMenuModel() MainMenuModel {
	items := []list.Item{
		item(ViewPasswordItem),
		item(AddPasswordItem),
		item(LogoutItem),
		item(QuitItem),
	}

	const defaultWidth = 40

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "MAIN MENU"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)
	l.Styles.Title = titleStyle.Copy().MarginBottom(1)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return MainMenuModel{
		list: l,
	}
}

// Init initializes the MainMenuModel
func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and user input for the main menu.
// It sends state change messages to the main application model.
func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, common.ChangeStateCmd(common.StateQuit)

		case "enter":
			selectedItem, ok := m.list.SelectedItem().(item)
			if !ok {
				return m, nil
			}

			var nextState common.MsgState
			choice := string(selectedItem)
			switch choice {
			case ViewPasswordItem:
				nextState = common.StateGoToViewPasswords
			case AddPasswordItem:
				nextState = common.StateGoToAddPassword
			case LogoutItem:
				nextState = common.StateLogout
			case QuitItem:
				m.quitting = true
				nextState = common.StateQuit
			default:
				return m, nil
			}
			return m, common.ChangeStateCmd(nextState)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the main menu UI.
func (m MainMenuModel) View() string {
	if m.quitting {
		return quitTextStyle.Render("Quitting.")
	}
	return docStyle.Render(m.list.View())
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)
