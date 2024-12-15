package main

import (
	"net/http"
	"strings"
)

func (app *application) handleDashboardGet(w http.ResponseWriter, r *http.Request) error {
	suid, err := app.getSessionUserID(r)
	if err != nil {
		return err
	}

	user, err := app.models.User.GetWithID(suid)
	if err != nil {
		return err
	}

	perms, err := app.models.Permission.GetForUser(user.ID)
	if err != nil {
		return err
	}

	var data struct {
		Username    string
		Permissions string
	}

	data.Username = user.Username
	data.Permissions = strings.Join(perms, ", ")

	return app.render(w, r, http.StatusOK, "dashboard.tmpl", data)
}
