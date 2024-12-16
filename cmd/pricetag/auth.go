package main

import (
	"errors"
	"net/http"

	"github.com/ferretcode/pricetag/internal/models"
)

type contextKey string

const (
	authenticatedUserIDSessionKey = "authenticatedUserID"
	isAuthenticatedContextKey     = contextKey("isAuthenticated")
)

func (app *application) login(r *http.Request, userID int) error {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		return err
	}

	app.sessionManager.Put(r.Context(), authenticatedUserIDSessionKey, userID)

	return nil
}

func (app *application) logout(r *http.Request) error {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		return err
	}

	app.sessionManager.Remove(r.Context(), authenticatedUserIDSessionKey)

	return nil
}

// Check the auth context set by the authenticate middleware
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *application) getSessionUserID(r *http.Request) (int, error) {
	id := app.sessionManager.GetInt(r.Context(), authenticatedUserIDSessionKey)

	return id, nil
}

func (app *application) handleAuthLoginGet(w http.ResponseWriter, r *http.Request) error {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return nil
	}

	return app.render(w, r, http.StatusOK, "login.tmpl", nil)
}

func (app *application) handleAuthLoginPost(w http.ResponseWriter, r *http.Request) error {
	if app.isAuthenticated(r) {
		return app.renderError(w, r, http.StatusBadRequest, "already authenticated")
	}

	var form struct {
		Username string `form:"username" validate:"required"`
		Password string `form:"password" validate:"required"`
	}

	err := app.parseForm(r, &form)
	if err != nil {
		return err
	}

	user, err := app.models.User.GetForCredentials(form.Username, form.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidCredentials):
			return app.renderError(w, r, http.StatusUnauthorized, "invalid credentials")
		default:
			return err
		}
	}

	err = app.login(r, user.ID)
	if err != nil {
		return err
	}

	// Redirect to homepage after authenticating the user.
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func (app *application) handleAuthLogoutPost(w http.ResponseWriter, r *http.Request) error {
	err := app.logout(r)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func (app *application) handleAuthSignupPost(w http.ResponseWriter, r *http.Request) error {
	if app.isAuthenticated(r) {
		return app.renderError(w, r, http.StatusBadRequest, "already authenticated")
	}

	var form struct {
		Username string `form:"username" validate:"required,max=254"`
		Password string `form:"password" validate:"required,min=8,max=72"`
	}

	err := app.parseForm(r, &form)
	if err != nil {
		return err
	}

	makeAdmin, err := app.models.User.IsFirstUser()
	if err != nil {
		return err
	}

	user, err := app.models.User.New(form.Username, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUsername) {
			return app.renderError(w, r, http.StatusUnauthorized, "invalid credentials")
		}

		return err
	}

	if makeAdmin {
		err = app.models.Permission.SetForUser(user.ID, "admin")
		if err != nil {
			return err
		}
	}

	// Login user
	app.sessionManager.Clear(r.Context())
	err = app.login(r, user.ID)
	if err != nil {
		return err
	}

	f := FlashMessage{
		Type:    FlashSuccess,
		Message: "Successfully created account. Welcome!",
	}
	app.putFlash(r, f)

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func (app *application) handleAuthChangePasswordGet(w http.ResponseWriter, r *http.Request) error {
	return app.render(w, r, http.StatusOK, "change-password.tmpl", nil)
}

func (app *application) handleAuthChangePasswordPost(w http.ResponseWriter, r *http.Request) error {
	var form struct {
		Password string `form:"password" validate:"required,min=8,max=72"`
		Confirm  string `form:"confirm" validate:"required,eqfield=Password"`
	}

	err := app.parseForm(r, &form)
	if err != nil {
		return err
	}

	suid, err := app.getSessionUserID(r)
	if err != nil {
		return err
	}

	user, err := app.models.User.GetWithID(suid)
	if err != nil {
		return err
	}

	// Change password
	err = user.SetPasswordHash(form.Password)
	if err != nil {
		return err
	}

	err = app.models.User.Update(user)
	if err != nil {
		return err
	}

	// Logout user
	app.sessionManager.Clear(r.Context())
	err = app.logout(r)
	if err != nil {
		return err
	}

	f := FlashMessage{
		Type:    FlashSuccess,
		Message: "Successfully changed password. Please login.",
	}
	app.putFlash(r, f)

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}
