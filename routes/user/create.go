package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"

	"github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RenderCreateUserPage(w http.ResponseWriter, r *http.Request, templates *template.Template) error {
	err := templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		return err
	}

	return nil
}

func Create(w http.ResponseWriter, r *http.Request, db *sql.DB) (status int, err error) {
	currentUsersQuery := squirrel.Select("*").From("users")

	sql, args, err := currentUsersQuery.ToSql()
	if err != nil {
		return 500, err
	}

	res, err := db.Query(sql, args)
	if err != nil {
		return 500, err
	}

	admin := false

	if !res.Next() {
		admin = true
	}

	createUserRequest := createUserRequest{}

	err = processBody(r, &createUserRequest)
	if err != nil {
		return 500, err
	}

	if len(createUserRequest.Username) < 3 {
		return 400, errors.New("your username must be at least 3 characters")
	}

	if len(createUserRequest.Username) > 16 {
		return 400, errors.New("your username cannot be over 16 charactres")
	}

	if len(createUserRequest.Password) > 74 {
		return 400, errors.New("your password cannot be over 74 characters")
	}

	if len(createUserRequest.Password) < 12 {
		return 400, errors.New("your password must be at least 12 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), 14)
	if err != nil {
		return 500, err
	}

	return 200, nil
}

func processBody(r *http.Request, to interface{}) error {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, to); err != nil {
		return err
	}

	return nil
}
