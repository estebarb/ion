// Package middleware contains general purpose middleware.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging provides a logging middleware
func Logging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

// DontPanic recovers from panics in other handlers
func DontPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// FormParser parses the forms in all the requests,
// so that you don't have to do it in the handlers/controllers.
func FormParser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
