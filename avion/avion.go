// Avion is a experimental web framework, trying to provide a higher level of
// abstraction to the programmer.
package avion

import (
	"fmt"
	"github.com/estebarb/ion"
	"github.com/estebarb/ion/context"
	"github.com/estebarb/ion/iongae"
	"github.com/estebarb/ion/middleware"
	"net/http"
	"path/filepath"
	"strings"
)

type Avion struct {
	Router *ion.Router
}

// Initializes a ion.Router, using a set of default middleware
func NewAvion() *Avion {
	return &Avion{Router: ion.NewRouterDefaults(
		middleware.PanicMiddleware,
		middleware.LoggingMiddleware,
		middleware.FormParserMiddleware,
		extractFormat,
	),
	}
}

func NewGAEAvion() *Avion {
	return &Avion{Router: ion.NewRouterDefaults(
		iongae.GAEPanicMiddleware,
		middleware.LoggingMiddleware,
		middleware.FormParserMiddleware,
		extractFormat,
	),
	}
}

func extractFormat(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "FORMAT",
			strings.ToLower(filepath.Ext(r.URL.Path)))
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func templateName(c Controller, method string, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "TEMPLATE",
			c.Name()+"_"+method)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Format(r *http.Request) (string, bool) {
	str, ok := context.Get(r, "FORMAT")
	return str.(string), ok
}

type Controller interface {
	Name() string
	HandlerIndex(w http.ResponseWriter, r *http.Request)
	HandlerNew(w http.ResponseWriter, r *http.Request)
	HandlerCreate(w http.ResponseWriter, r *http.Request)
	HandlerShow(w http.ResponseWriter, r *http.Request)
	HandlerEdit(w http.ResponseWriter, r *http.Request)
	HandlerUpdate(w http.ResponseWriter, r *http.Request)
	HandlerDestroy(w http.ResponseWriter, r *http.Request)
}

type Pluggable interface {
	PluginIndex(next http.Handler) http.Handler
	PluginNew(next http.Handler) http.Handler
	PluginCreate(next http.Handler) http.Handler
	PluginShow(next http.Handler) http.Handler
	PluginEdit(next http.Handler) http.Handler
	PluginUpdate(next http.Handler) http.Handler
	PluginDestroy(next http.Handler) http.Handler
}

type Model interface {
	Create(data interface{}) (string, error)
	Read(id string) (interface{}, error)
	Update(id string, data interface{}) error
	Delete(id string) error
	Query(name string, data interface{}) ([]interface{}, error)
}

type BasicController int

func (c *BasicController) HandlerIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Avion Default Index Handler")
}

func (c *BasicController) HandlerNew(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Avion Default New Handler")
}

func (c *BasicController) HandlerCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Avion Default Create Handler")
}

func (c *BasicController) HandlerShow(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Avion Default Show Handler")
}

func (c *BasicController) HandlerEdit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Avion Default Edit Handler")
}

func (c *BasicController) HandlerUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Avion Default Update Handler")
}

func (c *BasicController) HandlerDestroy(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Avion Default Destroy Handler")
}

type BasicPlugin int

func (p *BasicPlugin) PluginIndex(next http.Handler) http.Handler {
	return next
}
func (p *BasicPlugin) PluginNew(next http.Handler) http.Handler {
	return next
}
func (p *BasicPlugin) PluginCreate(next http.Handler) http.Handler {
	return next
}
func (p *BasicPlugin) PluginShow(next http.Handler) http.Handler {
	return next
}
func (p *BasicPlugin) PluginEdit(next http.Handler) http.Handler {
	return next
}
func (p *BasicPlugin) PluginUpdate(next http.Handler) http.Handler {
	return next
}
func (p *BasicPlugin) PluginDestroy(next http.Handler) http.Handler {
	return next
}

func generatePath(cs ...Controller) string {
	s := ""
	if len(cs) > 0 {
		for _, v := range cs[0 : len(cs)-1] {
			s = s + "/" + v.Name() + "/+{" + v.Name() + "_id}"
		}
	}
	return s
}

// Add a controller to the router, using the others controllers as ancestors
func (a *Avion) RouteController(c Controller, ancestors ...Controller) {
	base := generatePath(ancestors...)
	entity := c.Name()
	path := base + "/" + entity
	a.Router.GetFunc(path, c.HandlerIndex)
	a.Router.GetFunc(path+"/new", c.HandlerNew)
	a.Router.PostFunc(path, c.HandlerCreate)
	a.Router.GetFunc(path+"/{id}", c.HandlerShow)
	a.Router.GetFunc(path+"/{id}/edit", c.HandlerEdit)
	a.Router.PutFunc(path+"/{id}", c.HandlerUpdate)
	a.Router.PatchFunc(path+"/{id}", c.HandlerUpdate)
	a.Router.DeleteFunc(path+"/{id}", c.HandlerDestroy)
	a.Router.PostFunc(path+"/{id}/delete", c.HandlerDestroy)
}

type PluggableController interface {
	Pluggable
	Controller
}

func (a *Avion) RouteControllerWithPlugin(c PluggableController, ancestors ...Controller) {
	base := generatePath(ancestors...)
	entity := c.Name()
	path := base + "/" + entity
	a.Router.GET(path, templateName(c, "index", c.PluginIndex(http.HandlerFunc(c.HandlerIndex))))
	a.Router.GET(path+"/new", templateName(c, "new", c.PluginNew(http.HandlerFunc(c.HandlerNew))))
	a.Router.POST(path, templateName(c, "create", c.PluginCreate(http.HandlerFunc(c.HandlerCreate))))
	a.Router.GET(path+"/{id}", templateName(c, "show", c.PluginShow(http.HandlerFunc(c.HandlerShow))))
	a.Router.GET(path+"/{id}/edit", templateName(c, "edit", c.PluginEdit(http.HandlerFunc(c.HandlerEdit))))
	a.Router.PUT(path+"/{id}", templateName(c, "update", c.PluginUpdate(http.HandlerFunc(c.HandlerUpdate))))
	a.Router.PATCH(path+"/{id}", templateName(c, "update", c.PluginUpdate(http.HandlerFunc(c.HandlerUpdate))))
	a.Router.DELETE(path+"/{id}", templateName(c, "destroy", c.PluginDestroy(http.HandlerFunc(c.HandlerDestroy))))
	a.Router.POST(path+"/{id}/delete", templateName(c, "destroy", c.PluginDestroy(http.HandlerFunc(c.HandlerDestroy))))
}
