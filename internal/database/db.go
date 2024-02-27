package database

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

// DbFileName is a path to DB local file
const DbFileName = ".local/share/yubigo-pass/stores/root/yubigo-pass.db"

// MigrationPath is a path to migration directory
const MigrationPath = "file://assets/migrations"

var db *sqlx.DB

// CreateDB Creates DB instance
func CreateDB(dbFilePath, migrationPath string) (*sqlx.DB, error) {
	err := os.MkdirAll(filepath.Dir(dbFilePath), 0750)
	if err != nil {
		return nil, fmt.Errorf("error creating directory path: %w", err)
	}

	db, err = sqlx.Connect("sqlite3", dbFilePath)
	if err != nil {
		CloseDB()
		return nil, fmt.Errorf("error creating database instance: %w", err)
	}

	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		CloseDB()
		return nil, fmt.Errorf("error creating database driver: %w", err)
	}

	log.Info("Starting migration")

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "sqlite3", driver)
	if err != nil {
		CloseDB()
		return nil, fmt.Errorf("error creating migration instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		CloseDB()
		return nil, fmt.Errorf("error during migration: %w", err)
	}

	log.Info("Migration successful!")
	return db, nil
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Error("Error closing database connection:", err)
		}
	}
}
