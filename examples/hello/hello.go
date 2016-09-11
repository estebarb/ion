package main

import (
	"fmt"
	"github.com/estebarb/ion"
	"net/http"
)

type App struct {
	*ion.Ion
}

func NewApp() *App {
	app := &App{
		Ion: ion.New(),
	}
	app.GetFunc("/", app.hello)
	app.GetFunc("/:name", app.hello)
	return app
}

func (app *App) hello(w http.ResponseWriter, r *http.Request) {
	state := app.Router.GetState(r)
	value, exists := state.Get("name")
	if exists {
		fmt.Fprintf(w, "Hello, %v!", value)
	} else {
		fmt.Fprint(w, "Hello world!")
	}
}

func main() {
	http.ListenAndServe(":5500", NewApp())
}
