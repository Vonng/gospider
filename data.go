package gospider

import (
	"io"
	"bytes"
	"net/http"
)

/**************************************************************
* interface: Data
**************************************************************/

// Data is interface for Request/Response/Item
type Data interface {
	// Valid will tells whether it's a valid Data entry
	Valid() bool
}

/**************************************************************
* struct: Item
**************************************************************/

// Item holds result, but it's just a container
// the interpretation rely on user themselves
type Item map[string]interface{}

// Item_Valid implement Data interface
func (item Item) Valid() bool {
	return item != nil
}

// Maybe serialize & deserialize should be part of item definition

/**************************************************************
* struct: MetaMap
**************************************************************/
// MetaMap holds meta info that attach to request & response
type MetaMap map[string]interface{}

/**************************************************************
* struct: Request
**************************************************************/

// Request extend http.Request with additional meta info.
// Request implemented interface Data
type Request struct {
	*http.Request
	Meta MetaMap
}

// NewRequest create new request with http.Request and meta
func NewRequest(method, urlStr string, body io.Reader, meta MetaMap) (*Request, error) {
	if meta == nil {
		meta = make(MetaMap, 0)
	}
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	return &Request{
		Request: req,
		Meta:    meta,
	}, nil
}

// Request_Valid will tell whether request is valid
// implement interface Data.Valid
func (req *Request) Valid() bool {
	return req.Request != nil && req.Request.URL != nil
}

/**************************************************************
* struct: Response
**************************************************************/

// Response hold http.Response and corresponding meta. implemented interface Data
type Response struct {
	*http.Response
	Meta MetaMap
}

// NewResponse will create new response and attach request's meta
func NewResponse(res *http.Response, meta map[string]interface{}) *Response {
	return &Response{Response: res, Meta: meta}
}

// Response_Content will read response.body with a buffer holding it
// Can't be called twice
func (res *Response) Content() (buf *bytes.Buffer, err error) {
	if res.Response == nil || res.Body == nil {
		return nil, ErrNilResponse
	}
	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	res.Body.Close()
	return
}

// Response_Valid implement Data interface
func (res *Response) Valid() bool {
	return res.Response != nil && res.Body != nil
}

/**************************************************************
* type: fakeBody(string)
**************************************************************/

// fakeBody implement io.ReadCloser.
type fakeBody struct {
	*bytes.Buffer
}

func (self *fakeBody) Close() (error) {
	return nil
}

// NewFakeBody build io.ReadCloser from []byte
func NewFakeBody(body []byte) io.ReadCloser {
	return &fakeBody{
		Buffer: bytes.NewBuffer(body),
	}
}

// NewFakeBodyString build io.ReadCloser from string
func NewFakeBodyString(body string) io.ReadCloser {
	return &fakeBody{
		Buffer: bytes.NewBufferString(body),
	}
}

// FakeResponse : make response from string body
func FakeResponseString(url, content string) *Response {
	//return &Response{Res: res, Meta: meta}
	req, _ := http.NewRequest("GET", url, nil)
	return &Response{
		Response: &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Header:        make(http.Header, 0),
			ContentLength: int64(len(content)),
			Request:       req,
			Body:          NewFakeBodyString(content),
		},
		Meta: make(map[string]interface{}, 0),
	}
}

// FakeResponse : make response from []byte body
func FakeResponse(url string, content []byte) *Response {
	//return &Response{Res: res, Meta: meta}
	req, _ := http.NewRequest("GET", url, nil)
	return &Response{
		Response: &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Header:        make(http.Header, 0),
			ContentLength: int64(len(content)),
			Request:       req,
			Body:          NewFakeBody(content),
		},
		Meta: make(map[string]interface{}, 0),
	}
}
