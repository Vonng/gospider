package wdj_app

import (
	"testing"
	"fmt"
	. "github.com/Vonng/gospider"
)

func TestNewWdjAppRequestGenerator(t *testing.T) {
	generator, err := RequestGenerator("redis://localhost:6379/0")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(generator)

	for d := range generator {
		fmt.Println(d.(*Request).URL)
	}
}
