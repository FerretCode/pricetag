package models

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alexedwards/argon2id"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mattn/go-sqlite3"
)

type UserModel struct {
	db *sql.DB
}

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.PasswordHash, validation.Required))
}

func (u *User) SetPasswordHash(password string) error {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	u.PasswordHash = hash

	return nil
}

func (m *UserModel) New(username, password string) (*User, error) {
	user := &User{Username: username}

	err := user.SetPasswordHash(password)
	if err != nil {
		return nil, err
	}

	err = m.Insert(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *UserModel) Insert(user *User) error {
	err := user.Validate()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (username, password)
		VALUES(?, ?)
		RETURNING id;`

	args := []any{user.Username, user.PasswordHash}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	err = m.db.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		switch sqliteErrCode(err) {
		case sqlite3.ErrConstraintUnique:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) Usernames() ([]string, error) {
	query := "SELECT username FROM users;"

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return usernames, nil
}

func (m *UserModel) IsFirstUser() (bool, error) {
	query := "SELECT COUNT(*) FROM users;"

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	var count int
	err := m.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (m *UserModel) GetWithID(id int) (*User, error) {
	query := `
		SELECT id, username, password
		FROM users WHERE id = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	var u User
	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Username,
		&u.PasswordHash,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrInvalidCredentials
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (m *UserModel) GetWithUsername(username string) (*User, error) {
	query := `
		SELECT id, username, password
		FROM users WHERE username = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	var u User
	err := m.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID,
		&u.Username,
		&u.PasswordHash,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrInvalidCredentials
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (m *UserModel) GetForCredentials(username, password string) (*User, error) {
	u, err := m.GetWithUsername(username)
	if err != nil {
		return nil, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, string(u.PasswordHash))
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE id = ?
		);`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	var exists bool
	err := m.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m *UserModel) ExistsWithUsername(username string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE username = ?
		);`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	var exists bool
	err := m.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m UserModel) Update(user *User) error {
	err := user.Validate()
	if err != nil {
		return err
	}

	query := `
		UPDATE users 
        SET username = ?, password = ?
        WHERE id = ?;`

	args := []any{
		user.Username,
		user.PasswordHash,
		user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	_, err = m.db.ExecContext(ctx, query, args...)
	if err != nil {
		switch sqliteErrCode(err) {
		case sqlite3.ErrConstraintUnique:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}
