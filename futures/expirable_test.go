package futures

import (
	"testing"
	"time"
)

func TestNewExpirable(t *testing.T) {
	ttl := time.Millisecond * 100
	exp := NewExpirable(ttl, func() interface{} {
		return time.Now()
	})
	t1 := exp.Read().(time.Time)
	time.Sleep(ttl / 2)
	t2 := exp.Read().(time.Time)

	if t1 != t2 {
		t.Errorf("Expected %v as second read, given %v", t1, t2)
	}

	time.Sleep(ttl)
	t3 := exp.Read().(time.Time)
	if t1 == t3 {
		t.Errorf("Expected t1 != t3 but %v == %v", t1, t3)
	}
}

func TestNewExpirablePool(t *testing.T) {
	ttl := time.Millisecond * 100
	exp := NewExpirablePool(ttl, 1, func() interface{} {
		return time.Now()
	})
	t1 := exp.Read().(time.Time)
	time.Sleep(ttl / 2)
	t2 := exp.Read().(time.Time)

	if t1 != t2 {
		t.Errorf("Expected %v as second read, given %v", t1, t2)
	}

	time.Sleep(ttl)
	t3 := exp.Read().(time.Time)
	if t1 == t3 {
		t.Errorf("Expected t1 != t3 but %v == %v", t1, t3)
	}
}

func TestExpirable_ReadFunc(t *testing.T) {
	ttl := time.Millisecond * 100
	fun := func() interface{} {
		return time.Now()
	}
	fun2 := func() interface{} {
		return 0
	}
	exp := NewExpirablePool(ttl, 1, fun)
	t1 := exp.Read().(time.Time)
	time.Sleep(ttl / 2)
	t2 := exp.ReadFunc(fun2).(time.Time)

	if t1 != t2 {
		t.Errorf("Expected %v as second read, given %v", t1, t2)
	}

	time.Sleep(ttl)
	t3, ok := exp.ReadFunc(fun2).(int)
	if !ok || t3 != 0 {
		t.Errorf("Expected t3 == 0 but is %v", t3)
	}
}
