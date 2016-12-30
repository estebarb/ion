package reqctx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	sc := New(BasicContextFactory)
	handler := sc.Middleware

	dummy := func(w http.ResponseWriter, r *http.Request) {
		state := sc.Context(r).(map[string]interface{})
		for k, v := range state {
			if v2, exists := state[k]; !exists || v.(string) != v2.(string) {
				t.Error("GetAll key/values doesn't correspond with Get results")
			}
		}
		_, exists := state["k2"]
		if exists {
			t.Error("k2 should have been deleted by step m2")
		}
	}

	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			state := sc.Context(r).(map[string]interface{})
			state["k1"] = "true"
			state["k2"] = "true"
			state["k3"] = "true"
			r = sc.WithContext(state, r)
			next.ServeHTTP(w, r)
		})
	}

	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			state := sc.Context(r).(map[string]interface{})
			for _, key := range []string{"k1", "k2"} {
				field, exists := state[key]
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
			delete(state, "k2")
			r = sc.WithContext(state, r)
			next.ServeHTTP(w, r)
		})
	}

	s := httptest.NewServer(handler(m1(m2(http.HandlerFunc(dummy)))))

	http.Get(s.URL)

	if sc.size() != 0 {
		t.Error("Expecting empty StateContainer, but have ",
			sc.size(),
			"states")
	}
}
