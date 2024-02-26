//go:build unit

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"yubigo-pass/internal/database"

	"github.com/stretchr/testify/assert"
)

func TestCreatePathForDB(t *testing.T) {
	// given
	oldHomeDir := os.Getenv("HOME")
	expectedPath := filepath.Join(oldHomeDir, database.DbFileName)

	// when
	actualPath := CreatePathForDB()

	// then
	assert.Equal(t, expectedPath, actualPath, "Database file path does not match expected")
}
