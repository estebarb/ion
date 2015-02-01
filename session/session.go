/*
Package sessions provides the Gorilla Framework Sessions as a
middleware for Ion.

	package main

	import (
		"github.com/estebarb/ion"
		"github.com/estebarb/ion/session"
		"net/http"
		"github.com/gorilla/sessions"
		"fmt"
		"math/rand"
	)

	var store = sessions.NewCookieStore([]byte(""))

	func handler(w http.ResponseWriter, req *http.Request) {
		s := session.GetSession(req, "ion_session")
		oldValue := s.Values["randomNumber"]
		newValue := rand.Intn(100)
		s.Values["randomNumber"] = newValue
		fmt.Fprintf(w, `Refresh the page and the new value will become the old one.
				Old Value: %v
				New Value: %v
		Values: %v`, oldValue, newValue, s.Values)
	}

	func main() {
		store.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3600 * 8, // 8 hours
			HttpOnly: true,
		}
		r := ion.NewRouterDefaults(session.Sessions("ion_session", store))
		r.GetFunc("/", handler)
		http.ListenAndServe(":8010", r)
	}

*/
package session

import (
	"github.com/estebarb/ion/context"
	gs "github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"net/http"
	"net/http/httptest"
)

var sessionData = context.New()

// Sessions is a Middleware that maps a session.Session service into the Ion middleware chain.
// Sessions can use a number of storage solutions with the given store.
func Sessions(name string, store gs.Store) alice.Constructor {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			s, _ := store.Get(r, name)
			defer sessionData.Clear(r)
			sessionData.Set(r, name, s)
			wrec := httptest.NewRecorder()

			defer func() {
				// we copy the original headers first
				for k, v := range wrec.Header() {
					w.Header()[k] = v
				}
				store.Save(r, w, s)
				w.WriteHeader(wrec.Code)
				w.Write(wrec.Body.Bytes())
			}()

			next.ServeHTTP(wrec, r)
		}
		return http.HandlerFunc(fn)
	}
}

func GetSession(r *http.Request, name string) *gs.Session {
	s, _ := sessionData.Get(r, name)
	return s.(*gs.Session)
}
