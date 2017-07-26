package gospider

import (
	"testing"
	"math"
)

func TestSerial(t *testing.T) {
	s := NewSerial()
	if s == nil {
		t.Errorf("couldn't create new serial!")
	}

	if s.Min() != 0 {
		t.Errorf("serial's init min value not equal to zero!")
	}

	if s.Max() != math.MaxUint64 {
		t.Errorf("serial's init max value not equal to math.MaxUint64!")
	}

	if s.Peek() != 0 {
		t.Errorf("serial's init peek() != 0")
	}

	if s.Next() != 0 {
		t.Errorf("serial's first next() != 0")
	}

	if s.Next() != 1 {
		t.Errorf("serial's second next() != 1")
	}

	if s.Peek() != 2 {
		t.Error("serial's second peek() != 1")
	}

	if s.Cycle() != 0 {
		t.Error("serial's Cycle != 0")
	}

	if s.Set(math.MaxUint64); s.Next() != math.MaxUint64 {
		t.Error("serial set value failed")
	}

	if s.Next() != s.Min() && s.Cycle() != 1 {
		t.Error("serial's recycle failed")
	}

	s = NewSerialRange(100, 101)

	if s.Next() != 100 {
		t.Error("serial next() start ne 100")
	}

	if s.Next() != 101 {
		t.Error("serial next() second ne 101")
	}

	if s.Next() != 100 && s.Cycle() != 1 {
		t.Error("serial next() second ne 101")
	}
}

func TestNamedSerial(t *testing.T) {
	s1 := NewNamedSerial("s1")

	if s1.Next() != 0 {
		t.Error("s1 init next() ne 1")
	}

	if s1.Peek() != 1 {
		t.Error("s1 after next(), peek() should be 1")
	}

	// get serial variable by name. refer to same serial as s1
	s1ref := GetNamedSerial("s1")
	if s1ref.Peek() != 1 {
		t.Error("s1ref peek() should also be 1")
	}

	s1ref.Next()
	s1ref.Next()

	if s1ref.Peek() != 3 {
		t.Error("s1ref next() again, peek ne 2")
	}

	if s1.Peek() != 3 {
		t.Error("s1 peek() should also be 2")
	}
}
