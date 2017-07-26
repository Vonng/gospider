package gospider

import (
	"math"
	"sync"
)

/**************************************************************
* interface: Serial
**************************************************************/

// Serial
type Serial interface {
	// Min will return lower bound of serial
	Min() uint64

	// Max will return upper bound of serial
	Max() uint64

	// Set will set next value to given param
	Set(uint64)

	// Next will fetch next value and advanced
	Next() uint64

	// Peek will read next value without advance
	Peek() uint64

	// Cycle will return cycle count
	Cycle() uint64
}

/**************************************************************
* struct: defaultSerial
**************************************************************/

// defaultSerial implement interface Serial
type defaultSerial struct {
	min   uint64
	max   uint64
	next  uint64
	cycle uint64
	lock  sync.RWMutex
}

// NewSerial create a defaultSerial with range 0 ~ UINT64_MAX
func NewSerial() Serial {
	return &defaultSerial{
		min:  0,
		max:  math.MaxUint64,
		next: 0,
	}
}

// NewSerialRange create a Serial with given range min ~ max
// when given max = 0, max is set to math.MaxUint64
func NewSerialRange(min uint64, max uint64) Serial {
	if max == 0 {
		max = math.MaxUint64
	}
	return &defaultSerial{
		min:  min,
		max:  max,
		next: min,
	}
}

// defaultSerial_Min get serial's lower bound
func (s *defaultSerial) Min() uint64 {
	return s.min
}

// defaultSerial_Max get serial's upper bound
func (s *defaultSerial) Max() uint64 {
	return s.max
}

// defaultSerial_set set next value to given value
func (s *defaultSerial) Set(i uint64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.next = i
}

// defaultSerial_Next get next value and advance
func (s *defaultSerial) Next() uint64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	id := s.next
	if id == s.max {
		// Reset serial and advance cycle
		s.next = s.min
		s.cycle++
	} else {
		s.next++
	}
	return id
}

// defaultSerial_Next see next value without advance
func (s *defaultSerial) Peek() uint64 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.next
}

// defaultSerial_Cycle will fetch cycle count
func (s *defaultSerial) Cycle() uint64 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.cycle
}

/**************************************************************
* Serial Registry. NamedSerial can be register and fetch by name
**************************************************************/
var (
	serialRegistry = make(map[string]Serial, 5)
	registryLock   sync.RWMutex
)

// [PUBLIC]
// NewNamedSerialRange will create a new named serial with given range
func NewNamedSerialRange(name string, min, max uint64) Serial {
	s := NewSerialRange(min, max)
	registryLock.Lock()
	defer registryLock.Unlock()
	serialRegistry[name] = s
	return s
}

// NewNamedSerial will create a default named serial
func NewNamedSerial(name string) Serial {
	s := NewSerial()
	registryLock.Lock()
	defer registryLock.Unlock()
	serialRegistry[name] = s
	return s
}

// GetNamedSerial will lookup serial by name
func GetNamedSerial(name string) Serial {
	registryLock.RLock()
	defer registryLock.RUnlock()
	return serialRegistry[name]
}

// SetNamedSerial will register a serial with given name
func SetNamedSerial(name string, s Serial) {
	registryLock.Lock()
	defer registryLock.Unlock()
	serialRegistry[name] = s
}
