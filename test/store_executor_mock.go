package test

import "yubigo-pass/internal/app/model"

const (
	// ExistingUsername is username stub for StoreExecutorMock
	ExistingUsername = "existingUsername"
)

// StoreExecutorMock is a mock of StoreExecutor for testing purposes
type StoreExecutorMock struct {
}

// NewStoreExecutorMock returns StoreExecutorMock
func NewStoreExecutorMock() StoreExecutorMock {
	return StoreExecutorMock{}
}

// CreateUser mocks StoreExecutor CreateUser method
func (s StoreExecutorMock) CreateUser(input model.User) error {
	return nil
}

// GetUser mocks StoreExecutor GetUser method
func (s StoreExecutorMock) GetUser(username string) (model.User, error) {
	if username == ExistingUsername {
		return model.User{
			Username: username,
		}, nil
	}
	return model.User{}, nil
}
