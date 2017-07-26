package gospider

import (
	"testing"
	"fmt"
)

func TestDefaultDownloader_Download(t *testing.T) {
	downloader := NewDefaultDownloader()

	reqMeta := MetaMap{"hello": "world", "proxy": "http://localhost:8848"}
	req, err := NewRequest("GET",
		"http://www.baidu.com",
		nil,
		reqMeta,
	)

	if err != nil {
		t.Error("create new request failed")
	}

	res, err := downloader.Download(req)
	if err != nil || !res.Valid() {
		t.Error("download failed. %#v %#v", req, res)
	}

	for k, v := range res.Meta {
		if reqMeta[k] != v {
			t.Error("Request & Response should have same meta")
		}
	}

}

func TestFakeResponse(t *testing.T) {
	res := FakeResponseString("http://www.baidu.com", "what the heck is that")
	fmt.Println(res.Content())
}

func TestDulicateRequest(t *testing.T) {
	downloader, _ := NewDownloader(nil, NewMapFilter())

	req, err := NewRequest("GET",
		"http://www.baidu.com",
		nil,
		MetaMap{"hello": "world", "proxy": "http://localhost:8848"})

	if err != nil {
		t.Error("create new request failed")
	}

	res, err := downloader.Download(req)
	if err != nil || !res.Valid() {
		t.Error("download failed. %#v %#v", req, res)
	}

	res, err = downloader.Download(req)

	if err == nil {
		t.Error("downloading same request should raise ErrDuplicateRequest")
	}

}
