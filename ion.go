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
		"net/http"
	)

	func hello(w http.ResponseWriter, r *http.Request) {
		value := ion.URLArgs(r, "name")
		if value != "" {
			fmt.Fprintf(w, "Hello, %v!", value)
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
	"github.com/estebarb/ion/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"html/template"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"time"
	"log"
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

But the preferred way is the following:

	name := ion.URLArgs(r, "name")

*/
const urlargs = "urlargs"

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
		context.Set(r, urlargs, ps)
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
		// Adds URL params to context
		args, _ := context.Get(r, urlargs)
		context.Set(r, urlargs, paramsToMap(args.(httprouter.Params)))
		// Adds data in context to template context
		ctx, _ := context.GetAll(r)
		t.Execute(w, ctx)
	}
}

// Converts the httprouter.Params array to a map, that can
// be consumed easily from templates.
func paramsToMap(p httprouter.Params) map[string]string {
	ret := make(map[string]string)
	for _, v := range p {
		ret[v.Key] = v.Value
	}
	return ret
}

// This interface works with RegisterREST to provide a shortcut
// to register an RESTful endpoint.

type RESTendpoint interface {
	LIST(w http.ResponseWriter, r *http.Request)
	POST(w http.ResponseWriter, r *http.Request)
	PUT(w http.ResponseWriter, r *http.Request)
	GET(w http.ResponseWriter, r *http.Request)
	DELETE(w http.ResponseWriter, r *http.Request)
}

// Register a RESTendpoint in a router.
// It will register the following routes:
// - GET  path		(list function)
// - POST path		(post function)
// - GET  path/:id	(get function)
// - PUT  path/:id	(put function)
// - DELETE  path/:id	(delete function)
// The path MUST include the trailing slash.
func (r *Router) RegisterREST(path string, handler RESTendpoint) {
	r.GetFunc(path+":id", handler.GET)
	r.PutFunc(path+":id", handler.PUT)
	r.DeleteFunc(path+":id", handler.DELETE)
	r.GetFunc(path, handler.LIST)
	r.PostFunc(path, handler.POST)
}

// Returns a named argument from the request URL
func URLArgs(r *http.Request, name string) string {
	val, _ := context.Get(r, urlargs)
	v2 := val.(httprouter.Params)
	return v2.ByName(name)
}

// Writes a JSON value
func MarshalJSON(w http.ResponseWriter, value interface{}) error{
	out, err := json.Marshal(value)
	if err != nil{
		return err		
	}
	w.Write(out)
	return nil
}

// Returns the unmarshaled value from the request body
func UnmarshalJSON(r *http.Request, value interface{}) error{
	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		return err
	}
	return json.Unmarshal(body, value)
}

// Provides a logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

// Provides a recovery from panics in other handlers
func PanicMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}