package cli

import (
	"fmt"
	// Removed io import
	"strings"
	"time"
	"yubigo-pass/internal/app/common"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const viewTableHeight = 10
const revealTimeout = time.Second * 5
const statusTimeout = time.Second * 2

type clearStatusMsg struct{}

type hideRevealedPasswordMsg struct{}

func clearStatusCmd() tea.Cmd {
	return tea.Tick(statusTimeout, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func hideRevealedPasswordCmd() tea.Cmd {
	return tea.Tick(revealTimeout, func(t time.Time) tea.Msg {
		return hideRevealedPasswordMsg{}
	})
}

// ViewPasswordsModel is a model for displaying and managing passwords.
type ViewPasswordsModel struct {
	table              table.Model
	statusMsg          string
	statusStyle        lipgloss.Style
	showStatus         bool
	items              []PasswordListItem
	revealedPasswordID string
	revealedTimeout    tea.Cmd
}

// NewViewPasswordsModel creates a new instance using bubbles/table.
func NewViewPasswordsModel(items []PasswordListItem) ViewPasswordsModel {
	columns := []table.Column{
		{Title: "Title", Width: 20},
		{Title: "Username", Width: 20},
		{Title: "Password", Width: 30},
		{Title: "URL", Width: 30},
	}

	rows := make([]table.Row, len(items))
	for i, item := range items {
		rows[i] = table.Row{
			item.Title,
			item.Username,
			"********",
			item.Url,
		}
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(viewTableHeight),
		table.WithStyles(s),
	)

	statusOKStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateOk))

	return ViewPasswordsModel{
		table:              t,
		items:              items,
		revealedPasswordID: "",
		statusStyle:        statusOKStyle,
	}
}

// Init initializes the ViewPasswordsModel.
func (m ViewPasswordsModel) Init() tea.Cmd {
	m.statusMsg = ""
	m.showStatus = false
	m.revealedPasswordID = ""
	rows := m.table.Rows()
	for i := range rows {
		rows[i][2] = "********"
	}
	m.table.SetRows(rows)
	return nil
}

// Update handles messages and input for the view passwords table.
func (m ViewPasswordsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width - 4)
		m.table.SetHeight(msg.Height - 10)
		return m, nil

	case clearStatusMsg:
		m.statusMsg = ""
		m.showStatus = false
		return m, nil

	case hideRevealedPasswordMsg:
		if m.revealedPasswordID != "" {
			rowIndex := m.findRowIndexByID(m.revealedPasswordID)
			if rowIndex != -1 {
				rows := m.table.Rows()
				rows[rowIndex][2] = "********"
				m.table.SetRows(rows)
			}
			m.revealedPasswordID = ""
			m.revealedTimeout = nil
		}
		return m, nil

	case common.StateMsg:
		if msg.State == common.StatePasswordCopied {
			m.statusMsg = "Password copied to clipboard!"
			m.statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateOk))
			m.showStatus = true
			cmds = append(cmds, clearStatusCmd())
		}

	case common.ErrorMsg:
		m.statusMsg = fmt.Sprintf("Error: %v", msg.Err)
		m.statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorValidateErr))
		m.showStatus = true
		cmds = append(cmds, clearStatusCmd())

	case common.PasswordDecryptedMsg:
		rowIndex := m.findRowIndexByID(msg.PasswordID)
		if rowIndex != -1 {
			rows := m.table.Rows()
			if m.revealedPasswordID != "" && m.revealedPasswordID != msg.PasswordID {
				oldRowIndex := m.findRowIndexByID(m.revealedPasswordID)
				if oldRowIndex != -1 {
					rows[oldRowIndex][2] = "********"
				}
			}
			rows[rowIndex][2] = msg.Plaintext
			m.table.SetRows(rows)
			m.revealedPasswordID = msg.PasswordID
			m.revealedTimeout = hideRevealedPasswordCmd()
			cmds = append(cmds, m.revealedTimeout)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "b", "q", "esc":
			if m.revealedPasswordID != "" {
				rowIndex := m.findRowIndexByID(m.revealedPasswordID)
				if rowIndex != -1 {
					rows := m.table.Rows()
					rows[rowIndex][2] = "********"
					m.table.SetRows(rows)
				}
				m.revealedPasswordID = ""
			}
			return m, common.ChangeStateCmd(common.StateGoBack)

		case "enter":
			if len(m.items) > 0 && m.table.Cursor() < len(m.items) {
				selectedItem := m.items[m.table.Cursor()]
				return m, common.RequestDecryptAndCopyPasswordCmd(selectedItem.ID)
			}

		case "ctrl+s":
			if len(m.items) > 0 && m.table.Cursor() < len(m.items) {
				cursor := m.table.Cursor()
				selectedItem := m.items[cursor]

				if m.revealedPasswordID == selectedItem.ID {
					rows := m.table.Rows()
					rows[cursor][2] = "********"
					m.table.SetRows(rows)
					m.revealedPasswordID = ""
					m.revealedTimeout = nil
				} else {
					return m, common.RequestDecryptPasswordCmd(selectedItem.ID)
				}
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the view passwords table UI.
func (m ViewPasswordsModel) View() string {
	var builder strings.Builder

	builder.WriteString(tableBaseStyle.Render(m.table.View()))

	if m.showStatus && m.statusMsg != "" {
		builder.WriteString("\n\n" + m.statusStyle.Render(m.statusMsg))
	} else {
		builder.WriteString("\n\n ")
	}

	help := fmt.Sprintf("\n %s | %s | %s | %s",
		blurredStyle.Render("↑/↓: Navigate"),
		blurredStyle.Render("Enter: Copy Password"),
		blurredStyle.Render("Ctrl+S: Show/Hide Password"),
		blurredStyle.Render("Esc/q/b: Back"),
	)
	builder.WriteString(help)

	return builder.String()
}

// findRowIndexByID finds the table row index corresponding to a password ID.
func (m ViewPasswordsModel) findRowIndexByID(id string) int {
	for i, item := range m.items {
		if item.ID == id {
			return i
		}
	}
	return -1
}

var tableBaseStyle = lipgloss.NewStyle().Margin(1, 2)
