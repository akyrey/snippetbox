package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akyrey/snippetbox/internal/models"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	snippetModel := models.SnippetModel{DB: app.DB}
	snippets, err := snippetModel.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := NewTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", &data)
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
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

	data := NewTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", &data)
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData(r)

	app.render(w, r, http.StatusOK, "create.tmpl", &data)
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	snippets := models.SnippetModel{DB: app.DB}
	id, err := snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
