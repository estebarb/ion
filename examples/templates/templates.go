package main

import (
	"github.com/estebarb/ion"
	"net/http"
	"strings"
	"github.com/estebarb/ion/components/reqctx"
)

type App struct {
	*ion.Ion
	States *reqctx.StateContainer
}

func (app *App) hello(w http.ResponseWriter, r *http.Request) {
	state := app.States.GetState(r)
	state.Set("key", "value")
	app.Template.RenderTemplate("hello").ServeHTTP(w, r)
}

func main() {
	// Creates a new router
	app := &App{
		Ion: ion.New(),
		States: reqctx.NewStateContainer(),
	}

	app.Template.AddStateContainer("App", app.States)
	app.Template.AddStateContainer("Router", app.Router.GetStateContainer())

	// Adds the functions used in the templates.
	// Other packages can add their functions (eg:
	// user authentication, etc)
	app.Template.AddFunc("title", strings.Title)
	// Loads the templates
	app.Template.LoadPattern("views/*.html")

	// The usual Ion handlers
	handler := app.States.Middleware(http.HandlerFunc(app.hello))
	app.Get("/", handler)
	app.Get("/:name", handler)

	// Starts the server.
	http.ListenAndServe(":8000", app)
}
