package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ferretcode/pricetag/ui"
	"github.com/go-chi/chi/v5"
)

// App router
func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(app.recovery)
	r.Use(secureHeaders)

	// Static files
	r.Handle("/static/*", app.handleStatic())
	r.Get("/favicon.ico", app.handleFavicon)

	r.Route("/", func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave)
		r.Use(app.noSurf)
		r.Use(app.authenticate)

		r.Route("/auth", func(r chi.Router) {
			r.Get("/login", app.handle(app.handleAuthLoginGet))
			r.Post("/login", app.handle(app.handleAuthLoginPost))
			r.Post("/logout", app.handle(app.handleAuthLogoutPost))
			r.Post("/signup", app.handle(app.handleAuthSignupPost))
			r.Route("/change-password", func(r chi.Router) {
				r.Use(app.requireAuthentication)
				r.Get("/", app.handle(app.handleAuthChangePasswordGet))
				r.Post("/", app.handle(app.handleAuthChangePasswordPost))
			})
		})

		r.Route("/", func(r chi.Router) {
			r.Use(app.requireAuthentication)

			r.Get("/", app.handle(app.handleDashboardGet))
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(app.requirePermission("admin"))

			r.Get("/", app.handle(app.handleUsersGet))
			r.Post("/", app.handle(app.handleUsersPost))
			r.Post("/delete", app.handle(app.handleUsersDeletePost))
		})
	})

	r.NotFound(app.handle(app.handleNotFound))

	return r
}

func (app *application) refresh(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *application) handleStatic() http.Handler {
	if app.config.dev {
		fs := http.FileServer(http.Dir("./ui/static/"))

		return http.StripPrefix("/static", fs)
	}

	return http.FileServer(http.FS(ui.Files))
}

func (app *application) handleFavicon(w http.ResponseWriter, r *http.Request) {
	if app.config.dev {
		http.ServeFile(w, r, "./ui/static/favicon.ico")

		return
	}
	http.ServeFileFS(w, r, ui.Files, "static/favicon.ico")
}

func (app *application) handleNotFound(w http.ResponseWriter, r *http.Request) error {
	return app.renderError(w, r, http.StatusNotFound, MessageNotFound)
}

type withError func(w http.ResponseWriter, r *http.Request) error

// http.HandlerFunc wrapper with error handling
func (app *application) handle(h withError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// First, check for any expected errors
			var formErrors FormErrors
			if errors.As(err, &formErrors) {
				// Redirect to previous page with form errors in session data
				app.putFormErrors(r, formErrors)
				http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
			}

			// Log unexpected error and return internal server error
			app.logger.Error("handled unexpected error", slog.Any("err", err), slog.String("type", fmt.Sprintf("%T", err)))

			http.Error(w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			app.renderError(w, r, http.StatusInternalServerError, "something went wrong...")
		}
	}
}
