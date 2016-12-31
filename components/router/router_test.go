package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testEq(t *testing.T, a, b []string) {
	if len(a) != len(b) {
		t.Error("len(a) != len(b):", len(a), "vs", len(b))
	}
	for k, v := range a {
		if v != b[k] {
			t.Errorf("Item %d: <%s> != <%s>", k, v, b[k])
		}
	}
}

func TestSplitWithoutTrailingSlash(t *testing.T) {
	a := "/hello/world"
	b := "/hello/world/"
	x := splitWithoutTrailingSlash(a)
	y := splitWithoutTrailingSlash(b)
	testEq(t, x, y)
}

func TestRouter_Get(t *testing.T) {
	r := NewDefault()
	r.GetFunc("/hello",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello"))
		}).Name("hello")
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
	r := NewDefault()
	r.Get("/hello/:name/:number/world",
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			state := r.Context(req).(IPathParam)
			name, _ := state.PathParams()["name"]
			number, _ := state.PathParams()["number"]
			fmt.Fprintf(w, "/hello/%s/world/%s",
				name, number)
		})).Name("hello")
	ts := httptest.NewServer(r)

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
	r := NewDefault()
	r.Get("/hello/:name/:number/world",
		http.NotFoundHandler()).Name("hello")
	url := r.RouteFor("hello", "name", "test", "number", "001")
	if string(url) != "/hello/test/001/world" {
		t.Error("Expecting /hello/test/001/world, received:", string(url))
	}

	url = r.RouteFor("void", "name", "test", "number", "001")
	if string(url) != "" {
		t.Error("Expecting empty string, got:", url)
	}

	url = r.RouteFor("hello", "name", "test", "number")
	if string(url) != "" {
		t.Error("Expecting empty string, got:", url)
	}

	url = r.RouteFor("hello", "name", "test")
	if string(url) != "" {
		t.Error("Expecting empty string, got:", url)
	}
}

func dummy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func TestNotFound(t *testing.T) {
	r := NewDefault()
	r.GetFunc("/abc", dummy)
	r.GetFunc("/abcd", dummy)
	r.GetFunc("/:name", dummy)
	r.GetFunc("/:name/xyz", dummy)

	ts := httptest.NewServer(r)

	resp, err := http.Get(ts.URL + "/abc")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Expecting StatusOK, received", resp.Status)
	}

	resp, err = http.Get(ts.URL + "/abcd")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Expecting StatusOK, received", resp.Status)
	}

	resp, err = http.Get(ts.URL + "/hello")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Expecting StatusOK, received", resp.Status)
	}

	resp, err = http.Get(ts.URL + "/hello/xyz")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Expecting StatusOK, received", resp.Status)
	}

	resp, err = http.Get(ts.URL + "/hello/asdf")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("Expecting StatusNotFound, received", resp.Status)
	}

	resp, err = http.Get(ts.URL + "/abc/asdf")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("Expecting StatusNotFound, received", resp.Status)
	}
}

func TestHandlers(t *testing.T) {
	fun := func(method string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != method {
				t.Errorf("Expecting %s, got %s",
					method,
					r.Method)
			}
			w.Write([]byte(method))
		}
	}

	r := NewDefault()
	r.Get("/", http.HandlerFunc(fun("GET")))
	r.Post("/", http.HandlerFunc(fun("POST")))
	r.Put("/", http.HandlerFunc(fun("PUT")))
	r.Delete("/", http.HandlerFunc(fun("DELETE")))
	r.Patch("/", http.HandlerFunc(fun("PATCH")))
	r.Options("/", http.HandlerFunc(fun("OPTIONS")))

	r.GetFunc("/", fun("GET"))
	r.PostFunc("/", fun("POST"))
	r.PutFunc("/", fun("PUT"))
	r.DeleteFunc("/", fun("DELETE"))
	r.PatchFunc("/", fun("PATCH"))
	r.OptionsFunc("/", fun("OPTIONS"))

	ts := httptest.NewServer(r)

	for _, method := range []string{"GET",
		"POST",
		"PUT",
		"DELETE",
		"PATCH",
		"OPTIONS"} {
		req, err := http.NewRequest(method, ts.URL, nil)
		if err != nil {
			t.Error(err)
		}
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected StatusOK, got %s", response.Status)
		}
	}
}
