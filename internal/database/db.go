package database

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"

	// import for migration driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	// import for migration driver
	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = ".local/share/yubigo-pass/stores/root/yubigo-pass.db"

var db *sqlx.DB

// CreateDB Creates DB instance
func CreateDB() *sqlx.DB {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory:", err)
	}

	dbFilePath := filepath.Join(homeDir, dbFileName)

	err = os.MkdirAll(filepath.Dir(dbFilePath), 0750)
	if err != nil {
		log.Fatal("Error creating directory path:", err)
	}

	db, err := sqlx.Connect("sqlite3", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

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
	return db
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Error("Error closing database connection:", err)
		}
	}
}
