package gospider

import (
	"fmt"
	"strings"
	"net/http"
)

/**************************************************************
* struct: Response
**************************************************************/

// Response hold http.Response and corresponding Request
type Response struct {
	*http.Response
	Request *Request
}

// NewResponse create a new response and attach origin request to it
func NewResponse(res *http.Response, req *Request) *Response {
	return &Response{Response: res, Request: req}
}

// Response_Repr implement Data interface
func (res *Response) Repr() string {
	return fmt.Sprintf("%+v", res)
}

/**************************************************************
* type: fakeBody(string)
**************************************************************/

// fakeBody implement io.ReadCloser
// used for testing response.body
type fakeBody struct {
	reader *strings.Reader
}

func (fb fakeBody) Read(b []byte) (n int, err error) {
	return fb.reader.Read(b)
}

func (r fakeBody) Close() error {
	return nil
}

// FakeBody build a io.ReadCloser from given string
func FakeBody(s string) fakeBody {
	return fakeBody{strings.NewReader(s)}
}

// FakeResponse : make response from string body
func FakeResponse(url, content string) *Response {
	//return &Response{Res: res, Meta: meta}
	req, _ := NewRequest("GET", url, nil, nil)
	return &Response{
		Response: &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Header:        make(http.Header, 0),
			ContentLength: int64(len(content)),
			Request:       req.Request,
			Body:          FakeBody(content),
		},
		Request: req,
	}
}

// FakeResponse : make response from string body with additional meta attached
func FakeResponseMeta(url, content string, metaMap MetaMap) *Response {
	//return &Response{Res: res, Meta: meta}
	req, _ := NewRequest("GET", url, nil, metaMap)
	return &Response{
		Response: &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Header:        make(http.Header, 0),
			ContentLength: int64(len(content)),
			Request:       req.Request,
			Body:          FakeBody(content),
		},
		Request: req,
	}
}
