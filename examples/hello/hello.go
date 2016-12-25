package main

import (
	"fmt"
	"github.com/estebarb/ion"
	"net/http"
)

func NewApp() *ion.Ion {
	ion.GetFunc("/", hello)
	ion.GetFunc("/:name", hello)
	return ion.App
}

func hello(w http.ResponseWriter, r *http.Request) {
	state := ion.App.Router.GetState(r)
	value, exists := state.Get("name")
	if exists {
		fmt.Fprintf(w, "Hello, %v!", value)
	} else {
		fmt.Fprint(w, "Hello world!")
	}
}

func main() {
	http.ListenAndServe(":5500", NewApp())
}
