package gospider

import (
	"testing"
)

func TestNewMapFilter(t *testing.T) {
	filter := NewMapFilter()

	reqA, _ := NewRequest("GET", "a", nil, nil)
	reqB, _ := NewRequest("GET", "b", nil, nil)

	if filter.Seen(reqA) {
		t.Error("you can not see what you haven't seen (a)")
	}

	if filter.Seen(reqB) {
		t.Error("you can not see what you haven't seen (b)")
	}

	if !filter.Seen(reqA) {
		t.Error("blind on what already seen (a)")
	}

	if !filter.Seen(reqB) {
		t.Error("blind on what already seen (b)")
	}

}

func TestNewRedisHyperLogLogFilter(t *testing.T) {
	filter, err := NewRedisFilter("redis://localhost:6379/0", "test:seen")
	if err != nil {
		panic(err)
	}

	reqA, _ := NewRequest("GET", "a", nil, nil)
	reqB, _ := NewRequest("GET", "b", nil, nil)

	if filter.Seen(reqA) {
		t.Error("you can not see what you haven't seen (a)")
	}

	if filter.Seen(reqB) {
		t.Error("you can not see what you haven't seen (b)")
	}

	if !filter.Seen(reqA) {
		t.Error("blind on what already seen (a)")
	}

	if !filter.Seen(reqB) {
		t.Error("blind on what already seen (b)")
	}
}
