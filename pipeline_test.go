package gospider

import (
	"testing"
	"fmt"
)

func TestNewPipeline(t *testing.T) {
	p1 := func(item Item) error {
		fmt.Println("p1")
		return nil
	}

	p2 := func(item Item) error {
		fmt.Println("p2")
		return nil
	}

	p3 := func(item Item) error {
		fmt.Println("p3")
		if item["flag"].(bool) == true {
			item["flagSetToTrue"] = true
			return nil
		}
		return ErrDropItem
	}

	pipe, _ := NewPipeline([]Processor{p1, p2, p3, p2, p1})

	fmt.Println(pipe)
	i := map[string]interface{}{
		"flag":  false,
		"data":  "xixi",
		"hello": 3,
		"url":   "http://www.baidu.com",
	}

	fmt.Printf("Before %#v\n", i)
	pipe.Send(Item(i))
	fmt.Printf("After %#v\n", i)

}
