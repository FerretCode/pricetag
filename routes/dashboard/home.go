package dashboard

import (
	"html/template"
	"net/http"

	"github.com/ferretcode/pricetag/types"
)

type homeData struct {
	User       types.User
	Permission types.Permission
}

func Home(w http.ResponseWriter, r *http.Request, templates *template.Template) error {
	err := templates.ExecuteTemplate(w, "home.html", homeData{
		User:       r.Context().Value("user").(types.User),
		Permission: r.Context().Value("permission").(types.Permission),
	})
	if err != nil {
		return err
	}

	return nil
}
