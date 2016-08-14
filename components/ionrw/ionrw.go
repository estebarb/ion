// ionrw contains a http.ResponseWriter implementation
package ionrw

import (
	"net/http"
	"sync"
)

// ResponseWriter is a custom implementation of http.ResponseWriter
// that supports registering callbacks, that are called
// before writting the final response.
type ResponseWriter struct {
	http.ResponseWriter
	status        int
	size          int
	preconditions []http.HandlerFunc
	request       *http.Request
	once          sync.Once
}

// New creates a new *ResponseWriter based on the passed
// http.ResponseWriter and *http.Request
func New(w http.ResponseWriter, r *http.Request) *ResponseWriter {
	rw := &ResponseWriter{}
	rw.ResponseWriter = w
	rw.status = 0
	rw.size = 0
	rw.preconditions = make([]http.HandlerFunc, 0)
	rw.request = r
	return rw
}

// callPreconditions executes all the preconditions registered
// in this ResponseWriter.
func (rw *ResponseWriter) callPreconditions() {
	for _, fun := range rw.preconditions {
		fun(rw, rw.request)
	}
}

// Written returns if this ResponseWriter have been already
// written.
func (rw *ResponseWriter) Written() bool {
	return rw.status != 0
}

// Size returns the amount of bytes written to this ResponseWriter
// body.
func (rw *ResponseWriter) Size() int {
	return rw.size
}

// AddPrecondition adds a http.HandlerFunc that will be called
// before writting headers to the response.
func (rw *ResponseWriter) AddPrecondition(fun http.HandlerFunc) {
	rw.preconditions = append(rw.preconditions, fun)
}

// WriteHeader writes a header to the ResponseWriter, but before
// calls all the preconditions registered.
func (rwi *ResponseWriter) WriteHeader(status int) {
	rwi.status = status
	rwi.once.Do(rwi.callPreconditions)
	rwi.ResponseWriter.WriteHeader(status)
}

// Write writes the bytes passed to this ResponseWriter, but before
// calls all the preconditions registered.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.once.Do(rw.callPreconditions)
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// Interceptor returns or creates a *ResponseWriter associated
// with the passed http.ResponseWriter and *http.Request
func Interceptor(w http.ResponseWriter, r *http.Request) *ResponseWriter {
	if w, ok := w.(*ResponseWriter); ok {
		return w
	}
	return New(w, r)
}

func InterceptorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rw := Interceptor(w, r)
			next.ServeHTTP(rw, r)
		})
}
