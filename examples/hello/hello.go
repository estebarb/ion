package main

import (
	"fmt"
	"github.com/estebarb/ion"
	"github.com/estebarb/ion/components/router"
	"net/http"
)

func NewApp() *ion.Ion {
	ion.GetFunc("/", hello)
	ion.GetFunc("/:name", hello)
	return ion.App
}

func hello(w http.ResponseWriter, r *http.Request) {
	context := ion.App.Context(r)
	params := context.(router.IPathParam)
	value, exists := params.PathParams()["name"]
	if exists {
		fmt.Fprintf(w, "Hello, %v!", value)
	} else {
		fmt.Fprint(w, "Hello world!")
	}
}

func main() {
	http.ListenAndServe(":5500", NewApp())
}
