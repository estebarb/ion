// Middleware contains general purpose middleware.
package middleware

import (
	"github.com/estebarb/ion/context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

/*
Ion adds the arguments to the request context, using
this key. The application can retrieve the arguments
using:

	name := ion.URLArgs(r, "name")
*/
const urlargs = "urlargs"

// Inserts the path variables in the context
func ContextMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, urlargs, mux.Vars(r))
		defer context.Clear(r)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Provides a logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

// Provides a recovery from panics in other handlers
func PanicMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC]: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// FormParser parses the forms in all the requests,
// so that you don't have to do it in the handlers/controllers.
func FormParserMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
