package errors

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/ferretcode/pricetag/types"
)

func HandleError(w http.ResponseWriter, source string, status int, err string, templates *template.Template) {
	log.Error(
		fmt.Sprintf("error serving %s", source),
		"err",
		err,
	)

	serveErr := templates.ExecuteTemplate(w, "error.html", types.Error{
		Status: status,
		Error:  err,
	})

	if serveErr != nil {
		http.Error(w, serveErr.Error(), http.StatusInternalServerError)
	}
}
