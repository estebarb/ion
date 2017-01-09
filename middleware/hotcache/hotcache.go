// Package hotcache intercepts and group equal requests,
// perform a single server request.
// Is ideal for caching really hot pages, like front pages.
// Is not ideal for caching responses that depends on the logged user.
package hotcache

import (
	"github.com/estebarb/ion/futures"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

// Config contains the configuration for hot caching of requests
type Config struct {
	l       sync.Mutex
	timeout time.Duration
	content map[string]*futures.Expirable
}

// New creates a new configurated Config for hot caching
func New(timeout time.Duration) *Config {
	return &Config{
		timeout: timeout,
		content: make(map[string]*futures.Expirable),
	}
}

func execute(r *http.Request, next http.Handler) *httptest.ResponseRecorder {
	wrec := httptest.NewRecorder()
	next.ServeHTTP(wrec, r)
	return wrec
}

// Middleware wraps a request and hot caches it
func (c *Config) Middleware(next http.Handler) http.Handler {
	fun := func(w http.ResponseWriter, r *http.Request) {
		c.l.Lock()
		expirable, ok := c.content[r.URL.Path]
		if !ok {
			expirable = futures.NewExpirable(c.timeout,
				func() interface{} {
					return execute(r, next)
				})
			c.content[r.URL.Path] = expirable
		}
		c.l.Unlock()
		recorded := expirable.Read().(*httptest.ResponseRecorder)
		for k, v := range recorded.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(recorded.Code)
		w.Write(recorded.Body.Bytes())
	}
	return http.HandlerFunc(fun)
}
