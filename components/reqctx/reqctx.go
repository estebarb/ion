
// +build !go1.7

// reqctx provides an uniform way to access request context.
package reqctx

import (
	"golang.org/x/net/context"
	"net/http"
	"sync"
)

var contexts = make(map[*http.Request]*context.Context)
var contextsLock sync.Mutex

func Context(r *http.Request) *context.Context {
	contextsLock.Lock()
	defer contextsLock.Unlock()
	return contexts[r]
}

func ContextMiddleware(base *context.Context) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := base
				contextsLock.Lock()
				contexts[r] = ctx
				contextsLock.Unlock()

				defer func(){
					contextsLock.Lock()
					delete(contexts, r)
					contextsLock.Unlock()
				}()

				next.ServeHTTP(w, r)
			})
	}
}