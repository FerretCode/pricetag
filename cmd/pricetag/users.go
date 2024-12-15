package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ferretcode/pricetag/internal/models"
)

// Manage user permissions
func (app *application) handleUsersGet(w http.ResponseWriter, r *http.Request) error {
	suid, err := app.getSessionUserID(r)
	if err != nil {
		return err
	}

	sessionUser, err := app.models.User.GetWithID(suid)
	if err != nil {
		return err
	}

	selectedUsername := r.URL.Query().Get("username")
	if selectedUsername == "" {
		selectedUsername = sessionUser.Username
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

	userPerms, err := app.models.Permission.GetForUser(selectdUser.ID)
	if err != nil {
		return err
	}

	var data struct {
		Usernames   []string
		Permissions []string
		Username    string
		UserPerms   []string
		ShowDelete  bool
	}

	data.Username = selectedUsername
	data.Usernames = usernames
	data.Permissions = permissions
	data.UserPerms = userPerms
	// Can not delete admin unless self
	data.ShowDelete = !userPerms.Include("admin") || selectedUsername == sessionUser.Username

	return app.render(w, r, http.StatusOK, "users.tmpl", data)
}

// Update user permissions
func (app *application) handleUsersPost(w http.ResponseWriter, r *http.Request) error {
	var form struct {
		Username string `form:"username" validate:"required"`
	}

	err := app.parseForm(r, &form)
	if err != nil {
		return err
	}

	user, err := app.models.User.GetWithUsername(form.Username)
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

func (app *application) handleUsersDeletePost(w http.ResponseWriter, r *http.Request) error {
	var form struct {
		Username string `form:"username" validate:"required"`
	}

	err := app.parseForm(r, &form)
	if err != nil {
		return err
	}

	user, err := app.models.User.GetWithUsername(form.Username)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoRecord):
			return app.renderError(w, r, http.StatusBadRequest, "user does not exist")
		default:
			return err
		}
	}

	perms, err := app.models.Permission.GetForUser(user.ID)
	if err != nil {
		return err
	}

	if perms.Include("admin") {
		suid, err := app.getSessionUserID(r)
		if err != nil {
			return err
		}

		if user.ID != suid {
			return app.renderError(w, r, http.StatusForbidden, "can not delete admin user")
		}
	}

	err = app.models.User.Delete(user.ID)
	if err != nil {
		return err
	}

	f := FlashMessage{
		Type:    FlashSuccess,
		Message: fmt.Sprintf("Successfully deleted user: %s", user.Username),
	}
	app.putFlash(r, f)

	http.Redirect(w, r, "/users", http.StatusSeeOther)

	return nil
}
