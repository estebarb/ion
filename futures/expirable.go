package futures

import (
	"math/rand"
	"sync"
	"time"
)

type Expirable struct {
	sync.Mutex
	timeout   time.Duration
	timestamp time.Time
	value     interface{}
	generator func() interface{}
}

func NewExpirable(timeout time.Duration, generator func() interface{}) *Expirable {
	return &Expirable{
		timeout:   timeout,
		value:     nil,
		generator: generator,
	}
}

func (e *Expirable) Read() interface{} {
	return e.ReadFunc(e.generator)
}

func (e *Expirable) ReadFunc(f func() interface{}) interface{} {
	e.Lock()
	defer e.Unlock()
	if time.Since(e.timestamp) > e.timeout || e.value == nil {
		e.value = f()
		e.timestamp = time.Now()
	}
	return e.value
}

type ExpirablePool struct {
	Expirables []*Expirable
}

func NewExpirablePool(timeout time.Duration, size int, generator func() interface{}) *ExpirablePool {
	pool := &ExpirablePool{
		Expirables: make([]*Expirable, size),
	}
	for i := range pool.Expirables {
		pool.Expirables[i] = NewExpirable(timeout, generator)
	}
	return pool
}

func (p ExpirablePool) Pick() *Expirable {
	return p.Expirables[rand.Intn(len(p.Expirables))]
}

func (p *ExpirablePool) Read() interface{} {
	return p.Pick().Read()
}

func (p *ExpirablePool) ReadFunc(f func() interface{}) interface{} {
	return p.Pick().ReadFunc(f)
}
