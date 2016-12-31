// Package chain provides functions for managing chains of middleware
// that can be used to provide reusable functionality that wraps
// http Handlers.
package chain

import "net/http"

// Chain contains a list of Middleware that will wrap a handler
type Chain struct {
	middleware []Middleware
}

// Middleware is a function that wraps an http.Handler
type Middleware func(next http.Handler) http.Handler

// New allocates a new Chain
func New() *Chain {
	return &Chain{}
}

// Join merges several chains into a single one
func Join(chains ...*Chain) *Chain {
	merged := New()
	for _, ch := range chains {
		merged.middleware = append(merged.middleware, ch.middleware...)
	}
	return merged
}

// Add appents a middleware to the end of a Chain
func (m *Chain) Add(h Middleware) *Chain {
	m.middleware = append(m.middleware, h)
	return m
}

// Then wraps a Handler with the middleware in the chain and returns the
// wrapped handler
func (m *Chain) Then(h http.Handler) http.Handler {
	f := h
	for i := len(m.middleware) - 1; i >= 0; i-- {
		f = m.middleware[i](f)
	}
	return f
}

// ThenFunc wraps a HandlerFunc with the middleware in the chain and returns the
// wrapped handler
func (m *Chain) ThenFunc(h http.HandlerFunc) http.Handler {
	return m.Then(http.HandlerFunc(h))
}

// Sequence returns a Handler that calls in order each handler passed
// as argument
func Sequence(handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			for _, fun := range handlers {
				fun.ServeHTTP(w, r)
			}
		})
}
