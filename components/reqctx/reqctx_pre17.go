// +build !go1.7

// reqctx provides an uniform way to access request context.
package reqctx

import (
	"net/http"
	"sync"
)

type State struct {
	sync.RWMutex
	data       map[*http.Request]interface{}
	ctxFactory func() interface{}
}

func New(ctxFactory func() interface{}) *State {
	return &State{
		data:       make(map[*http.Request](interface{})),
		ctxFactory: ctxFactory,
	}
}

func (s *State) Context(r *http.Request) interface{} {
	s.RLock()
	defer s.RUnlock()
	state, ok := s.data[r]
	if !ok {
		return s.ctxFactory()
	}
	return state
}

func (s *State) WithContext(ctx interface{}, r *http.Request) *http.Request {
	s.Lock()
	defer s.Unlock()
	s.data[r] = ctx
	return r
}

func (s *State) DestroyContext(r *http.Request) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, r)
}

func (s *State) size() int {
	return len(s.data)
}
