package model

// UserAlreadyExistsError is an error if user already exists in db
type UserAlreadyExistsError struct {
	err error
}

func (e UserAlreadyExistsError) Error() string {
	return e.err.Error()
}

// NewUserAlreadyExistsError returns new UserAlreadyExistsError instance
func NewUserAlreadyExistsError(err error) UserAlreadyExistsError {
	return UserAlreadyExistsError{err: err}
}
