package middleware

import (
	"context"
	"html/template"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/ferretcode/pricetag/errors"
	"github.com/ferretcode/pricetag/session"
	"github.com/ferretcode/pricetag/types"
	"github.com/jmoiron/sqlx"
)

func CheckUser(db *sqlx.DB, session *session.SessionManager, templates *template.Template) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")

			if err != nil {
				if err == http.ErrNoCookie {
					errors.HandleError(w, r.URL.Path, 403, "please log in", templates)
					return
				}
			}

			session, err := session.GetSession(cookie.Value)
			if err != nil {
				errors.HandleError(w, r.URL.Path, 403, err.Error(), templates)
				return
			}

			user := types.User{}
			permission := types.Permission{}

			selectUserQuery := squirrel.
				Select("*").
				From("User").
				Where(squirrel.Eq{"ID": session.UserID})

			sql, args, err := selectUserQuery.ToSql()
			if err != nil {
				errors.HandleError(w, r.URL.Path, 500, err.Error(), templates)
				return
			}

			err = db.Get(&user, sql, args...)
			if err != nil {
				errors.HandleError(w, r.URL.Path, 403, "you may not access this resource", templates)
				return
			}

			selectPermissionQuery := squirrel.
				Select("*").
				From("Permission").
				Where(squirrel.Eq{"UserID": session.UserID})

			sql, args, err = selectPermissionQuery.ToSql()
			if err != nil {
				errors.HandleError(w, r.URL.Path, 500, err.Error(), templates)
				return
			}

			err = db.Get(&permission, sql, args...)
			if err != nil {
				errors.HandleError(w, r.URL.Path, 403, "you may not access this resource", templates)
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			ctx = context.WithValue(ctx, "permission", permission)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
