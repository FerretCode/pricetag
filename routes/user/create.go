package user

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/charmbracelet/log"
	"github.com/ferretcode/pricetag/session"
	"github.com/ferretcode/pricetag/types"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	Username string
	Password string
}

func RenderCreateUserPage(w http.ResponseWriter, r *http.Request, templates *template.Template) error {
	err := templates.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		return err
	}

	return nil
}

func Create(w http.ResponseWriter, r *http.Request, db *sqlx.DB, session *session.SessionManager) (status int, err error) {
	err = deleteExistingSession(w, r, session)
	if err != nil {
		return 500, err
	}

	err = r.ParseForm()
	if err != nil {
		return 500, err
	}

	currentUsersQuery := squirrel.Select("*").From("User")

	sql, _, err := currentUsersQuery.ToSql()
	if err != nil {
		return 500, err
	}

	res, err := db.Query(sql)
	if err != nil {
		return 500, err
	}

	admin := false

	if !res.Next() {
		admin = true
	}

	createUserRequest := createUserRequest{
		Username: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	}

	status, err = validateUsernameAndPassword(createUserRequest.Username, createUserRequest.Password)
	if err != nil {
		return status, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), 14)
	if err != nil {
		return 500, err
	}

	userID, err := createUserDBRecord(createUserRequest, string(hash), admin, db)
	if err != nil {
		return 500, err
	}

	sessionID := session.CreateSession(userID)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Domain:  os.Getenv("COOKIE_DOMAIN"),
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	})

	log.Info("successfully created user", "id", userID, "username", createUserRequest.Username)

	http.Redirect(w, r, "/dashboard/home", http.StatusFound)

	return 200, nil
}

func deleteExistingSession(w http.ResponseWriter, r *http.Request, session *session.SessionManager) error {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil
		}
		return err
	}

	_, err = session.GetSession(cookie.Value)
	if err == nil {
		session.DeleteSession(cookie.Value)
	}

	deleteCookie := &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Domain:  os.Getenv("COOKIE_DOMAIN"),
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, deleteCookie)

	return nil
}

func createUserDBRecord(cur createUserRequest, hash string, admin bool, db *sqlx.DB) (userID int, err error) {
	createNewUserQuery := squirrel.
		Insert("User").
		Columns("Username", "Password").
		Values(cur.Username, hash).
		Suffix("RETURNING ID")

	sql, args, err := createNewUserQuery.ToSql()
	if err != nil {
		return 0, err
	}

	user := types.User{}

	tx, err := db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	rows, err := tx.Queryx(sql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.New("error getting returned user")
	}

	createNewPermissionQuery := squirrel.
		Insert("Permission").
		Columns("UserID", "Admin", "ManageServices", "ManageTags", "ManageForwarding", "ViewLogs").
		Values(user.ID, admin, false, false, false, false).
		Suffix("RETURNING ID")

	sql, args, err = createNewPermissionQuery.ToSql()
	if err != nil {
		return 0, err
	}

	permission := types.Permission{}

	rows, err = tx.Queryx(sql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&permission)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.New("error getting returned permission")
	}

	updateUserPermissionQuery := squirrel.
		Update("User").
		Set("PermissionID", permission.ID).
		Where("ID", user.ID)

	sql, args, err = updateUserPermissionQuery.ToSql()
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func validateUsernameAndPassword(username string, password string) (status int, err error) {
	if len(username) < 3 {
		return 400, errors.New("your username must be at least 3 characters")
	}

	if len(username) > 16 {
		return 400, errors.New("your username cannot be over 16 charactres")
	}

	if len(password) > 74 {
		return 400, errors.New("your password cannot be over 74 characters")
	}

	if len(password) < 12 {
		return 400, errors.New("your password must be at least 12 characters")
	}

	return 200, nil
}
