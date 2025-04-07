package main

import (
	"fmt"
	"os"
	"yubigo-pass/internal/app/cli"
	"yubigo-pass/internal/app/services"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirupsen/logrus"
)

// main is the entry point of the yubigo-pass application.
// It sets up services, initializes the main Bubble Tea application model,
// and runs the TUI program loop.
func main() {
	setupLogging() // Configure logging early

	container, err := services.Build()
	if err != nil {
		logrus.Fatalf("Failed to build application services: %v", err)
		fmt.Fprintf(os.Stderr, "Error: Failed to build application services: %v\n", err)
		os.Exit(1)
	}

	logrus.Info("Application starting...")

	appModel := cli.NewAppModel(container)
	program := tea.NewProgram(appModel, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		logrus.Errorf("Application error during run: %v", err)
		fmt.Fprintf(os.Stderr, "Error: Application exited unexpectedly: %v\n", err)
		os.Exit(1)
	}
}

// setupLogging configures the logrus logger.
func setupLogging() {
	logFile, err := os.OpenFile("yubigo-pass.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // #nosec G302
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to open log file: %v. Logging to stderr.\n", err)
		logrus.SetOutput(os.Stderr)
	} else {
		logrus.SetOutput(logFile)
		logrus.SetLevel(logrus.InfoLevel)
	}
}
