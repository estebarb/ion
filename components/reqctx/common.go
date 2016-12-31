// Package reqctx provides an uniform way to access request context.
//
// For this package a context is a single value that is "carried" by
// the request. This single value is an interface{}, and is expected
// that the programs provide its own ContextFactory, and have a
// corresponding interface that could be used to make safe type assertion.
//
// This can be both in Go 1.7+ or previous versions. If used in Go 1.7+ it
// stores the context in the Request context, but if used in Go <1.7 then
// it stores it in a map[*Request] protected by mutexes. The user code
// don't need to any code change to work in both versions, as long as
// it uses reqctx in the recommended way.
//
// If used with the Ion Router then you don't have to worry about allocating
// the context or destroying it after usage.
package reqctx

import (
	"net/http"
)

// BasicContextFactory creates a map[string]interface{} where the request
// context could be stored.
//
// Note that this one is intended for trivial usage or reqctx testing, and
// must not be used with the rest of Ion.
//
// In Ion is expected that the application provides its own ContextFactory,
// or at least uses the default ContextFactory of the router package.
func BasicContextFactory() interface{} {
	return make(map[string]interface{})
}

// NewDefault returns a new State manager that uses BasicContextFactory
func NewDefault() *State {
	return New(BasicContextFactory)
}

// Middleware wraps a handler with context creation and destruction.
//
// Note that the user shouldn't use this one when using Ion, as the router
// inserts it by default (the router puts the path arguments in the context).
func (s *State) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := s.Context(r)
		r = s.WithContext(ctx, r)
		defer s.DestroyContext(r)
		next.ServeHTTP(w, r)
	})
}
