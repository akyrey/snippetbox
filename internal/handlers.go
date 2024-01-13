package internal

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/akyrey/snippetbox/internal/models"
)

func (app *Application) HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/home.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			app.serverError(w, r, err)
		}
	}
}

func (app *Application) SnippetViewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.notFound(w)
			return
		}

		fmt.Fprintf(w, "Display a specific snippet with ID %d", id)
	}
}

func (app *Application) SnipperCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			app.clientError(w, http.StatusMethodNotAllowed)
			return
		}

		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
		expires := 7

		snippets := models.SnippetModel{DB: app.DB}
		id, err := snippets.Insert(title, content, expires)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	}
}
