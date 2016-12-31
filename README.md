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
    	"net/http"
    )
    
    type App struct {
    	*ion.Ion
    }
    
    func newApp() *App {
    	app := &App{
    		Ion: ion.New(),
    	}
    	app.GetFunc("/", app.hello)
    	app.GetFunc("/:name", app.hello)
    	return app
    }
    
    func (app *App) hello(w http.ResponseWriter, r *http.Request) {
    	state := app.Router.GetState(r)
    	value, exists := state.Get("name")
    	if exists {
    		fmt.Fprintf(w, "Hello, %v!", value)
    	} else {
    		fmt.Fprint(w, "Hello world!")
    	}
    }
    
    func main() {
    	http.ListenAndServe(":5500", newApp())
    }
	
At this point the framework is highly experimental, so please don't
use it in production for now... I'm planning to add more features,
but maybe I will break things. Don't say I didn't tell you! :p

## License

Ion is released under the MIT License, as specified in LICENSE.
