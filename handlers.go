package main

import (
	"database/sql"
	"net/http"

	"github.com/ferretcode/pricetag/errors"
	"github.com/ferretcode/pricetag/routes/user"
	"github.com/ferretcode/pricetag/routes/views"
	"github.com/go-chi/chi/v5"
)

func registerHandlers(r chi.Router, db *sql.DB) {
	r.Route("/dashboard", func(r chi.Router) {
		// TODO: write admin middleware
		r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
			err := views.Home(w, r, templates)
			if err != nil {
				errors.HandleError(w, http.StatusInternalServerError, err.Error(), templates)
			}
		})
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/create", func(w http.ResponseWriter, r *http.Request) {
			err := user.RenderCreateUserPage(w, r, templates)
			if err != nil {
				errors.HandleError(w, http.StatusInternalServerError, err.Error(), templates)
			}
		})

		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			status, err := user.Create(w, r, db)
			if err != nil {
				errors.HandleError(w, status, err.Error(), templates)
			}
		})

		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {})

		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {

		})
	})
}
