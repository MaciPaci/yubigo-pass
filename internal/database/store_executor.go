package database

import (
	"yubigo-pass/internal/app/model"
)

// StoreExecutor is an interface for DB access
type StoreExecutor interface {
	CreateUser(input model.User) error
	GetUser(username string) (model.User, error)
	AddPassword(password model.Password) error
	GetPassword(userID, title, username string) (model.Password, error)
	GetAllUserPasswords(userID string) ([]model.Password, error)
}
