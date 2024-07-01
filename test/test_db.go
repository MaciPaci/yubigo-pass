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
		_ = db.Close()
		return nil, err
	}

	driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "sqlite3://:memory:", driver)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

// TeardownTestDB closes in-memory test database
func TeardownTestDB(db *sqlx.DB) {
	_ = db.Close()
}

// InsertIntoUsers inserts record into users database for testing purposes
func InsertIntoUsers(t *testing.T, db *sqlx.DB, input model.User) {
	query := `INSERT INTO users (id, username, password, salt) VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(query, input.UserID, input.Username, input.Password, input.Salt)
	if err != nil {
		t.Fatalf("failed to create user: %s", err)
	}
}

// GetUser fetches user by username for testing purposes
func GetUser(t *testing.T, db *sqlx.DB, username string) model.User {
	query := `SELECT * FROM users where username = $1`

	var user model.User
	err := db.QueryRowx(query, username).StructScan(&user)
	if err != nil {
		t.Fatalf("failed to get user: %s", err)
	}

	return user
}
