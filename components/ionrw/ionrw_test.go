package ionrw

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var TESTKEY = "TEST"
var TESTVAL = "SUCCESS"
var TESTBODY = "Hello World!"

func precondition(w http.ResponseWriter, r *http.Request) {

	w.Header().Add(TESTKEY, TESTVAL)
}

func example(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			wr := Interceptor(w, r)
			wr.AddPrecondition(precondition)
			next.ServeHTTP(w, r)
		})
}

func TestResponseWriter_AddPrecondition(t *testing.T) {
	dummy := func(w http.ResponseWriter, r *http.Request) {
		wr := Interceptor(w, r)
		if wr.Size() != 0 {
			t.Errorf("Expected size 0, got: %v", wr.Size())
		}
		if wr.Written() {
			t.Error("Expecting unwritten response")
		}

		w.Write([]byte(TESTBODY))

		if wr.Size() != len(TESTBODY) {
			t.Errorf("Expected size %v, got: %v",
				len(TESTBODY),
				wr.Size())
		}
		if !wr.Written() {
			t.Error("Expecting written response")
		}
	}

	s := httptest.NewServer(
		InterceptorMiddleware(
			example(
				http.HandlerFunc(dummy))))
	res, err := http.Get(s.URL)

	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	if res.Header.Get(TESTKEY) != TESTVAL {
		t.Errorf("Expecting header <%s:%s>, received <%s:%s>",
			TESTKEY, TESTVAL,
			TESTKEY, res.Header.Get(TESTKEY))
	}

	received := string(body)
	if received != TESTBODY {
		t.Errorf("Expecting <%s>, received <%s>",
			TESTBODY,
			received)
	}
}
