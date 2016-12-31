// +build !go1.7

package reqctx

import (
	"net/http"
	"sync"
)

// State is used to manage the requests context
type State struct {
	sync.RWMutex
	data       map[*http.Request]interface{}
	ctxFactory func() interface{}
}

// New creates a new context manager
func New(ctxFactory func() interface{}) *State {
	return &State{
		data:       make(map[*http.Request](interface{})),
		ctxFactory: ctxFactory,
	}
}

// Context returns the context of the request
func (s *State) Context(r *http.Request) interface{} {
	s.RLock()
	defer s.RUnlock()
	state, ok := s.data[r]
	if !ok {
		return s.ctxFactory()
	}
	return state
}

// WithContext associates the context with the request. It must be done once.
func (s *State) WithContext(ctx interface{}, r *http.Request) *http.Request {
	s.Lock()
	defer s.Unlock()
	s.data[r] = ctx
	return r
}

// DestroyContext destroys the context. Note that in Go 1.7+ this is a No-OP,
// as the context is stored with the Request own context. In <1.7 this removes
// the context from the internal map
func (s *State) DestroyContext(r *http.Request) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, r)
}

// size returns the number contexts stored in the internal map request. In 1.7+
// is always zero, but is left here to make the tests compatible between
// versions.
func (s *State) size() int {
	return len(s.data)
}
