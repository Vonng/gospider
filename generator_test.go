package gospider

import (
	"testing"
	"fmt"
	"time"
)

func TestNewURLRequestGenerator(t *testing.T) {
	g := NewGenerator(0)
	g.SendURLs([]string{"http://www.baidu.com", "http://www.wandoujia.com/"})

	go func() {
		time.Sleep(10 * time.Second)
		g.Close()
	}()

	for data := range g.Generator().Gen() {
		switch v := data.(type) {
		case *Request:
			fmt.Println(data.(*Request).URL)
			fmt.Println(v)
		case *Response:
			fmt.Println(v)
		case Item:
			fmt.Println(v)
		default:
			fmt.Println("not a valid data")
		}
	}
}
