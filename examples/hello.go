package main

import (
	"fmt"
	"github.com/estebarb/ion"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	val := context.Get(r, ion.Urlargs).(httprouter.Params)
	if val != nil {
		fmt.Fprintf(w, "Hello, %v!", val.ByName("name"))
	} else {
		fmt.Fprint(w, "Hello world!")
	}
}

func main() {
	r := ion.NewRouter()
	r.GetFunc("/", hello)
	r.GetFunc("/:name", hello)
	http.ListenAndServe(":8080", r)
}
