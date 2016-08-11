# Ion Web Framework
[![GoDoc](https://godoc.org/github.com/estebarb/ion?status.svg)](http://godoc.org/github.com/estebarb/ion)
[![Build Status](https://travis-ci.org/estebarb/ion.svg?branch=master)](https://travis-ci.org/estebarb/ion)


Ion is a small web framework for people in a hurry.
Ion provides a flexible router (gorilla/mux),
a integrated middleware support (justinas/alice), automatic
context support and some handful helpers.

A short example:

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
	
At this point the framework is highly experimental, so please don't
use it in production for now... I'm planning to add more features,
but maybe I will break things. Don't say I don't tell you! :p

## License

Ion is released under the MIT License, as specified in LICENSE.
