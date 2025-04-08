package database

import (
	"context"
	"yubigo-pass/internal/app/model"
)

// StoreExecutor defines the interface for all database operations.
type StoreExecutor interface {
	// User management
	CreateUser(user *model.User) error
	GetUser(username string) (*model.User, error)

	// Password management
	AddPassword(password *model.Password) error
	GetPassword(userID, title, username string) (*model.Password, error)
	GetPasswordByID(id string) (*model.Password, error)
	GetAllUserPasswords(userID string) ([]model.Password, error)

	// Connection management
	Close() error
	PingContext(ctx context.Context) error
}
