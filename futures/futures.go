// futures allow adding incomplete computations
// in contexts and templates.
// This allows to start the rendering of the
// template before all the other operations
// (for example, database queries) had been finished.
package futures

import "sync"

// Future represents an maybe incomplete operation, that is
// being processed
type Future struct {
	sync.Once
	input chan interface{}
	value interface{}
}

// NewFuture creates a Future from an input channel
func NewFuture(input chan interface{}) *Future {
	f := &Future{
		input: input,
		value: nil,
	}
	return f
}

// NewFutureFunc creates a Future from a function that
// returns an interface{}
func NewFutureFunc(f func()interface{}) *Future{
	c := make(chan interface{})
	go func(){
		c <- f()
	}()
	return NewFuture(c)
}

// Read blocks until the computation finish,
// and then returns the value.
func (f *Future) Read() interface{} {
	f.Do(func(){
		f.value = <- f.input
	})
	return f.value
}

// Read reads the Future value. Is intented to
// be used on templates.
func Read(f *Future) interface{} {
	return f.Read()
}