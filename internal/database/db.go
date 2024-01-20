package database

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = ".local/share/yubigo-pass/stores/root/yubigo-pass.db"

// CreateDB Creates DB instance
func CreateDB() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory:", err)
	}

	dbFilePath := filepath.Join(homeDir, dbFileName)

	err = os.MkdirAll(filepath.Dir(dbFilePath), 0755)
	if err != nil {
		log.Fatal("Error creating directory path:", err)
	}

	db, err := sqlx.Connect("sqlite3", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Starting migration")

	m, err := migrate.NewWithDatabaseInstance("file://assets/migrations", "sqlite3", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	log.Info("Migration successful!")
}
