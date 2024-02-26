//go:build integration

package database

import (
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
	db := CreateDB(tempDBFilePath, migrationsDir)

	// then
	assert.NotNil(t, db)
	_, err = os.Stat(tempDBFilePath)
	assert.NoError(t, err)

	// and close DB connection
	CloseDB()
	err = db.Ping()

	// then connection is closed
	assert.EqualError(t, err, "sql: database is closed")
}
