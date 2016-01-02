package main

import (
	"github.com/estebarb/ion"
	"github.com/estebarb/ion/context"
	"github.com/estebarb/ion/ionvc"
	"net/http"
	"strings"
)

func hello(w http.ResponseWriter, r *http.Request) {
	context.Set(r, "key", "value")
}

func main() {
	// Creates a new router
	r := ion.NewRouter()

	// Adds the functions used in the templates.
	// Other packages can add their functions (eg:
	// user authentication, etc)
	ionvc.AddFunc("title", strings.Title)
	// Loads the templates
	ionvc.LoadTemplates("views/*.html")

	// The usual Ion handlers, but using ionvc.ControllerFunc
	// wrapper. The first parameter is the handler, and
	// the second is the name of the template.
	r.GetFunc("/", ionvc.ControllerFunc(hello, "hello"))
	r.GetFunc("/{name}", ionvc.ControllerFunc(hello, "hello"))

	// Starts the server.
	http.ListenAndServe(":8080", r)
}
