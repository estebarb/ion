package hotcache

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

const ExpectedDuration = time.Millisecond * 50

func SlowResponse(w http.ResponseWriter, r *http.Request) {
	time.Sleep(ExpectedDuration)
	w.Header().Add("hello", "world")
	w.Write([]byte("hello"))
}

func TestConfig_Middleware(t *testing.T) {
	hc := New(time.Second * 3)
	h := hc.Middleware(http.HandlerFunc(SlowResponse))
	start := time.Now()

	ts := httptest.NewServer(h)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err := http.Get(ts.URL + "/")
			if err != nil {
				t.Error(err)
			}
			if req.Header.Get("hello") != "world" {
				t.Errorf("Expecting header hello:world, got %v", req.Header)
			}
			buf, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Error(err)
			}
			if string(buf) != "hello" {
				t.Errorf("Expecting 'hello', got %v", string(buf))
			}
			req.Body.Close()
		}()
	}
	wg.Wait()
	duration := time.Since(start)
	if duration < ExpectedDuration {
		t.Errorf("Expected duration was at least %v, takes %v", ExpectedDuration, duration)
	}

	if duration > ExpectedDuration+time.Millisecond*25 {
		t.Errorf("Expected duration was more or less %v, takes %v", ExpectedDuration, duration)
	}
}
