/*
Ion is a small web framework for people in a hurry.
Ion provides a flexible router (gorilla/mux),
a integrated middleware support (justinas/alice), automatic
context support and some handful helpers.

A short example:

	package main

	import (
		"fmt"
		"github.com/estebarb/ion"
		"net/http"
	)

	type App struct {
		*ion.Ion
	}

	func (app *App) hello(w http.ResponseWriter, r *http.Request) {
		state := app.Router.GetState(r)
	    value, exists := state.Get("name")
	    if exists
			fmt.Fprintf(w, "Hello, %v!", value)
		} else {
			fmt.Fprint(w, "Hello world!")
		}
	}

	func main() {
		app := &App{
			Ion: ion.New(),
		}
		app.GetFunc("/", app.hello)
		app.GetFunc("/:name", app.hello)
		http.ListenAndServe(":5500", app)
	}


At this point the framework is highly experimental, so please take that
in consideration if you want to use it.
*/
package ion

import (
	"encoding/json"
	"github.com/estebarb/ion/components/chain"
	"github.com/estebarb/ion/components/router"
	"github.com/estebarb/ion/components/templates"
	"io/ioutil"
	"net/http"
)

// Ion represents an Ion web application
type Ion struct {
	Router     *router.Router
	Middleware []*chain.Chain
	Template   *templates.Templates
}

var App *Ion

func init(){
	App = New()
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	App.Router.ServeHTTP(w, r)
}

// MethodHandle registers a request handler for the given path, and adds the current
// middleware in Ion settings.
func MethodHandle(method, path string, handle http.Handler) *router.Route {
	return App.Router.Handler(method, path, App.generateHandler(handle))
}

// MethodHandleFunc registers a request handler for the given path, and adds the current
// middleware in Ion settings.
func MethodHandleFunc(method, path string, handle http.HandlerFunc) *router.Route {
	return App.Router.Handler(method, path, App.generateHandlerFunc(handle))
}

func Delete(path string, handler http.Handler) *router.Route {
	return App.Router.Delete(path, App.generateHandler(handler))
}

func Get(path string, handler http.Handler) *router.Route {
	return App.Router.Get(path, App.generateHandler(handler))
}

func Post(path string, handler http.Handler) *router.Route {
	return App.Router.Post(path, App.generateHandler(handler))
}

func Patch(path string, handler http.Handler) *router.Route {
	return App.Router.Patch(path, App.generateHandler(handler))
}

func Put(path string, handler http.Handler) *router.Route {
	return App.Router.Put(path, App.generateHandler(handler))
}

func DeleteFunc(path string, handler http.HandlerFunc) *router.Route {
	return App.Router.Delete(path, App.generateHandlerFunc(handler))
}

func GetFunc(path string, handler http.HandlerFunc) *router.Route {
	return App.Router.Get(path, App.generateHandlerFunc(handler))
}

func PostFunc(path string, handler http.HandlerFunc) *router.Route {
	return App.Router.Post(path, App.generateHandlerFunc(handler))
}

func PatchFunc(path string, handler http.HandlerFunc) *router.Route {
	return App.Router.Patch(path, App.generateHandlerFunc(handler))
}

func PutFunc(path string, handler http.HandlerFunc) *router.Route {
	return App.Router.Put(path, App.generateHandlerFunc(handler))
}



/*
Returns a new router, with no middleware.
*/
func New() *Ion {
	return &Ion{
		Router:     router.New(),
		Middleware: []*chain.Chain{chain.New()},
		Template:   templates.New(),
	}
}

func (a *Ion) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Router.ServeHTTP(w, r)
}

func (a *Ion) generateHandler(handle http.Handler) http.Handler {
	return chain.Join(a.Middleware...).Then(handle)
}

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

func (a *Ion) Delete(path string, handler http.Handler) *router.Route {
	return a.Router.Delete(path, a.generateHandler(handler))
}

func (a *Ion) Get(path string, handler http.Handler) *router.Route {
	return a.Router.Get(path, a.generateHandler(handler))
}

func (a *Ion) Post(path string, handler http.Handler) *router.Route {
	return a.Router.Post(path, a.generateHandler(handler))
}

func (a *Ion) Patch(path string, handler http.Handler) *router.Route {
	return a.Router.Patch(path, a.generateHandler(handler))
}

func (a *Ion) Put(path string, handler http.Handler) *router.Route {
	return a.Router.Put(path, a.generateHandler(handler))
}

func (a *Ion) DeleteFunc(path string, handler http.HandlerFunc) *router.Route {
	return a.Router.Delete(path, a.generateHandlerFunc(handler))
}

func (a *Ion) GetFunc(path string, handler http.HandlerFunc) *router.Route {
	return a.Router.Get(path, a.generateHandlerFunc(handler))
}

func (a *Ion) PostFunc(path string, handler http.HandlerFunc) *router.Route {
	return a.Router.Post(path, a.generateHandlerFunc(handler))
}

func (a *Ion) PatchFunc(path string, handler http.HandlerFunc) *router.Route {
	return a.Router.Patch(path, a.generateHandlerFunc(handler))
}

func (a *Ion) PutFunc(path string, handler http.HandlerFunc) *router.Route {
	return a.Router.Put(path, a.generateHandlerFunc(handler))
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
// - GET  path/:id:	(get function)
// - PUT  path/:id:	(put function)
// - DELETE  path/:id	(delete function)
// The path MUST include the trailing slash.
func (a *Ion) RegisterREST(path string, handler RESTendpoint) {
	a.GetFunc(path+"/:id", handler.GET)
	a.PutFunc(path+"/:id", handler.PUT)
	a.DeleteFunc(path+"/:id", handler.DELETE)
	a.GetFunc(path, handler.LIST)
	a.PostFunc(path, handler.POST)
}

// This is the simpler possible handler: DoNothing do nothing.
// It can be used as a placeholder, or when we want to run
// the middleware for the collateral effects, but we don't want
// to do something specific.
// Also can be used with IonMVC.
func DoNothing(w http.ResponseWriter, r *http.Request) {
}

// Writes a JSON value
func MarshalJSON(w http.ResponseWriter, value interface{}) error {
	out, err := json.Marshal(value)
	if err != nil {
		return err
	}
	w.Write(out)
	return nil
}

// Returns the unmarshaled value from the request body
func UnmarshalJSON(r *http.Request, value interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, value)
}
