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
