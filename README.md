Ion Web Framework [![GoDoc](https://godoc.org/github.com/estebarb/ion?status.svg)](http://godoc.org/github.com/estebarb/ion)
=================

Ion is a small web framework for people in a hurry.
Ion provides a fast trie based router (julienschmidt/httprouter),
a integrated middleware support (justinas/alice), automatic
context support (via gorilla/context) and some handful helpers.

A short example:

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
	
At this point the framework is highly experimental, so please don't
use it in production for now... I'm planning to add more features,
but maybe I will break things. Don't say I don't tell you! :p
