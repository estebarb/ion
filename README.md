# Ion Web Framework
[![GoDoc](https://godoc.org/github.com/estebarb/ion?status.svg)](http://godoc.org/github.com/estebarb/ion)    [![Build Status](https://travis-ci.org/estebarb/ion.svg?branch=master)](https://travis-ci.org/estebarb/ion)    [![codecov](https://codecov.io/gh/estebarb/ion/branch/master/graph/badge.svg)](https://codecov.io/gh/estebarb/ion)    [![Go Report Card](https://goreportcard.com/badge/github.com/estebarb/ion)](https://goreportcard.com/report/github.com/estebarb/ion)


Ion is a small web framework written in Go.

Ion leverages component composition to allow functionality reuse, and
build complex behaviors based on simple components.

A short example:

    package main
    
    import (
        "fmt"
        "github.com/estebarb/ion"
        "github.com/estebarb/ion/components/router"
        "net/http"
    )
    
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
        ion.GetFunc("/", hello)
        ion.GetFunc("/:name", hello)
        http.ListenAndServe(":5500", ion.App)
    }
    
At this point the framework is highly experimental, so please take that
in consideration if you want to use it.

## License

Ion is released under the MIT License, as specified in LICENSE.
