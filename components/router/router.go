// Package router contains a flexible router,
// with integrated context management per request
package router

import (
	"context"
	"net/http"
	"strings"
)

// Router register routes to be matched and
// dispatches its corresponding handler.
// As the router accepts path arguments then it fills
// the context with the them.
type Router struct {
	routeByName   map[string]*route
	routeByMethod map[string][]*route
}

// New creates a new router, with the given ContextFactory
func New() *Router {
	return &Router{
		routeByName:   make(map[string]*route),
		routeByMethod: make(map[string][]*route),
	}
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

// RouteFor returns a route corresponding to the requested
// route name.
// The arguments have the format:
// RouteFor(name, [key, value]*)
func (r *Router) RouteFor(name string, args ...string) string {
	route, ok := r.routeByName[name]
	if !ok || len(args)%2 != 0 {
		return ""
	}

	dst := make([]string, len(route.parsedPath))
	copy(dst, route.parsedPath)
	for i := 0; i < len(args); i += 2 {
		for k, v := range dst {
			if len(v) > 1 && v[0] == ':' && v[1:] == args[i] {
				dst[k] = args[i+1]
			}
		}
	}

	for _, v := range dst {
		if len(v) > 0 && v[0] == ':' {
			return ""
		}
	}
	return strings.Join(dst, "/")
}

// Handler register a handler to be dispatched when a request
// matches with the method and the path.
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

// HandleFunc register a handler to be dispatched when a request
// matches with the method and the path.
func (r *Router) HandleFunc(method string,
	path string,
	handler http.HandlerFunc) *Route {
	return r.Handler(method, path, handler)
}

// ServeHTTP dispatches the handler that matches with the request
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	routes := r.routeByMethod[method]
	path := splitWithoutTrailingSlash(req.URL.Path)

	for _, route := range routes {
		values, eq := equalPath(path, route.parsedPath)
		if eq {
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				ctx := req.Context()
				for k, v := range values {
					ctx = context.WithValue(ctx, k, v)
				}
				route.handler.ServeHTTP(w, req.WithContext(ctx))
			}).ServeHTTP(w, req)
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

// Route represents a Router matching rule, to be further refined.
type Route struct {
	router *Router
	route  *route
}

// Name assigns an identifier to the Route. This allows to use RouteFor
// to construct a path that could match this rule.
func (r *Route) Name(name string) *Route {
	if name != "" {
		r.router.routeByName[name] = r.route
	}
	return r
}

// Get register the handler in the router, after wrapping it with the middleware
func (r *Router) Get(path string, handler http.Handler) *Route {
	return r.Handler(http.MethodGet, path, handler)
}

// Post register the handler in the router, after wrapping it with the middleware
func (r *Router) Post(path string, handler http.Handler) *Route {
	return r.Handler(http.MethodPost, path, handler)
}

// Put register the handler in the router, after wrapping it with the middleware
func (r *Router) Put(path string, handler http.Handler) *Route {
	return r.Handler(http.MethodPut, path, handler)
}

// Delete register the handler in the router, after wrapping it with the middleware
func (r *Router) Delete(path string, handler http.Handler) *Route {
	return r.Handler(http.MethodDelete, path, handler)
}

// Patch register the handler in the router, after wrapping it with the middleware
func (r *Router) Patch(path string, handler http.Handler) *Route {
	return r.Handler(http.MethodPatch, path, handler)
}

// Options register the handler in the router, after wrapping it with the middleware
func (r *Router) Options(path string, handler http.Handler) *Route {
	return r.Handler(http.MethodOptions, path, handler)
}

// GetFunc register the handler in the router, after wrapping it with the middleware
func (r *Router) GetFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc(http.MethodGet, path, handler)
}

// PostFunc register the handler in the router, after wrapping it with the middleware
func (r *Router) PostFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc(http.MethodPost, path, handler)
}

// PutFunc register the handler in the router, after wrapping it with the middleware
func (r *Router) PutFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc(http.MethodPut, path, handler)
}

// DeleteFunc register the handler in the router, after wrapping it with the middleware
func (r *Router) DeleteFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc(http.MethodDelete, path, handler)
}

// PatchFunc register the handler in the router, after wrapping it with the middleware
func (r *Router) PatchFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc(http.MethodPatch, path, handler)
}

// OptionsFunc register the handler in the router, after wrapping it with the middleware
func (r *Router) OptionsFunc(path string, handler http.HandlerFunc) *Route {
	return r.HandleFunc(http.MethodOptions, path, handler)
}
