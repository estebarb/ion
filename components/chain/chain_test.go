package chain

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func GenerateHandlers(tag string) http.Handler {
	return http.HandlerFunc(GenerateHandlerFunc(tag))
}

func GenerateHandlerFunc(tag string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(tag))
	}
}

func GenerateMiddleware(tag string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tag))
			next.ServeHTTP(w, r)
			w.Write([]byte(tag))
		})
	}
}

func TestAddThen(t *testing.T) {
	chain := New()
	chain.Add(GenerateMiddleware("aaa"))
	handler := chain.Then(GenerateHandlers("bbb"))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, nil)
	if rec.Body.String() != "aaabbbaaa" {
		t.Error("Wrong response. Expecting aaabbaaa, received", rec.Body.String())
	}
}

func TestAddThenFunc(t *testing.T) {
	chain := New()
	chain.Add(GenerateMiddleware("aaa"))
	handler := chain.ThenFunc(GenerateHandlerFunc("bbb"))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, nil)
	if rec.Body.String() != "aaabbbaaa" {
		t.Error("Wrong response. Expecting aaabbaaa, received", rec.Body.String())
	}
}

func TestSequence(t *testing.T) {
	handler := Sequence(GenerateHandlers("a"),
		GenerateHandlers("b"),
		GenerateHandlers("c"),
		GenerateHandlers("d"),
	)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, nil)
	if rec.Body.String() != "abcd" {
		t.Error("Wrong response. Expecting abcd, received", rec.Body.String())
	}
}
