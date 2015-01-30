/*
Package context provides a shared storage for values
used during the lifetime of a request.

It can be used for passing values between middleware,
handlers, templates, etc.

Because it exports a Context type each middleware can
keep a private context for each request, and share it as
necessary. We could use package especific functions to
retrieve the data, without contaminating a global context.
In this way, we also solve the eventual problem that one middleware
can overwrite the data of others.

It is based on https://groups.google.com/forum/#!msg/golang-nuts/teSBtPvv1GQ/U12qA9N51uIJ and
Gorilla Toolkit Context.
 */
package context

import (
	"net/http"
	"sync"
)

type Context struct{
	sync.RWMutex
	data map[*http.Request]map[interface{}]interface{}
}

/*
Creates a new context
 */
func New()Context{
		return Context{
			data: make(map[*http.Request]map[interface{}]interface{}),
		}
}

var globalContext Context

/*
Clears the context data associated to the given request.
 */
func (c *Context)Clear(r *http.Request){
	c.Lock()
	defer c.Unlock()
	delete(c.data, r)
}

/*
Returns a http.Handler that automatically clears the data associated
to the request, after the handler has been processed.
 */
func (c *Context) ClearHandler(h http.Handler) http.Handler{
	fn := func(w http.ResponseWriter, r *http.Request){
		defer c.Clear(r)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

/*
Sets a value associated with a request in this context.
 */
func (c *Context) Set(r *http.Request, key interface{}, value interface{}){
	c.Lock()
	defer c.Unlock()
	
	if _, exists := c.data[r]; !exists{
		c.data[r] = make(map[interface{}]interface{})
	}
	c.data[r][key] = value
}

/*
Deletes a key-value from the context associated to the given request.
 */
func (c *Context) Delete(r *http.Request, key interface{}){
	c.Lock()
	defer c.Unlock()
	
	delete(c.data[r], key)
}

/*
Return a value from the context associated to the given request.
 */
func (c *Context) Get(r *http.Request, key interface{}) (interface{}, bool){
	c.RLock()
	defer c.RUnlock()
	data, ok := c.data[r][key]
	return data, ok
}

/*
Returns all the values from the context associated to the given request.
 */
func (c *Context) GetAll(r *http.Request) (map[interface{}]interface{}, bool){
	c.RLock()
	defer c.RUnlock()
	data, ok := c.data[r]
	return data, ok
}

/*
Clears the context data associated to the given request.
 */
func Clear(r *http.Request){
	globalContext.Lock()
	defer globalContext.Unlock()
	delete(globalContext.data, r)
}

/*
Returns a http.Handler that automatically clears the data associated
to the request, after the handler has been processed.
 */
func ClearHandler(h http.Handler) http.Handler{
	fn := func(w http.ResponseWriter, r *http.Request){
		defer globalContext.Clear(r)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

/*
Sets a value associated with a request in this context.
 */
func Set(r *http.Request, key interface{}, value interface{}){
	globalContext.Lock()
	defer globalContext.Unlock()

	if _, exists := globalContext.data[r]; !exists{
		globalContext.data[r] = make(map[interface{}]interface{})
	}
	globalContext.data[r][key] = value
}

/*
Deletes a key-value from the context associated to the given request.
 */
func Delete(r *http.Request, key interface{}){
	globalContext.Lock()
	defer globalContext.Unlock()

	delete(globalContext.data[r], key)
}

/*
Return a value from the context associated to the given request.
 */
func Get(r *http.Request, key interface{}) (interface{}, bool){
	globalContext.RLock()
	defer globalContext.RUnlock()
	data, ok := globalContext.data[r][key]
	return data, ok
}

/*
Returns all the values from the context associated to the given request.
 */
func GetAll(r *http.Request) (map[interface{}]interface{}, bool){
	globalContext.RLock()
	defer globalContext.RUnlock()
	data, ok := globalContext.data[r]
	return data, ok
}
