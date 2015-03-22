package main

import (
	"fmt"
	"github.com/estebarb/ion"
	"github.com/estebarb/ion/session"
	"github.com/gorilla/sessions"
	"math/rand"
	"net/http"
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
