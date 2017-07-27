package gospider

import (
	"testing"
)

func TestNewPipeline(t *testing.T) {
	p1 := func(item Item) error {
		item["p1"] = "p1"
		return nil
	}

	p2 := func(item Item) error {
		item["p2"] = "p2"
		return ErrDropItem
	}

	p3 := func(item Item) error {
		item["p3"] = "p3"
		return nil
	}

	pipe, _ := NewPipeline([]Processor{p1, p2, p3})
	// have "p1", "p2", but not have p3

	i := make(Item, 3)
	pipe.Send(Item(i))

	if i["p1"] != "p1" {
		t.Error("pipeline should have first processor executed")
	}

	if i["p2"] != "p2" {
		t.Error("pipeline should have second processor executed")
	}

	if _, ok := i["p3"]; ok {
		t.Error("pipeline should drop item after p2, so there shouldn't be p3")
	}

}
