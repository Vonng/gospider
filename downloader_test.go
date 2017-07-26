package gospider

import (
	"testing"
	"fmt"
)

func TestDefaultDownloader_Download(t *testing.T) {
	dld := NewDownloader(nil, nil)

	req, err := NewRequest("GET",
		"http://www.baidu.com",
		nil,
		MetaMap{"hello": "world", "proxy": "http://localhost:8848"})

	if err != nil {
		t.Error("create new request failed")
	}

	res, err := dld.Download(req)
	if err != nil || !res.Valid() {
		t.Error("download failed. %#v %#v", req, res)
	}

	b, _ := res.Content()
	fmt.Printf("%+v %+v \n", b, res.Meta)

}

func TestFakeResponse(t *testing.T) {
	res := FakeResponseString("http://www.baidu.com", "what the heck is that")
	fmt.Println(res.Content())
}

func TestDulicateRequest(t *testing.T) {
	dld := NewDownloader(nil, nil)

	req, err := NewRequest("GET",
		"http://www.baidu.com",
		nil,
		MetaMap{"hello": "world", "proxy": "http://localhost:8848"})

	if err != nil {
		t.Error("create new request failed")
	}

	res, err := dld.Download(req)
	if err != nil || !res.Valid() {
		t.Error("download failed. %#v %#v", req, res)
	}

	res, err = dld.Download(req)

	if err == nil {
		t.Error("downloading same request should raise ErrDuplicateRequest")
	}

}
