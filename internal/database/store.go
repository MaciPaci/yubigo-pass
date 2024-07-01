package database

import (
	"database/sql"
	"errors"
	"fmt"
	"yubigo-pass/internal/app/model"

	"github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
)

// Store is a DB gateway
type Store struct {
	db *sqlx.DB
}

// NewStore returns new Store instance
func NewStore(db *sqlx.DB) Store {
	return Store{
		db: db,
	}
}

// CreateUser adds new user in DB
func (s Store) CreateUser(input model.User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `INSERT INTO users (id, username, password, salt) VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(query, input.UserID, input.Username, input.Password, input.Salt)
	if err != nil {
		_ = tx.Rollback()
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return model.NewUserAlreadyExistsError(input.Username)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	_ = tx.Commit()
	return nil
}

// GetUser fetches a user by username from DB
func (s Store) GetUser(username string) (model.User, error) {
	query := `SELECT * FROM users where username = $1`

	var user model.User
	err := s.db.QueryRowx(query, username).StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.NewUserNotFoundError(username)
		}
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s Store) AddPassword(input model.Password) error {
	tx, err := s.db.Begin()
	query := `INSERT INTO passwords (user_id, title, username, password, url, nonce) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.Exec(query, input.UserID, input.Title, input.Username, input.Password, input.Url, input.Nonce)
	if err != nil {
		_ = tx.Rollback()
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
			return model.NewPasswordAlreadyExistsError(input.UserID, input.Title, input.Username)
		}
		return fmt.Errorf("failed to create password: %w", err)
	}

	_ = tx.Commit()
	return nil
}

func (s Store) GetPassword(userID, title, username string) (model.Password, error) {
	query := `SELECT * FROM passwords WHERE user_id = $1 AND title = $2 AND username = $3`

	var password model.Password
	err := s.db.QueryRowx(query, userID, title, username).StructScan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Password{}, model.NewPasswordNotFoundError(userID, title, username)
		}
		return model.Password{}, fmt.Errorf("failed to get password: %w", err)
	}
	return password, nil
}

func (s Store) GetAllUserPasswords(username string) ([]model.Password, error) {
	query := `SELECT * FROM passwords WHERE user_id = $1`

	var passwords []model.Password
	err := s.db.Select(&passwords, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get passwords: %w", err)
	}
	return passwords, nil
}
