package errors

import (
	"html/template"
	"net/http"

	"github.com/ferretcode/pricetag/types"
)

func HandleError(w http.ResponseWriter, status int, err string, templates *template.Template) {
	serveErr := templates.ExecuteTemplate(w, "error.html", types.Error{
		Status: status,
		Error:  err,
	})

	if serveErr != nil {
		http.Error(w, serveErr.Error(), http.StatusInternalServerError)
	}
}
