package utils

import (
	"os"
	"path/filepath"
	"yubigo-pass/internal/database"

	log "github.com/sirupsen/logrus"
)

// CreatePathForDB creates the absolute file path for the database file
func CreatePathForDB() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory:", err)
	}

	return filepath.Join(homeDir, database.DbFileName)
}
