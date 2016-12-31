package futures

import (
	"math/rand"
	"sync"
	"time"
)

// Expirable manages the retrieval and recalculation of a value that expires
// after some time
type Expirable struct {
	sync.Mutex
	timeout   time.Duration
	timestamp time.Time
	value     interface{}
	generator func() interface{}
}

// NewExpirable creates a new Expirable with the given TTL and generated with
// the given generator
func NewExpirable(timeout time.Duration, generator func() interface{}) *Expirable {
	return &Expirable{
		timeout:   timeout,
		value:     nil,
		generator: generator,
	}
}

// Read returns the current value of the Expirable, and if it has
// expired then it recalculates the value using the default generator
func (e *Expirable) Read() interface{} {
	return e.ReadFunc(e.generator)
}

// ReadFunc returns the current value of the Expirable, and if it has
// expired then it recalculates the value using the given generator
func (e *Expirable) ReadFunc(f func() interface{}) interface{} {
	e.Lock()
	defer e.Unlock()
	if time.Since(e.timestamp) > e.timeout || e.value == nil {
		e.value = f()
		e.timestamp = time.Now()
	}
	return e.value
}

// ExpirablePool manages the retrieval and recalculation of a pool of values
// that expires after some time.
type ExpirablePool struct {
	Expirables []*Expirable
}

// NewExpirablePool creates a pool of several Expirable, that could be requested
// concurrently avoiding a bit the overhead of locking a single Expirable for
// all the requests.
// Note that potentially the returned values could go out of synchronization.
func NewExpirablePool(timeout time.Duration, size int, generator func() interface{}) *ExpirablePool {
	pool := &ExpirablePool{
		Expirables: make([]*Expirable, size),
	}
	for i := range pool.Expirables {
		pool.Expirables[i] = NewExpirable(timeout, generator)
	}
	return pool
}

// Pick selects an Expirable at random from the pool
func (p ExpirablePool) Pick() *Expirable {
	return p.Expirables[rand.Intn(len(p.Expirables))]
}

// Read reads the value from an Expirable picked at random, and if it
// has expired then it recalculates the value using the default generator
func (p *ExpirablePool) Read() interface{} {
	return p.Pick().Read()
}

// ReadFunc reads the value from an Expirable picked at random, and if it
// has expired then it recalculates the value using the given generator
func (p *ExpirablePool) ReadFunc(f func() interface{}) interface{} {
	return p.Pick().ReadFunc(f)
}
