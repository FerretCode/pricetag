package user

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/ferretcode/pricetag/session"
	"github.com/ferretcode/pricetag/types"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type loginUserRequest struct {
	Username string
	Password string
}

func RenderLoginPage(w http.ResponseWriter, r *http.Request, templates *template.Template) error {
	err := templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		return err
	}

	return nil
}

func Login(w http.ResponseWriter, r *http.Request, db *sqlx.DB, session *session.SessionManager) (status int, err error) {
	err = r.ParseForm()
	if err != nil {
		return 500, err
	}

	cookie, err := r.Cookie("session_id")
	if err != nil && err != http.ErrNoCookie {
		return 500, err
	}

	if cookie != nil {
		_, err := session.GetSession(cookie.Value)
		if err == nil {
			http.Redirect(w, r, "/dashboard/home", http.StatusFound)
			return 200, nil
		}
	}

	loginUserRequest := loginUserRequest{
		Username: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	}

	loginUsersQuery := squirrel.
		Select("*").
		From("User").
		Where(squirrel.Eq{"Username": loginUserRequest.Username})

	query, args, err := loginUsersQuery.ToSql()
	if err != nil {
		return 500, err
	}

	user := types.User{}
	err = db.Get(&user, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return 404, errors.New("user not found")
		}
		return 500, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUserRequest.Password))
	if err != nil {
		return 403, errors.New("the password is not correct")
	}

	sessionID := session.CreateSession(user.ID)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Domain:  os.Getenv("COOKIE_DOMAIN"),
		Path:    "/",
		Value:   sessionID,
		Expires: time.Now().Add(24 * time.Hour),
	})

	http.Redirect(w, r, "/dashboard/home", http.StatusFound)

	return 200, nil
}
