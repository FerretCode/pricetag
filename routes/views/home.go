package views

import (
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request, templates *template.Template) error {
	err := templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		return err
	}

	return nil
}
