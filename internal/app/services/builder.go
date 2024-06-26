package services

import (
	"yubigo-pass/internal/app/utils"
	"yubigo-pass/internal/database"

	"github.com/sirupsen/logrus"
)

// Build builds all app services
func Build() Container {
	db, err := database.CreateDB(utils.CreatePathForDB(), database.MigrationPath)
	if err != nil {
		logrus.Fatalf("error building database: %s", err)
	}
	store := database.NewStore(db)
	teaModels := InitTeaModels(store)

	return Container{
		Store:  store,
		Models: teaModels,
	}
}
