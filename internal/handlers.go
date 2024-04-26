package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akyrey/snippetbox/internal/models"
	"github.com/akyrey/snippetbox/internal/validator"
)

// snippetCreateForm is a struct that represents the form fields for creating a new snippet.
// All fields must be exported so that the template will be able to access them.
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	validator.Validator `form:"-"`
	Expires             int `form:"expires"`
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", &data)
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(
		validator.NotBlank(form.Name),
		"name",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.MaxChars(form.Name, 255),
		"name",
		"This field cannot be more than 255 characters long",
	)
	form.CheckField(
		validator.NotBlank(form.Email),
		"email",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.MaxChars(form.Email, 255),
		"email",
		"This field cannot be more than 255 characters long",
	)
	form.CheckField(
		validator.Matches(form.Email, validator.EmailRegExp),
		"email",
		"This field must be a valid email address",
	)
	form.CheckField(
		validator.NotBlank(form.Password),
		"password",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.MinChars(form.Password, 8),
		"password",
		"This field cannot be less than 8 characters long",
	)

	if !form.IsValid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", &data)
		return
	}

	err = app.Users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.NewTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", &data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.SessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", &data)
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(
		validator.NotBlank(form.Email),
		"email",
		"This field cannot be blank",
	)
	form.CheckField(
		validator.Matches(form.Email, validator.EmailRegExp),
		"email",
		"This field must be a valid email address",
	)
	form.CheckField(
		validator.NotBlank(form.Password),
		"password",
		"This field cannot be blank",
	)

	if !form.IsValid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", &data)
		return
	}

	id, err := app.Users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.NewTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", &data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	err = app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.SessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.SessionManager.Remove(r.Context(), "authenticatedUserID")
	app.SessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
