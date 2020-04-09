package main

import (
	"fmt"
	"github.com/estebarb/ion"
	"net/http"
)

func newApp() http.Handler {
	routes := ion.Routes{
		"/:name": {
			Middleware:  []ion.Middleware{ion.PathEnd},
			HttpHandler: http.HandlerFunc(hello),
		},
	}
	return routes.Build()
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := r.Context().Value("name")
	if name != "" {
		fmt.Fprintf(w, "Hello, %v!", name)
	} else {
		fmt.Fprint(w, "Hello world!")
	}
}

func main() {
	app := newApp()
	http.ListenAndServe(":5500", app)
}
