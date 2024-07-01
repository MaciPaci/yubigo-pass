package cli

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle       = focusedStyle.Copy()
	noStyle           = lipgloss.NewStyle()
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("205"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	titleStyle        = lipgloss.NewStyle().Margin(1, 0, 0, 4).Foreground(lipgloss.Color("205")).Background(lipgloss.Color("235")).Bold(true)
)

var (
	focusedSubmitButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredSubmitButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	focusedBackButton   = focusedStyle.Copy().Render("[ Back ]")
	blurredBackButton   = fmt.Sprintf("[ %s ]", blurredStyle.Render("Back"))
	focusedAddButton    = focusedStyle.Copy().Render("[ Add ]")
	blurredAddButton    = fmt.Sprintf("[ %s ]", blurredStyle.Render("Add"))
)

const (
	validateOkPrefix  = "✔"
	validateErrPrefix = "✘"
	colorValidateOk   = "2"
	colorValidateErr  = "1"
)
