package cli

import (
	"fmt"
	"io"
	"strings"
	"time"
	"yubigo-pass/internal/app/common"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const viewListHeight = 18
const statusMessageTimeout = time.Second * 2

type passwordListItem PasswordListItem

// FilterValue allows list filtering on Title and Username.
func (i passwordListItem) FilterValue() string {
	return i.Title + " " + i.Username
}

// passwordDelegate handles rendering items in the view passwords list.
type passwordDelegate struct{}

func (d passwordDelegate) Height() int                             { return 1 }
func (d passwordDelegate) Spacing() int                            { return 0 }
func (d passwordDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d passwordDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(passwordListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s (%s)", i.Title, i.Username)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// Define a message type for clearing the status message after a timeout
type clearStatusMsg struct{}

// clearStatusCmd returns a command that sends a clearStatusMsg after a delay.
func clearStatusCmd() tea.Cmd {
	return tea.Tick(statusMessageTimeout, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

// ViewPasswordsModel is the Bubble Tea model for displaying and interacting with the password list.
type ViewPasswordsModel struct {
	list         list.Model
	statusMsg    string
	statusStyle  lipgloss.Style
	clearCmd     tea.Cmd
	showStatus   bool
	initialItems []PasswordListItem
}

// NewViewPasswordsModel creates a new instance of the ViewPasswordsModel.
// It expects a slice of items containing only identifiers (ID, Title, Username).
func NewViewPasswordsModel(items []PasswordListItem) ViewPasswordsModel {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = passwordListItem(item)
	}

	const defaultWidth = 50

	l := list.New(listItems, passwordDelegate{}, defaultWidth, viewListHeight)
	l.Title = "VIEW PASSWORDS"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)
	l.Styles.Title = titleStyle.Copy().MarginBottom(1)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.Styles.FilterPrompt = focusedStyle
	l.Styles.FilterCursor = cursorStyle

	statusOKStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateOk))

	return ViewPasswordsModel{
		list:         l,
		statusStyle:  statusOKStyle,
		initialItems: items,
	}
}

// Init initializes the ViewPasswordsModel.
func (m ViewPasswordsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and input for the view passwords screen.
func (m ViewPasswordsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := listDocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-2)
		return m, nil

	case clearStatusMsg:
		m.statusMsg = ""
		m.showStatus = false
		m.clearCmd = nil
		return m, nil

	case common.StateMsg:
		if msg.State == common.StatePasswordCopied {
			m.statusMsg = "Password copied to clipboard!"
			m.statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateOk))
			m.showStatus = true
			m.clearCmd = clearStatusCmd()
			return m, m.clearCmd
		}

	case common.ErrorMsg:
		m.statusMsg = fmt.Sprintf("Error: %v", msg.Err)
		m.statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateErr))
		m.showStatus = true
		m.clearCmd = clearStatusCmd()
		return m, m.clearCmd

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch keypress := msg.String(); keypress {
		case "b", "q", "esc":
			return m, common.ChangeStateCmd(common.StateGoBack)

		case "enter":
			selectedItem, ok := m.list.SelectedItem().(passwordListItem)
			if !ok {
				m.statusMsg = "Error: Could not get selected item."
				m.statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateErr))
				m.showStatus = true
				m.clearCmd = clearStatusCmd()
				return m, m.clearCmd
			}
			return m, common.RequestDecryptAndCopyPasswordCmd(selectedItem.ID)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the view passwords UI.
func (m ViewPasswordsModel) View() string {
	var builder strings.Builder

	listRender := listDocStyle.Render(m.list.View())
	builder.WriteString(listRender)

	if m.showStatus && m.statusMsg != "" {
		if !strings.HasSuffix(listRender, "\n") {
			builder.WriteString("\n")
		}
		builder.WriteString("\n" + m.statusStyle.Render(m.statusMsg))
	}

	return builder.String()
}

// Define a docStyle specific for this view's layout, assuming styles are defined elsewhere
var listDocStyle = lipgloss.NewStyle().Margin(1, 2)
