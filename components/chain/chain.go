package chain

import "net/http"

type Chain struct {
	middlewares []Middleware
}

type Middleware func(next http.Handler) http.Handler

func New() *Chain {
	return &Chain{}
}

func Join(chains ...*Chain) *Chain {
	merged := New()
	for _, ch := range chains {
		merged.middlewares = append(merged.middlewares, ch.middlewares...)
	}
	return merged
}

func (m *Chain) Add(h Middleware) *Chain {
	m.middlewares = append(m.middlewares, h)
	return m
}

func (m *Chain) Then(h http.Handler) http.Handler {
	f := h
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		f = m.middlewares[i](f)
	}
	return f
}

func (m *Chain) ThenFunc(h http.HandlerFunc) http.Handler {
	return m.Then(http.HandlerFunc(h))
}

func Sequence(handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			for _, fun := range handlers {
				fun.ServeHTTP(w, r)
			}
		})
}
