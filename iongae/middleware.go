// +build appengine
package iongae

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
)


func GAEPanicMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				c := appengine.NewContext(r)
				log.Errorf(c, "[PANIC] %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
