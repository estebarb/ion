/*
Ion is a small web framework for people in a hurry.
Ion provides a fast trie based router (julienschmidt/httprouter),
a integrated middleware support (justinas/alice), automatic
context support (via gorilla/context) and some handful helpers.

A short example:

	package main

	import (
		"fmt"
		"github.com/estebarb/ion"
		"github.com/gorilla/context"
		"github.com/julienschmidt/httprouter"
		"net/http"
	)

	func hello(w http.ResponseWriter, r *http.Request) {
		val := context.Get(r, ion.Urlargs).(httprouter.Params)
		if val != nil {
			fmt.Fprintf(w, "Hello, %v!", val.ByName("name"))
		} else {
			fmt.Fprint(w, "Hello world!")
		}
	}

	func main() {
		r := ion.NewRouter()
		r.GetFunc("/", hello)
		r.GetFunc("/:name", hello)
		http.ListenAndServe(":8080", r)
	}

At this point the framework is highly experimental, so please don't
use it in production for now...
*/
package ion

import (
	"github.com/gorilla/context"
	httprouter "github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"html/template"
	"net/http"
)

// This is the router used by Ion. It contains the high performance
// trie based httprouter, the middleware manager alice and adapter
// functions that make trivial to use the Go http.Handler
type Router struct {
	*httprouter.Router
	Middleware alice.Chain
}

/*
Ion adds the path arguments by httprouter to
context, so they can be retrieved using:

	// asuming a path like /:name
	// r is a *http.Request
	params := context.Get(r, ion.Urlargs)
	name := params.ByName("name")
*/
const Urlargs = "ion_urlargs"

/*
Returns a new router, with no middleware.
*/
func NewRouter() *Router {
	return &Router{httprouter.New(), alice.New()}
}

/*
Returns a new router, configured with the middlewares
provided.
*/
func NewRouterDefaults(middleware ...alice.Constructor) *Router {
	r := &Router{
		httprouter.New(),
		alice.New(middleware...),
	}
	return r
}

/*
wrapHandler transforms a http.Handler handler to a httprouter.Handle
*/
func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, Urlargs, ps)
		h.ServeHTTP(w, r)
	}
}

// Registers a new request handler (http.Handler) for the given path and method.
// It also executes the current Middleware in the settings.
func (r *Router) MethodHandle(method, path string, handle http.Handler) {
	r.Handle(method, path, wrapHandler(r.Middleware.Then(handle)))
}

// Registers a new request handler (http.HandlerFunc) for the given path and method.
// It also executes the current Middleware in the settings.
func (r *Router) MethodHandleFunc(method, path string, handle http.HandlerFunc) {
	r.Handle(method, path, wrapHandler(r.Middleware.ThenFunc(handle)))
}

// Shortcut for router.MethodHandle("DELETE", path, handler)
func (r *Router) DELETE(path string, handler http.Handler) {
	r.MethodHandle("DELETE", path, handler)
}

// Shortcut for router.MethodHandle("GET", path, handler)
func (r *Router) GET(path string, handler http.Handler) {
	r.MethodHandle("GET", path, handler)
}

// Shortcut for router.MethodHandle("POST", path, handler)
func (r *Router) POST(path string, handler http.Handler) {
	r.MethodHandle("POST", path, handler)
}

// Shortcut for router.MethodHandle("PATCH", path, handler)
func (r *Router) PATCH(path string, handler http.Handler) {
	r.MethodHandle("PATCH", path, handler)
}

// Shortcut for router.MethodHandle("PUT", path, handler)
func (r *Router) PUT(path string, handler http.Handler) {
	r.MethodHandle("PUT", path, handler)
}

// Shortcut for router.MethodHandleFunc("DELETE", path, handler)
func (r *Router) DeleteFunc(path string, handler http.HandlerFunc) {
	r.MethodHandleFunc("DELETE", path, handler)
}

// Shortcut for router.MethodHandleFunc("GET", path, handler)
func (r *Router) GetFunc(path string, handler http.HandlerFunc) {
	r.MethodHandleFunc("GET", path, handler)
}

// Shortcut for router.MethodHandleFunc("POST", path, handler)
func (r *Router) PostFunc(path string, handler http.HandlerFunc) {
	r.MethodHandleFunc("POST", path, handler)
}

// Shortcut for router.MethodHandleFunc("PATCH", path, handler)
func (r *Router) PatchFunc(path string, handler http.HandlerFunc) {
	r.MethodHandleFunc("PATCH", path, handler)
}

// Shortcut for router.MethodHandleFunc("PUT", path, handler)
func (r *Router) PutFunc(path string, handler http.HandlerFunc) {
	r.MethodHandleFunc("PUT", path, handler)
}

// Returns a handler that can render the given templates. The templates
// receives as parameters the context associated to the current
// request.
func RenderTemplate(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.GetAll(r)
		t.Execute(w, ctx)
	}
}
