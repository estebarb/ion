package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func compare(t *testing.T, path, expected string) {
	ts := httptest.NewServer(NewApp())
	defer ts.Close()
	response, err := http.Get(ts.URL + path)
	if err != nil {
		t.Error(err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	strbody := string(body)

	if strbody != expected {
		t.Errorf("Expected '%s', given '%s'", expected, strbody)
	}
}

func TestNewApp(t *testing.T) {
	compare(t, "/", "Hello world!")
	compare(t, "/asdf", "Hello, asdf!")
}
