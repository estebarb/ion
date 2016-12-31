/*
Package ion is a small web framework written in Go.

Ion provides a request router, middleware management, automatic
context support and some handful helpers.

A short example:

	package main

	import (
		"fmt"
		"github.com/estebarb/ion"
		"github.com/estebarb/ion/components/router"
		"net/http"
	)

	func hello(w http.ResponseWriter, r *http.Request) {
		context := ion.App.Context(r)
		params := context.(router.IPathParam)
		value, exists := params.PathParams()["name"]
		if exists {
			fmt.Fprintf(w, "Hello, %v!", value)
		} else {
			fmt.Fprint(w, "Hello world!")
		}
	}

	func main() {
		ion.GetFunc("/", hello)
		ion.GetFunc("/:name", hello)
		http.ListenAndServe(":5500", ion.App)
	}

At this point the framework is highly experimental, so please take that
in consideration if you want to use it.
*/
package ion

import (
	"encoding/json"
	"github.com/estebarb/ion/components/chain"
	"github.com/estebarb/ion/components/router"
	"io/ioutil"
	"net/http"
)

// Ion represents an Ion web application
type Ion struct {
	*router.Router
	Middleware []*chain.Chain
}

// App is the default Ion application created when importing the package.
// Instead of creating and passing a ion.Ion instance in the application
// this global variable could be used.
var App *Ion

func init() {
	App = New(router.ContextFactory)
}

// ServeHTTP calls ion.App.ServeHTTP(w, r), the request dispatcher
// of the default Ion application.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	App.ServeHTTP(w, r)
}

// MethodHandle registers a request handler for the given path, and adds the current
// middleware in Ion settings.
func MethodHandle(method, path string, handle http.Handler) *router.Route {
	return App.Handler(method, path, App.generateHandler(handle))
}

// MethodHandleFunc registers a request handler for the given path, and adds the current
// middleware in Ion settings.
func MethodHandleFunc(method, path string, handle http.HandlerFunc) *router.Route {
	return App.Handler(method, path, App.generateHandlerFunc(handle))
}

// Delete register a Delete action in the default application
func Delete(path string, handler http.Handler) *router.Route {
	return App.Delete(path, App.generateHandler(handler))
}

// Get register a Get action in the default application
func Get(path string, handler http.Handler) *router.Route {
	return App.Get(path, App.generateHandler(handler))
}

// Post register a Post action in the default application
func Post(path string, handler http.Handler) *router.Route {
	return App.Post(path, App.generateHandler(handler))
}

// Patch register a Patch action in the default application
func Patch(path string, handler http.Handler) *router.Route {
	return App.Patch(path, App.generateHandler(handler))
}

// Put register a Put action in the default application
func Put(path string, handler http.Handler) *router.Route {
	return App.Put(path, App.generateHandler(handler))
}

// DeleteFunc register a DeleteFunc action in the default application
func DeleteFunc(path string, handler http.HandlerFunc) *router.Route {
	return App.Delete(path, App.generateHandlerFunc(handler))
}

// GetFunc register a GetFunc action in the default application
func GetFunc(path string, handler http.HandlerFunc) *router.Route {
	return App.Get(path, App.generateHandlerFunc(handler))
}

// PostFunc register a PostFunc action in the default application
func PostFunc(path string, handler http.HandlerFunc) *router.Route {
	return App.Post(path, App.generateHandlerFunc(handler))
}

// PatchFunc register a PatchFunc action in the default application
func PatchFunc(path string, handler http.HandlerFunc) *router.Route {
	return App.Patch(path, App.generateHandlerFunc(handler))
}

// PutFunc register a PutFunc action in the default application
func PutFunc(path string, handler http.HandlerFunc) *router.Route {
	return App.Put(path, App.generateHandlerFunc(handler))
}

// New creates an Ion application, with router and middleware support.
// It receives a contextFactory that creates a context per request.
func New(contextFactory func() interface{}) *Ion {
	app := &Ion{
		Router:     router.New(contextFactory),
		Middleware: []*chain.Chain{chain.New()},
	}
	return app
}

// generateHandler wraps the given Handler with the middleware layers
func (a *Ion) generateHandler(handle http.Handler) http.Handler {
	return chain.Join(a.Middleware...).Then(handle)
}

// generateHandler wraps the given HandlerFunc with the middleware layers
func (a *Ion) generateHandlerFunc(handle http.HandlerFunc) http.Handler {
	return chain.Join(a.Middleware...).ThenFunc(handle)
}

