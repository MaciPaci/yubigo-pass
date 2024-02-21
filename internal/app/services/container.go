package services

import (
	"yubigo-pass/internal/database"
)

// Container is a struct holding all app services
type Container struct {
	Store    database.StoreExecutor
	Programs Programs
}
