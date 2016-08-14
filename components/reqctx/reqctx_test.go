package reqctx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	sc := NewStateContainer()
	handler := sc.Middleware

	dummy := func(w http.ResponseWriter, r *http.Request) {
		state := sc.GetState(r)
		for k, v := range state.GetAll() {
			if v2, exists := state.Get(k); !exists || v.(string) != v2.(string) {
				t.Error("GetAll key/values doesn't correspond with Get results")
			}
		}
		_, exists := state.Get("k2")
		if exists {
			t.Error("k2 should have been deleted by step m2")
		}
	}

	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			state := sc.GetState(r)
			state.Set("k1", "true")
			state.Set("k2", "true")
			state.Set("k3", "true")
			next.ServeHTTP(w, r)
		})
	}

	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			state := sc.GetState(r)
			for _, key := range []string{"k1", "k2"} {
				field, exists := state.Get(key)
				if !exists {
					t.Errorf("Expecting %s=true, but %s doesn't exists",
						key, key,
						field)
				}
				value, ok := field.(string)
				if !ok || value != "true" {
					t.Errorf("Expecting %s=true, got %s=%v",
						key, key,
						field)
				}
			}
			state.Delete("k2")
			next.ServeHTTP(w, r)
		})
	}

	s := httptest.NewServer(handler(m1(m2(http.HandlerFunc(dummy)))))

	http.Get(s.URL)

	if len(sc.data) != 0 {
		t.Error("Expecting empty StateContainer, but have ",
			len(sc.data),
			"states")
	}
}
