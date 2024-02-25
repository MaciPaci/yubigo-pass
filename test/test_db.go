package test

import (
	"errors"
	"path/filepath"
	"runtime"
	"testing"
	"yubigo-pass/internal/app/model"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// SetupTestDB sets up in-memory database for testing purposes
func SetupTestDB() (*sqlx.DB, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	migrationsDir := filepath.Join(basePath, "..", "assets", "migrations")

	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		db.Close()
		return nil, err
	}

	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		db.Close()
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "sqlite3://:memory:", driver)
	if err != nil {
		db.Close()
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		db.Close()
		return nil, err
	}

	return db, nil
}

// TeardownTestDB closes in-memory test database
func TeardownTestDB(db *sqlx.DB) {
	db.Close()
}

// InsertIntoUsers inserts record into users database for tests only
func InsertIntoUsers(t *testing.T, db *sqlx.DB, input model.User) {
	query := `INSERT INTO users (id, username, password, salt) VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(query, input.Uuid, input.Username, input.Password, input.Salt)
	if err != nil {
		t.Fatalf("failed to create user: %s", err)
	}
}
