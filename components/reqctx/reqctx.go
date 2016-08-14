// reqctx provides an uniform way to access request context.
package reqctx

import (
	"net/http"
	"sync"
)

type State struct {
	data map[string]interface{}
}

func NewState() *State {
	return &State{
		data: make(map[string]interface{}),
	}
}

func (s *State) GetAll() map[string]interface{} {
	return s.data
}

func (s *State) Set(key string, value interface{}) {
	s.data[key] = value
}

func (s *State) Get(key string) (interface{}, bool) {
	v, ok := s.data[key]
	return v, ok
}

func (s *State) Delete(key string) {
	delete(s.data, key)
}

type StateContainer struct {
	lock sync.RWMutex
	data map[*http.Request](*State)
}

func NewStateContainer() *StateContainer {
	return &StateContainer{
		lock: sync.RWMutex{},
		data: make(map[*http.Request](*State)),
	}
}

func (s *StateContainer) GetState(r *http.Request) *State {
	s.lock.Lock()
	defer s.lock.Unlock()
	state, ok := s.data[r]
	if !ok {
		state = NewState()
		s.data[r] = state
	}
	return state
}

func (s *StateContainer) DestroyState(r *http.Request) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, r)
}

func (s *StateContainer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.GetState(r)
		defer s.DestroyState(r)
		next.ServeHTTP(w, r)
	})
}
