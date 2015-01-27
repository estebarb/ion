package main

import (
	"github.com/estebarb/ion"
	"html/template"
	"log"
	"net/http"
)

func main() {
	r := ion.NewRouter()
	t := template.Must(template.ParseFiles("template.html"))
	r.GetFunc("/", ion.RenderTemplate(t))
	r.GetFunc("/:name", ion.RenderTemplate(t))
	log.Fatal(http.ListenAndServe(":8080", r))
}
