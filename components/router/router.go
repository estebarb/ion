// Package router contains a flexible router,
// with advanced middleware handling
package router

import (
	"github.com/estebarb/ion/components/reqctx"
	"html/template"
	"net/http"
	"strings"
)

type Router struct {
	*reqctx.State
	routeByName   map[string]*route
	routeByMethod map[string]([]*route)
}

type IPathParam interface {
	PathParams() map[string]string
	SetPathParams(values map[string]string)
}

type Path struct {
	Params map[string]string
}

func (c *Path) PathParams() map[string]string {
	return c.Params
}

func (c *Path) SetPathParams(values map[string]string) {
	c.Params = values
}

func New(ContextFactory func() interface{}) *Router {
	return &Router{
		routeByName:   make(map[string]*route),
		routeByMethod: make(map[string]([]*route)),
		State:         reqctx.New(ContextFactory),
	}
}

func NewDefault() *Router {
	return New(ContextFactory)
}

func ContextFactory() interface{} {
	return &Path{}
}

type route struct {
	handler    http.Handler
	path       string
	parsedPath []string
	name       string
	method     string
}

func splitWithoutTrailingSlash(str string) []string {
	parsedPath := strings.Split(str, "/")
	if parsedPath[len(parsedPath)-1] == "" {
		parsedPath = parsedPath[:len(parsedPath)-1]
	}
	return parsedPath
}

// BuildRoute returns a route corresponding to de requested
// route name.
// The arguments have the format:
// BuildRoute(name, [key, value]*)
func (r *Router) BuildRoute(name string, args ...string) template.URL {
	route, ok := r.routeByName[name]
	if !ok || len(args)%2 != 0 {
		return ""
	}

	dst := make([]string, len(route.parsedPath))
	copy(dst, route.parsedPath)
	for i := 0; i < len(args); i += 2 {
		for k, v := range dst {
			if len(v) > 1 && v[0] == ':' && string(v[1:]) == args[i] {
				dst[k] = args[i+1]
			}
		}
	}

	for _, v := range dst {
		if len(v) > 0 && v[0] == ':' {
			return template.URL("")
		}
	}
	return template.URL(strings.Join(dst, "/"))
}

func (r *Router) Handler(method string,
	path string,
	handler http.Handler) *Route {
	routes, ok := r.routeByMethod[method]
	if !ok {
		routes = make([]*route, 0)
	}

	newRoute := &route{
		handler:    handler,
		path:       path,
		parsedPath: splitWithoutTrailingSlash(path),
		method:     method,
	}
	routes = append(routes, newRoute)
	r.routeByMethod[method] = routes
	return &Route{
		route:  newRoute,
		router: r,
	}
}

func (r *Router) HandleFunc(method string,
	path string,
	handler http.HandlerFunc) *Route {
	return r.Handler(method, path, http.HandlerFunc(handler))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	routes := r.routeByMethod[method]
	path := splitWithoutTrailingSlash(req.URL.Path)

	for _, route := range routes {
		values, eq := equalPath(path, route.parsedPath)
		if eq {
			r.Middleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				context := r.Context(req)
				context.(IPathParam).SetPathParams(values)
				route.handler.ServeHTTP(w, req)
			})).ServeHTTP(w, req)
			return
		}
	}
	http.NotFound(w, req)
}

func equalPath(path, pattern []string) (map[string]string, bool) {
	values := make(map[string]string)
	if len(path) != len(pattern) {
		return nil, false
	}
	for k, v := range path {
		pat := pattern[k]
		if len(pat) > 0 && pat[0] == ':' && len(v) > 0 {
			// We got a variable
			values[pat[1:]] = v
		} else if pat != v {
			// We found an invalid path
			return nil, false
		}
	}
	return values, true
}

type Route struct {
	router *Router
	route  *route
}

func (r *Route) Name(name string) *Route {
	if name != "" {
		r.router.routeByName[name] = r.route
	}
	return r
}

func (r *Router) Get(path string, handler http.Handler) *Route {
	return r.Handler("GET", path, handler)
}

func (r *Router) Post(path string, handler http.Handler) *Route {
	return r.Handler("POST", path, handler)
}

func (r *Router) Put(path string, handler http.Handler) *Route {
	return r.Handler("PUT", path, handler)
}

func (r *Router) Delete(path string, handler http.Handler) *Route {
	return r.Handler("DELETE", path, handler)
}

func (r *Router) Patch(path string, handler http.Handler) *Route {
	return r.Handler("PATCH", path, handler)
}

func (r *Router) Options(path string, handler http.Handler) *Route {
	return r.Handler("OPTIONS", path, handler)
}

func (r *Router) GetFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc("GET", path, handler)
}

func (r *Router) PostFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc("POST", path, handler)
}

func (r *Router) PutFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc("PUT", path, handler)
}

func (r *Router) DeleteFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc("DELETE", path, handler)
}

func (r *Router) PatchFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc("PATCH", path, handler)
}

func (r *Router) OptionsFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc("OPTIONS", path, handler)
}