// MethodHandle registers a request handler for the given path, and adds the current
// middleware in Ion settings.
func (a *Ion) MethodHandle(method, path string, handle http.Handler) *router.Route {
	return a.Router.Handler(method, path, a.generateHandler(handle))
}

// MethodHandleFunc registers a request handler for the given path, and adds the current
// middleware in Ion settings.
func (a *Ion) MethodHandleFunc(method, path string, handle http.HandlerFunc) *router.Route {
	return a.Router.Handler(method, path, a.generateHandlerFunc(handle))
}

// Delete register the handler in the router, after wrapping it with the middleware
func (a *Ion) Delete(path string, handler http.Handler) *router.Route {
	return a.Router.Delete(path, a.generateHandler(handler))
}

// Get register the handler in the router, after wrapping it with the middleware
func (a *Ion) Get(path string, handler http.Handler) *router.Route {
	return a.Router.Get(path, a.generateHandler(handler))
}

// Post register the handler in the router, after wrapping it with the middleware
func (a *Ion) Post(path string, handler http.Handler) *router.Route {
	return a.Router.Post(path, a.generateHandler(handler))
}

// Patch register the handler in the router, after wrapping it with the middleware
func (a *Ion) Patch(path string, handler http.Handler) *router.Route {
	return a.Router.Patch(path, a.generateHandler(handler))
}

// Put register the handler in the router, after wrapping it with the middleware
func (a *Ion) Put(path string, handler http.Handler) *router.Route {
	return a.Router.Put(path, a.generateHandler(handler))
}

// DeleteFunc register the handler in the router, after wrapping it with the middleware
func (a *Ion) DeleteFunc(path string, handler http.HandlerFunc) *router.Route {
	return a.Router.Delete(path, a.generateHandlerFunc(handler))
}

// GetFunc register the handler in the router, after wrapping it with the middleware
func (a *Ion) GetFunc(path string, handler http.HandlerFunc) *router.Route {
	return a.Router.Get(path, a.generateHandlerFunc(handler))
}

// PostFunc register the handler in the router, after wrapping it with the middleware
func (a *Ion) PostFunc(path string, handler http.HandlerFunc) *router.Route {
	return a.Router.Post(path, a.generateHandlerFunc(handler))
}

// PatchFunc register the handler in the router, after wrapping it with the middleware
func (a *Ion) PatchFunc(path string, handler http.HandlerFunc) *router.Route {
	return a.Router.Patch(path, a.generateHandlerFunc(handler))
}

// PutFunc register the handler in the router, after wrapping it with the middleware
func (a *Ion) PutFunc(path string, handler http.HandlerFunc) *router.Route {
	return a.Router.Put(path, a.generateHandlerFunc(handler))
}

// RESTendpoint works with RegisterREST to provide a shortcut
// to register an RESTful endpoint.
type RESTendpoint interface {
	LIST(w http.ResponseWriter, r *http.Request)
	POST(w http.ResponseWriter, r *http.Request)
	PUT(w http.ResponseWriter, r *http.Request)
	GET(w http.ResponseWriter, r *http.Request)
	DELETE(w http.ResponseWriter, r *http.Request)
}

// RegisterREST register a RESTendpoint in the router with some default paths
// It will register the following routes:
//
//    - GET  path		(list function)
//    - POST path		(post function)
//    - GET  path/:id	        (get function)
//    - PUT  path/:id	        (put function)
//    - DELETE  path/:id	(delete function)
//
// The path MUST include the trailing slash.
func (a *Ion) RegisterREST(path string, handler RESTendpoint) {
	a.GetFunc(path+":id", handler.GET)
	a.PutFunc(path+":id", handler.PUT)
	a.DeleteFunc(path+":id", handler.DELETE)
	a.GetFunc(path, handler.LIST)
	a.PostFunc(path, handler.POST)
}

// DoNothing is a handler that do nothing.
// It can be used as a placeholder, or when we want to run
// the middleware for the collateral effects, but we don't want
// to do something specific.
// Also can be used with IonMVC.
func DoNothing(w http.ResponseWriter, r *http.Request) {
}

// MarshalJSON writes a JSON value to the ResponseWriter
func MarshalJSON(w http.ResponseWriter, value interface{}) error {
	out, err := json.Marshal(value)
	if err != nil {
		return err
	}
	w.Write(out)
	return nil
}

// UnmarshalJSON parses the request body as a JSON value
func UnmarshalJSON(r *http.Request, value interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, value)
}
