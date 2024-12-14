package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mattn/go-sqlite3"
)

const ctxTimeout = 3 * time.Second

type Models struct {
	User       *UserModel
	Permission *PermissionModel
}

func New(db *sql.DB) Models {
	return Models{
		User:       &UserModel{db},
		Permission: &PermissionModel{db},
	}
}

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateUsername  = errors.New("models: duplicate username")
	ErrEditConflict       = errors.New("models: edit conflict")
)

func sqliteErrCode(err error) sqlite3.ErrNoExtended {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.ExtendedCode
	}

	return 0
}
