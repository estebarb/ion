package main

import (
	"fmt"
	"github.com/estebarb/ion"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	value := ion.URLArgs(r, "name")
	if value != "" {
		fmt.Fprintf(w, "Hello, %v!", value)
	} else {
		fmt.Fprint(w, "Hello world!")
	}
}

func main() {
	r := ion.NewRouter()
	r.GetFunc("/", hello)
	r.GetFunc("/{name}", hello)
	http.ListenAndServe(":8080", r)
}
