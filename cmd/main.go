package main

import (
	"fmt"
	"os"
	"yubigo-pass/internal/cli"
	"yubigo-pass/internal/database"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	database.CreateDB()
	_, err := tea.NewProgram(cli.NewCreateUserModel()).Run()
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
