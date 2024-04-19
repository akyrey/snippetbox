package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akyrey/snippetbox/internal"
	"github.com/akyrey/snippetbox/internal/models"
	"github.com/akyrey/snippetbox/internal/validator"
)

// snippetCreateForm is a struct that represents the form fields for creating a new snippet.
// All fields must be exported so that the template will be able to access them.
type snippetCreateForm struct {
	validator.Validator `form:"-"`
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
}

func (app *internal.Application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.Snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.NewTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", &data)
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.Snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}

		app.serverError(w, r, err)
		return
	}

	data := app.NewTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", &data)
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl", &data)
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(
		validator.NotBlank(form.Title),
		"title",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.MaxChars(form.Title, 100),
		"title",
		"This field cannot be more than 100 characters long",
	)
	form.CheckField(
		validator.NotBlank(form.Content),
		"content",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.PermittedValue(form.Expires, 1, 7, 365),
		"expires",
		"This field must equal 1, 7 or 365",
	)

	if !form.IsValid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", &data)
		return
	}

	id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.SessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
