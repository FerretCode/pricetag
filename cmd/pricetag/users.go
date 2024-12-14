package main

import (
	"net/http"
)

// Manage user permissions
func (app *application) handleUsersGet(w http.ResponseWriter, r *http.Request) error {
	selectedUsername := r.URL.Query().Get("user")
	if selectedUsername == "" {
		suid, err := app.getSessionUserID(r)
		if err != nil {
			return err
		}

		user, err := app.models.User.GetWithID(suid)
		if err != nil {
			return err
		}

		selectedUsername = user.Username
	}

	usernames, err := app.models.User.Usernames()
	if err != nil {
		return err
	}

	permissions, err := app.models.Permission.AllCodes()
	if err != nil {
		return err
	}

	selectdUser, err := app.models.User.GetWithUsername(selectedUsername)
	if err != nil {
		return err
	}

	userPerms, err := app.models.Permission.GetAllForUser(selectdUser.ID)
	if err != nil {
		return err
	}

	var data struct {
		Usernames   []string
		Selected    string
		Permissions []string
		UserPerms   []string
	}

	data.Usernames = usernames
	data.Selected = selectedUsername
	data.Permissions = permissions
	data.UserPerms = userPerms

	return app.render(w, r, http.StatusOK, "users.tmpl", data)
}

// Update user permissions
func (app *application) handleUsersPost(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	username := r.Form.Get("username")
	if username == "" {
		return app.renderError(w, r, http.StatusBadRequest, nil)
	}

	user, err := app.models.User.GetWithUsername(username)
	if err != nil {
		return err
	}

	permissions, err := app.models.Permission.AllCodes()
	if err != nil {
		return err
	}

	var codes []string
	for _, code := range permissions {
		if r.Form.Get(code) == "on" {
			codes = append(codes, code)
		}
	}

	err = app.models.Permission.SetForUser(user.ID, codes...)
	if err != nil {
		return err
	}

	app.refresh(w, r)

	return nil
}
