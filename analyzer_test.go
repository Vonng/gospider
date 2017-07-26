package gospider

import (
	"testing"
	"fmt"
)

func TestNewAnalyzer(t *testing.T) {

	res := FakeResponseString("http://www.baidu.com", "what the heck is that?")

	fmt.Println(res.Content())
	parsers := map[string]Parser{
		"default": func(res *Response) ([]Data, error) {
			content, err := res.Content()
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(content.String())

			return []Data{Item{"content": "hahahah"}}, nil
		},
	}
	//
	ana := NewAnalyzer(parsers)

	datas, err := ana.Analyze(res)
	if err != nil {
		panic(err)
	}
	fmt.Println((datas[0]).(Item))

}
