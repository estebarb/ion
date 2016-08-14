package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_Get(t *testing.T) {
	r := New()
	r.Get("/hello",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello"))
		}),
		"hello")
	ts := httptest.NewServer(r)

	res, err := http.Get(ts.URL + "/hello")
	if err != nil {
		t.Error("Unexpected errror:", err)
	}
	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error("Unexpected errror:", err)
	}
	if string(response) != "Hello" {
		t.Error("Expecting <Hello>, received:", string(response))
	}
}

func TestRouter_GetWithArguments(t *testing.T) {
	router := New()
	router.Get("/hello/:name/:number/world",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			state := router.GetState(r)
			args, ok := state.Get("path")
			if !ok {
				http.Error(w, "Can't read path arguments", 500)
			}
			arguments := args.(map[string]string)
			fmt.Fprintf(w, "/hello/%s/world/%s",
				arguments["name"],
				arguments["number"])
		}),
		"hello")
	ts := httptest.NewServer(router)

	res, err := http.Get(ts.URL + "/hello/test/001/world")
	if err != nil {
		t.Error("Unexpected errror:", err)
	}
	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error("Unexpected errror:", err)
	}
	if string(response) != "/hello/test/world/001" {
		t.Error("Expecting /hello/test/world/001, received:", string(response))
	}
}

func TestBuildRoute(t *testing.T) {
	router := New()
	router.Get("/hello/:name/:number/world",
		http.NotFoundHandler(),
		"hello")
	url := router.BuildRoute("hello", "name", "test", "number", "001")
	if string(url) != "/hello/test/001/world" {
		t.Error("Expecting /hello/test/001/world, received:", string(url))
	}
}
