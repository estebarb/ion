// +build appengine

// Middleware specific for GAE
package iongae

import (
	"appengine"
	"net/http"
)


func GAEPanicMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				c := appengine.NewContext(r)
				c.Errorf("[PANIC] %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}