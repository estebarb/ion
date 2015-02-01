package main

import (
	"fmt"
	"github.com/estebarb/ion"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type restHandler int

func (c restHandler) LIST(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "list things")
}

func (c restHandler) POST(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "post new thing")
}

func (c restHandler) PUT(w http.ResponseWriter, r *http.Request) {
	val := context.Get(r, ion.Urlargs).(httprouter.Params)
	fmt.Fprint(w, "upsert thing as %v", val.ByName("id"))
}

func (c restHandler) GET(w http.ResponseWriter, r *http.Request) {
	val := context.Get(r, ion.Urlargs).(httprouter.Params)
	fmt.Fprintf(w, "get thing %v", val.ByName("id"))
}

func (c restHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	val := context.Get(r, ion.Urlargs).(httprouter.Params)
	fmt.Fprintf(w, "delete thing %v", val.ByName("id"))
}

func main() {
	r := ion.NewRouter()
	var rest restHandler
	r.RegisterREST("/test/", rest)
	http.ListenAndServe(":8080", r)
}
