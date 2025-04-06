package model

import "fmt"

// UserAlreadyExistsError is an error if user already exists in db
type UserAlreadyExistsError struct {
	Username string
}

func (e UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("failed to create user: user already exists: %s", e.Username)
}

// NewUserAlreadyExistsError returns new UserAlreadyExistsError instance
func NewUserAlreadyExistsError(username string) UserAlreadyExistsError {
	return UserAlreadyExistsError{Username: username}
}

// UserNotFoundError is an error if user is not found in db
type UserNotFoundError struct {
	Username string
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user not found for username %s", e.Username)
}

// NewUserNotFoundError returns new UserNotFoundError instance
func NewUserNotFoundError(username string) UserNotFoundError {
	return UserNotFoundError{Username: username}
}

// PasswordNotFoundError is an error if password is not found in db
type PasswordNotFoundError struct {
	UserID   string
	Title    string
	Username string
}

// NewPasswordNotFoundError returns new PasswordNotFoundError instance
func NewPasswordNotFoundError(userID, title, username string) PasswordNotFoundError {
	return PasswordNotFoundError{
		UserID:   userID,
		Title:    title,
		Username: username,
	}
}

func (e PasswordNotFoundError) Error() string {
	return fmt.Sprintf("password not found for user %s, title %s, username %s", e.UserID, e.Title, e.Username)
}

// PasswordAlreadyExistsError is an error if password already exists in db
type PasswordAlreadyExistsError struct {
	UserID   string
	Title    string
	Username string
}

// NewPasswordAlreadyExistsError returns new PasswordAlreadyExistsError instance
func NewPasswordAlreadyExistsError(userID, title, username string) PasswordAlreadyExistsError {
	return PasswordAlreadyExistsError{
		UserID:   userID,
		Title:    title,
		Username: username,
	}
}

func (e PasswordAlreadyExistsError) Error() string {
	return fmt.Sprintf("password already exists for user %s, title %s, username %s", e.UserID, e.Title, e.Username)
}
