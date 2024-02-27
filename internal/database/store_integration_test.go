//go:build integration

package database

import (
	"fmt"
	"testing"
	"yubigo-pass/internal/app/model"
	"yubigo-pass/test"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateUserInDB(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)
	store := NewStore(db)

	// given
	input := model.User{
		Uuid:     test.RandomString(),
		Username: test.RandomString(),
		Password: test.RandomString(),
		Salt:     test.RandomString(),
	}

	// when
	err = store.CreateUser(input)

	// then
	assert.Nil(t, err)
	user, err := store.GetUser(input.Username)
	assert.Equal(t, input, user)
}

func TestShouldNotCreateUserIfOneWithTheSameUsernameIsAlreadyInDB(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)
	store := NewStore(db)

	// given
	input := model.User{
		Uuid:     test.RandomString(),
		Username: test.RandomString(),
		Password: test.RandomString(),
		Salt:     test.RandomString(),
	}
	test.InsertIntoUsers(t, db, input)

	// expected
	expectedError := model.NewUserAlreadyExistsError(input.Username)

	// when
	err = store.CreateUser(input)

	// then
	assert.EqualError(t, err, expectedError.Error())
}

func TestShouldGetUserFromDB(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)
	store := NewStore(db)

	// given
	input := model.User{
		Uuid:     test.RandomString(),
		Username: test.RandomString(),
		Password: test.RandomString(),
		Salt:     test.RandomString(),
	}
	test.InsertIntoUsers(t, db, input)

	// when
	user, err := store.GetUser(input.Username)

	// then
	assert.Nil(t, err)
	assert.Equal(t, input, user)
}

func TestShouldNotGetUserFromDBIfNoMatchingUsernameFound(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)
	store := NewStore(db)

	// given
	username := test.RandomString()

	// expected
	expectedError := fmt.Errorf("user not found for username %s", username)

	// when
	user, err := store.GetUser(username)

	// then
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, user)
}
