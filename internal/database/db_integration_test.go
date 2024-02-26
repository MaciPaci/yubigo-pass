//go:build integration

package database

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDBAndThenCloseConnection(t *testing.T) {
	// given
	tempDir := t.TempDir()
	tempDBFilePath := filepath.Join(tempDir, "test.db")
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	migrationsDir := filepath.Join("file://", cwd, "../../assets/migrations")

	// when
	db, err := CreateDB(tempDBFilePath, migrationsDir)

	// then
	assert.Nil(t, err)
	assert.NotNil(t, db)
	_, err = os.Stat(tempDBFilePath)
	assert.NoError(t, err)

	// and close DB connection
	CloseDB()
	err = db.Ping()

	// then connection is closed
	assert.EqualError(t, err, "sql: database is closed")
}

func TestCreateDBShouldFailCreatingDirectory(t *testing.T) {
	// given
	incorrectPath := t.TempDir()

	// expected
	expectedError := fmt.Errorf("error creating database instance: unable to open database file: is a directory")

	// when
	db, err := CreateDB(incorrectPath, "")

	// then
	assert.EqualError(t, err, expectedError.Error())
	assert.Nil(t, db)
}

func TestCreateDBShouldFailMigration(t *testing.T) {
	// given
	tempDir := t.TempDir()
	tempDBFilePath := filepath.Join(tempDir, "test.db")

	// expected
	expectedError := fmt.Errorf("error creating migration instance: URL cannot be empty")

	// when
	db, err := CreateDB(tempDBFilePath, "")

	// then
	assert.EqualError(t, err, expectedError.Error())
	assert.Nil(t, db)
}
