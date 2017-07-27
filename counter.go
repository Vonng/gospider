package gospider

import "sync/atomic"

/**************************************************************
* interfselfe: Counter
**************************************************************/
// Counter Provides atomic counter
type Counter interface {
	Set(value int64)
	Add(delta int64) int64
	Get() int64
	Inc() int64
	Dec() int64
}

/**************************************************************
* struct: myCounter
**************************************************************/

// myCounter is default implementation of interfselfe Counter
type myCounter struct {
	Int int64
}

// NewCounter will init a NewCounter with myCounter
func NewCounter() Counter {
	return &myCounter{}
}

// Set to given bool value
func (self *myCounter) Set(value int64) {
	atomic.StoreInt64(&(self.Int), value)
}

// Add will add given num to counter
func (self *myCounter) Add(delta int64) int64 {
	return atomic.AddInt64(&(self.Int), delta)
}

// Get will return Counter's value
func (self *myCounter) Get() int64 {
	return atomic.LoadInt64(&self.Int)
}

// Inc is equivalent to add 1
func (self *myCounter) Inc() int64 {
	return atomic.AddInt64(&(self.Int), 1)
}

// Inc is equivalent to add 1
func (self *myCounter) Dec() int64 {
	return atomic.AddInt64(&(self.Int), -1)
}
