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
		UserID:   test.RandomString(),
		Username: test.RandomString(),
		Password: test.RandomString(),
		Salt:     test.RandomString(),
	}

	// when
	err = store.CreateUser(input)

	// then
	assert.NoError(t, err)
	user := test.GetUser(t, db, input.Username)
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
		UserID:   test.RandomString(),
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
		UserID:   test.RandomString(),
		Username: test.RandomString(),
		Password: test.RandomString(),
		Salt:     test.RandomString(),
	}
	test.InsertIntoUsers(t, db, input)

	// when
	user, err := store.GetUser(input.Username)

	// then
	assert.NoError(t, err)
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

func TestShouldCreatePasswordInDB(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)
	store := NewStore(db)

	// given
	input := model.Password{
		UserID:   test.RandomString(),
		Title:    test.RandomString(),
		Username: test.RandomString(),
		Password: test.RandomString(),
		Url:      test.RandomString(),
		Nonce:    []byte(test.RandomString()),
	}

	// when
	err = store.AddPassword(input)

	// then
	assert.NoError(t, err)
	password := test.GetPassword(t, db, input.UserID, input.Title, input.Username)
	assert.Equal(t, input, password)
}

func TestShouldNotCreatePasswordIfOneWithTheSameIsAlreadyInDB(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)
	store := NewStore(db)

	// given
	input := model.Password{
		UserID:   test.RandomString(),
		Title:    test.RandomString(),
		Username: test.RandomString(),
		Password: test.RandomString(),
		Url:      test.RandomString(),
		Nonce:    []byte(test.RandomString()),
	}
	test.InsertIntoPasswords(t, db, input)

	// expected
	expectedError := model.NewPasswordAlreadyExistsError(input.UserID, input.Title, input.Username)

	// when
	err = store.AddPassword(input)

	// then
	assert.EqualError(t, err, expectedError.Error())
}

func TestShouldGetPasswordFromDB(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)
	store := NewStore(db)

	// given
	input := model.Password{
		UserID:   test.RandomString(),
		Title:    test.RandomString(),
		Username: test.RandomString(),
		Password: test.RandomString(),
		Url:      test.RandomString(),
		Nonce:    []byte(test.RandomString()),
	}
	test.InsertIntoPasswords(t, db, input)

	// when
	user, err := store.GetPassword(input.UserID, input.Title, input.Username)

	// then
	assert.NoError(t, err)
	assert.Equal(t, input, user)
}

func TestShouldNotGetPasswordFromDBIfNoMatchingPasswordFound(t *testing.T) {
	// setup
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer test.TeardownTestDB(db)
	store := NewStore(db)

	// given
	userID := test.RandomString()
	title := test.RandomString()
	username := test.RandomString()

	// expected
	expectedError := model.NewPasswordNotFoundError(userID, title, username)

	// when
	password, err := store.GetPassword(userID, title, username)

	// then
	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, password)
}
