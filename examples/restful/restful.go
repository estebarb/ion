package main

import (
	"fmt"
	"github.com/estebarb/ion"
	"github.com/estebarb/ion/components/router"
	"net/http"
)

type restHandler struct {
	*ion.Ion
}

func (c *restHandler) URLArgs(r *http.Request, key string) string {
	state := c.Context(r).(router.IPathParam)
	value, ok := state.PathParams()[key]
	if !ok {
		return ""
	}
	return value
}

func (c restHandler) LIST(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "list things")
}

func (c restHandler) POST(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "post new thing")
}

func (c restHandler) PUT(w http.ResponseWriter, r *http.Request) {
	val := c.URLArgs(r, "id")
	fmt.Fprintf(w, "upsert thing as %v", val)
}

func (c restHandler) GET(w http.ResponseWriter, r *http.Request) {
	val := c.URLArgs(r, "id")
	fmt.Fprintf(w, "get thing %v", val)
}

func (c restHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	val := c.URLArgs(r, "id")
	fmt.Fprintf(w, "delete thing %v", val)
}

func main() {
	var rest restHandler
	ion.GetFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `Please try the following endpoints:
	- GET       /test/      Lists the items
	- POST      /test/      Post a new item
	- GET       /test/{id}   Gets a item
	- DELETE    /test/{id}   Deletes a item
	- PUT       /test/{id}   Updates a item`)
	})
	ion.App.RegisterREST("/test/", rest)
	http.ListenAndServe(":8080", ion.App)
}
