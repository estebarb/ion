package main

import (
	"github.com/estebarb/ion"
	"github.com/estebarb/ion/context"
	"html/template"
	"log"
	"net/http"
)

var t = template.Must(template.ParseFiles("context.tmpl"))

func handler(w http.ResponseWriter, r *http.Request) {
	context.Set(r, "message", "Hello world from context value!")
	ion.RenderTemplate(t).ServeHTTP(w, r)
}

func main() {
	r := ion.NewRouter()
	r.GetFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
