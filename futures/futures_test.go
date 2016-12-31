package futures

import (
	"testing"
	"time"
)

var counter = 0

func generator() interface{} {
	time.Sleep(time.Second)
	counter++
	return counter
}

func TestNewFutureFunc(t *testing.T) {
	before := time.Now()
	valueF := NewFutureFunc(generator)
	afterCreate := time.Now()
	x := valueF.Read()
	readTime := time.Now()

	if afterCreate.Sub(before) > time.Second/2 {
		t.Error("NewFutureFunc should not block, but takes", afterCreate.Sub(before))
	}

	// The timings aren't exact, so we must take that into account
	if readTime.Sub(afterCreate) < time.Second-(time.Second/100) {
		t.Error("Reading the value should have taken 1 "+
			"second, but takes: ", readTime.Sub(afterCreate))
	}

	y := Read(valueF)

	if x.(int) != 1 {
		t.Error("Future value must be 1")
	}

	if y.(int) != 1 {
		t.Error("Future value must be 1")
	}
}
