package internal

import (
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/akyrey/snippetbox/internal/models"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/justinas/alice"
)

type Application struct {
	DB             *sql.DB
	FormDecoder    *form.Decoder
	Logger         *slog.Logger
	SessionManager *scs.SessionManager
	Snippets       *models.SnippetModel
	TemplateCache  map[string]*template.Template
	Users          *models.UserModel
}

func (app *Application) Routes(config Config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(config.StaticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.SessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// logRequest ↔ secureHeaders ↔ servemux ↔ application handler
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)
}
