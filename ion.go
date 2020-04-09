// Package ion implements a small web framework that allows to
// easily connect reusable components.
//
// This framework is based on https://blog.gopheracademy.com/advent-2016/go-syntax-for-dsls/
// idea of using a DSL. This approach naturally removes the need to implement
// a router, allowing the framework to just reuse Go standard mux.
//
// Ion have the following features:
//
// - Can use easily any http.Handler or http.HandlerFunc
// - Easily describe paths (with arguments) and method handlers
// - Compatible with Middlewares
// - Use context for passing path arguments
//
package ion

import (
	"context"
	"log"
	"net/http"
	"strings"
)

// Middleware is a function that wrap an http.Handler and returns a value
// that implements the http.Handler interface
type Middleware func(next http.Handler) http.Handler

// Chain describes a secuence of Middleware
type Chain []Middleware

// Then composes all the middlewares wrapping the given
// http.Handler, and returns a new http.Handler
func (c Chain) Then(h http.Handler) http.Handler {
	f := h
	for i := len(c) - 1; i >= 0; i-- {
		f = c[i](f)
	}
	return f
}

// Endpoint describes a http request handler, that may
// have optional Middleware
type Endpoint struct {
	Middleware  []Middleware
	Handler     Builder
	HttpHandler http.Handler
}

// Build generates an http.Handler from an Endpoint
func (e Endpoint) Build() http.Handler {
	if (e.HttpHandler != nil) == (e.Handler != nil) {
		panic("Endpoint support only Handler or HttpHandler, not both")
	}

	if e.HttpHandler != nil {
		return Chain(e.Middleware).Then(e.HttpHandler)
	}
	return Chain(e.Middleware).Then(e.Handler.Build())
}

// Builder interface is implemented by objects that can be build
// into an http.Handler
type Builder interface {
	Build() http.Handler
}

// Routes describe a request router that handles request according
// to its path
type Routes map[string]Endpoint

// Build returns an http.Handler that can handle requests by path
func (r Routes) Build() http.Handler {
	mux := http.NewServeMux()
	for prefix, endpoint := range r {
		if strings.HasPrefix(prefix, "/:") {
			parts := strings.Split(strings.TrimPrefix(prefix, "/:"), "/")
			name := parts[0]
			mux.Handle("/", captureArgument(name)(endpoint.Build()))
		} else {
			mux.Handle(prefix, http.StripPrefix(prefix, endpoint.Build()))
		}
	}
	return mux
}

// PathEnd is a middleware used to "cut" the requests path at the current level.
// The request is handled if the path is "/" or "".
func PathEnd(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		if r.URL.Path != "/" && r.URL.Path != "" {
			http.NotFound(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func captureArgument(name string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
			var value string
			if len(parts) > 0 {
				value = parts[0]
			}
			ctx := context.WithValue(r.Context(), name, value)
			log.Println(name, value)
			r2 := r.WithContext(ctx)

			http.StripPrefix("/"+value, next).ServeHTTP(w, r2)
		})
	}
}

// Methods implement an http.Handler that handles requests according to
// the request method.
type Methods map[string]Endpoint

// Build generates an http.Handler
func (m Methods) Build() http.Handler {
	handlers := make(map[string]http.Handler)
	for k, v := range m {
		handlers[k] = v.Build()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := handlers[r.Method]
		if handler != nil {
			handler.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}
