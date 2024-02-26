package internal

import (
	"errors"
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

		snippetModel := models.SnippetModel{DB: app.DB}
		snippets, err := snippetModel.Latest()
		if err != nil {
			app.serverError(w, r, err)
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

		data := templateData{
			Snippets: snippets,
		}

		err = ts.ExecuteTemplate(w, "base", data)
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

		snippets := models.SnippetModel{DB: app.DB}
		snippet, err := snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
				return
			}

			app.serverError(w, r, err)
			return
		}

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/view.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		data := templateData{
			Snippet: snippet,
		}

		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			app.serverError(w, r, err)
		}
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
