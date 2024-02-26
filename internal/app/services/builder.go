package services

import (
	"yubigo-pass/internal/app/utils"
	"yubigo-pass/internal/database"
)

// Build builds all app services
func Build() Container {
	db := database.CreateDB(utils.CreatePathForDB(), database.MigrationPath)
	store := database.NewStore(db)
	programs := InitPrograms(store)

	return Container{
		Store:    store,
		Programs: programs,
	}
}
