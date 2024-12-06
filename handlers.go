package main

import (
	"net/http"

	"github.com/ferretcode/pricetag/errors"
	"github.com/ferretcode/pricetag/routes/views"
	"github.com/go-chi/chi/v5"
)

func registerHandlers(r chi.Router) {
	r.Route("/dashboard", func(r chi.Router) {
		// TODO: write admin middleware
		r.Get("/home", func(w http.ResponseWriter, r *http.Request) {
			err := views.Home(w, r, templates)
			if err != nil {
				errors.HandleError(w, http.StatusInternalServerError, err.Error(), templates)
			}
		})
	})
}
