// +build go1.7

// reqctx provides an uniform way to access request context.
package reqctx

import (
	"context"
	"net/http"
)

type State struct {
	ctxFactory func() interface{}
}

func New(ctxFactory func() interface{}) *State {
	return &State{
		ctxFactory: ctxFactory,
	}
}

func (s *State) Context(r *http.Request) interface{} {
	state := r.Context().Value(s)
	if state == nil {
		return s.ctxFactory()
	}
	return state
}

func (s *State) WithContext(ctx interface{}, r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), s, ctx))
}

func (s *State) DestroyContext(r *http.Request) {
}

func (s *State) size() int {
	return 0
}
