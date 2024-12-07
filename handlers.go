package main

import (
	"net/http"

	"github.com/ferretcode/pricetag/errors"
	"github.com/ferretcode/pricetag/middleware"
	"github.com/ferretcode/pricetag/routes/dashboard"
	"github.com/ferretcode/pricetag/routes/user"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func registerHandlers(r chi.Router, db *sqlx.DB) {
	r.Route("/dashboard", func(r chi.Router) {
		r.Use(middleware.CheckAdmin(db, sessionManager, templates))

		r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
			err := dashboard.Home(w, r, templates)
			if err != nil {
				errors.HandleError(w, "/dashboard/home", http.StatusInternalServerError, err.Error(), templates)
			}
		})
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/create", func(w http.ResponseWriter, r *http.Request) {
			err := user.RenderCreateUserPage(w, r, templates)
			if err != nil {
				errors.HandleError(w, "GET /user/create", http.StatusInternalServerError, err.Error(), templates)
			}
		})

		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			status, err := user.Create(w, r, db, sessionManager)
			if err != nil {
				errors.HandleError(w, "POST /user/create", status, err.Error(), templates)
			}
		})

		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			err := user.RenderLoginPage(w, r, templates)
			if err != nil {
				errors.HandleError(w, "GET /user/login", http.StatusInternalServerError, err.Error(), templates)
			}
		})

		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			status, err := user.Login(w, r, db, sessionManager)
			if err != nil {
				errors.HandleError(w, "POST /user/login", status, err.Error(), templates)
			}
		})
	})
}
