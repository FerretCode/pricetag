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

	log.Info("existing session deleted")

	err = r.ParseForm()
	if err != nil {
		return 500, err
	}

	log.Info("form parsed successfully")

	currentUsersQuery := squirrel.Select("*").From("User")

	sql, _, err := currentUsersQuery.ToSql()
	if err != nil {
		return 500, err
	}

	res, err := db.Query(sql)
	if err != nil {
		return 500, err
	}

	exists := res.Next()

	log.Info("are there existing users", "exists", exists)

	admin := false

	if !exists {
		admin = true
	}

	log.Info("admin is set to", "admin", admin)

	createUserRequest := createUserRequest{
		Username: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	}

	log.Info("create request populated", "username", createUserRequest.Username, "password", createUserRequest.Password)

	status, err = validateUsernameAndPassword(createUserRequest.Username, createUserRequest.Password)
	if err != nil {
		return status, err
	}

	log.Info("the user was successfully validated")

	hash, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), 14)
	if err != nil {
		return 500, err
	}

	log.Info("password was successfully hashed", "hash", string(hash))

	userID, err := createUserDBRecord(createUserRequest, string(hash), admin, db)
	if err != nil {
		return 500, err
	}

	log.Info("user record was created", "user_id", userID)

	sessionID := session.CreateSession(userID)

	log.Info("session was created", "session_id", sessionID)

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

	log.Info("create user query created", "query", sql, "args", args)

	user := types.User{}

	tx, err := db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	log.Info("transaction has started")

	rows, err := tx.Queryx(sql, args...)
	if err != nil {
		log.Error("failed to begin transaction", "err", err)
		return 0, err
	}
	defer rows.Close()

	log.Info("query was successful")

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.New("error getting returned user")
	}

	log.Info("user record", "user", user)

	createNewPermissionQuery := squirrel.
		Insert("Permission").
		Columns("UserID", "Admin", "ManageServices", "ManageTags", "ManageForwarding", "ViewLogs").
		Values(user.ID, admin, false, false, false, false).
		Suffix("RETURNING ID")

	sql, args, err = createNewPermissionQuery.ToSql()
	if err != nil {
		return 0, err
	}

	log.Info("create permission query", "query", sql, "args", args)

	permission := types.Permission{}

	rows, err = tx.Queryx(sql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	log.Info("permission request was successfully executed")

	if rows.Next() {
		err = rows.StructScan(&permission)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.New("error getting returned permission")
	}

	log.Info("permission was created", "permission", permission)

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

	log.Info("user was updated", "permission_id", permission.ID)

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	log.Info("tx was committed")

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
