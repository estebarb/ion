package main

import (
	"fmt"
	"github.com/estebarb/ion"
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
	val := ion.URLArgs(r, "id")
	fmt.Fprint(w, "upsert thing as %v", val)
}

func (c restHandler) GET(w http.ResponseWriter, r *http.Request) {
	val := ion.URLArgs(r, "id")
	fmt.Fprintf(w, "get thing %v", val)
}

func (c restHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	val := ion.URLArgs(r, "id")
	fmt.Fprintf(w, "delete thing %v", val)
}

func main() {
	r := ion.NewRouter()
	var rest restHandler
	r.GetFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, `Please try the following endpoints:
	- GET       /test/      Lists the items
	- POST      /test/      Post a new item
	- GET       /test/:id   Gets a item
	- DELETE    /test/:id   Deletes a item
	- PUT       /test/:id   Updates a item`)
	})
	r.RegisterREST("/test/", rest)
	http.ListenAndServe(":8080", r)
}
