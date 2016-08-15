package main

import (
	"fmt"
	"net/http"
	"github.com/estebarb/ion/components/router"
)

type App struct {
	*router.Router
}

func NewApp() *App {
	app := &App{
		Router: router.New(),
	}
	app.GetFunc("/", app.hello)
	app.GetFunc("/:name", app.hello)
	return app
}

func (app *App) hello(w http.ResponseWriter, r *http.Request) {
	state := app.GetState(r)
	value, exists := state.Get("name")
	if exists {
		fmt.Fprintf(w, "Hello, %v!", value)
	} else {
		fmt.Fprint(w, "Hello world!")
	}
}

func main() {
	http.ListenAndServe(":5000", NewApp())
}
