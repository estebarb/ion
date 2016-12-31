// reqctx provides an uniform way to access request context.
package reqctx

import (
	"net/http"
)

func BasicContextFactory() interface{} {
	return make(map[string]interface{})
}

func NewDefault() *State {
	return New(BasicContextFactory)
}

func (s *State) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := s.Context(r)
		r = s.WithContext(ctx, r)
		defer s.DestroyContext(r)
		next.ServeHTTP(w, r)
	})
}
