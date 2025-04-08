package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"yubigo-pass/internal/app/model"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

// Store provides a concrete implementation of StoreExecutor using sqlx.
type Store struct {
	db *sqlx.DB
}

// Compile-time check to ensure *Store implements StoreExecutor.
var _ StoreExecutor = (*Store)(nil)

// NewStore creates and returns a new Store instance backed by the provided sqlx.DB.
func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

// Close terminates the underlying database connection.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// PingContext checks the database connection.
func (s *Store) PingContext(ctx context.Context) error {
	if s.db == nil {
		return errors.New("database connection is not initialized")
	}
	return s.db.PingContext(ctx)
}

// CreateUser inserts a new user record into the database.
func (s *Store) CreateUser(input *model.User) error {
	if s.db == nil {
		return errors.New("database is not initialized")
	}
	query := `INSERT INTO users (id, username, password, salt) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, input.UserID, input.Username, input.Password, input.Salt)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) &&
			errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) &&
			errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return model.NewUserAlreadyExistsError(input.Username)
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

// GetUser retrieves a user by their username.
// Returns model.UserNotFoundError if the user is not found.
func (s *Store) GetUser(username string) (*model.User, error) {
	if s.db == nil {
		return nil, errors.New("database is not initialized")
	}
	query := `SELECT id, username, password, salt FROM users WHERE username = ?`
	var user model.User
	err := s.db.Get(&user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.NewUserNotFoundError(username)
		}
		return nil, fmt.Errorf("failed to get user by username '%s': %w", username, err)
	}
	return &user, nil
}

// AddPassword adds a new password entry to the database.
func (s *Store) AddPassword(input *model.Password) error {
	if s.db == nil {
		return errors.New("database is not initialized")
	}
	query := `INSERT INTO passwords (id, user_id, title, username, password, url, nonce) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, input.ID, input.UserID, input.Title, input.Username, input.Password, input.Url, input.Nonce)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return model.NewPasswordAlreadyExistsError(input.UserID, input.Title, input.Username)
		}
		return fmt.Errorf("failed to insert password: %w", err)
	}
	return nil
}

// GetPassword retrieves a specific password entry by user ID, title, and username.
// Returns model.PasswordNotFoundError if not found.
func (s *Store) GetPassword(userID, title, username string) (*model.Password, error) {
	if s.db == nil {
		return nil, errors.New("database is not initialized")
	}
	query := `SELECT user_id, title, username, password, url, nonce FROM passwords WHERE user_id = ? AND title = ? AND username = ?`
	var password model.Password
	err := s.db.Get(&password, query, userID, title, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.NewPasswordNotFoundError(userID, title, username)
		}
		return nil, fmt.Errorf("failed to get password by title/username: %w", err)
	}
	return &password, nil
}

// GetPasswordByID retrieves a single password entry by its primary key (ID).
// Returns sql.ErrNoRows if not found.
func (s *Store) GetPasswordByID(id string) (*model.Password, error) {
	if s.db == nil {
		return nil, errors.New("database connection is not initialized")
	}
	query := `SELECT id, user_id, title, username, password, url, nonce FROM passwords WHERE id = ?`
	var password model.Password
	err := s.db.Get(&password, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get password by ID %s: %w", id, err)
	}
	return &password, nil
}

// GetAllUserPasswords retrieves all password entries for a given user ID.
// Returns an empty slice if no passwords are found.
func (s *Store) GetAllUserPasswords(userID string) ([]model.Password, error) {
	if s.db == nil {
		return nil, errors.New("database is not initialized")
	}
	query := `SELECT id, user_id, title, username, password, url, nonce FROM passwords WHERE user_id = ? ORDER BY title, username`
	var passwords []model.Password
	err := s.db.Select(&passwords, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all passwords for user %s: %w", userID, err)
	}
	return passwords, nil
}
