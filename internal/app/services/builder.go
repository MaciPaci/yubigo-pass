package services

import (
	"fmt"
	"yubigo-pass/internal/app/utils"
	"yubigo-pass/internal/database"
)

// Build initializes and wires up foundational application dependencies.
// It now only focuses on services like the database store.
func Build() (Container, error) {
	db, err := database.CreateDB(utils.CreatePathForDB(), database.MigrationPath)
	if err != nil {
		return Container{}, fmt.Errorf("error initializing database: %w", err)
	}

	store := database.NewStore(db)

	return Container{
		Store: store,
	}, nil
}
