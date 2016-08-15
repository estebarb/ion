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
	r := New()
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
	router := New()
	router.Get("/hello/:name/:number/world",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			state := router.GetState(r)
			name, _ := state.Get("name")
			number, _ := state.Get("number")
			fmt.Fprintf(w, "/hello/%s/world/%s",
				name, number)
		})).Name("hello")
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
		http.NotFoundHandler()).Name("hello")
	url := router.BuildRoute("hello", "name", "test", "number", "001")
	if string(url) != "/hello/test/001/world" {
		t.Error("Expecting /hello/test/001/world, received:", string(url))
	}

	url = router.BuildRoute("void", "name", "test", "number", "001")
	if string(url) != "" {
		t.Error("Expecting empty string, got:", url)
	}

	url = router.BuildRoute("hello", "name", "test", "number")
	if string(url) != "" {
		t.Error("Expecting empty string, got:", url)
	}

	url = router.BuildRoute("hello", "name", "test")
	if string(url) != "" {
		t.Error("Expecting empty string, got:", url)
	}
}

func dummy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func TestNotFound(t *testing.T) {
	router := New()
	router.GetFunc("/abc", dummy)
	router.GetFunc("/abcd", dummy)
	router.GetFunc("/:name", dummy)
	router.GetFunc("/:name/xyz", dummy)

	ts := httptest.NewServer(router)

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

	router := New()
	router.Get("/", http.HandlerFunc(fun("GET")))
	router.Post("/", http.HandlerFunc(fun("POST")))
	router.Put("/", http.HandlerFunc(fun("PUT")))
	router.Delete("/", http.HandlerFunc(fun("DELETE")))
	router.Patch("/", http.HandlerFunc(fun("PATCH")))
	router.Options("/", http.HandlerFunc(fun("OPTIONS")))

	router.GetFunc("/", fun("GET"))
	router.PostFunc("/", fun("POST"))
	router.PutFunc("/", fun("PUT"))
	router.DeleteFunc("/", fun("DELETE"))
	router.PatchFunc("/", fun("PATCH"))
	router.OptionsFunc("/", fun("OPTIONS"))

	ts := httptest.NewServer(router)

	for _, method := range []string{http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions} {
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
