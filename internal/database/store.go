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

	_, err = tx.Exec(query, input.Uuid, input.Username, input.Password, input.Salt)
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
