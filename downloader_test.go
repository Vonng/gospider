package gospider

import "testing"

func TestNewDownloader(t *testing.T) {
	downloader, err := NewDownloader(nil)
	if err != nil {
		t.Error(err)
	}

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
	if err != nil || res == nil {
		t.Error("download failed. %#v %#v", req, res)
	}

	for k, v := range res.Request.Meta {
		if reqMeta[k] != v {
			t.Error("Request & Response should have same meta")
		}
	}
}
