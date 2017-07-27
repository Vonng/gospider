package gospider

import (
	"net/http"
	"io"
	"fmt"
)

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
	Meta     MetaMap
	Callback string
	Errback  string
	// IgnoreDupe = true will ignore Dupe filter when scheduled
	IgnoreDupe bool
	Priority   int32
}

// NewRequest create new request with http.Request and meta
func NewRequest(method, urlStr string, body io.Reader, meta MetaMap) (req *Request, err error) {
	if meta == nil {
		meta = make(MetaMap, 2)
	}

	if httpReq, err := http.NewRequest(method, urlStr, body); err == nil {
		req = &Request{
			Request: httpReq,
			Meta:    meta,
		}
	}
	return
}

// NewGetRequest set method to GET and omit body & meta
func NewGetRequest(url string) (*Request, error) {
	return NewRequest("GET", url, nil, nil)
}

// Request_Repr will tell whether request is valid
// implement interface Data.Valid
func (req *Request) Repr() string {
	return fmt.Sprintf("(REQ:%s)", req.URL)
}

// Request_Data wraps itself as Data
func (req *Request) Data() Data {
	return req
}

// Request_DataList wraps itself as a slice of Data
// useful for those functions who yield only one Item
func (req *Request) DataList() []Data {
	return []Data{req}
}

// Request_SetCallback will set callback name to Meta
func (req *Request) SetCallback(parserName string) *Request {
	req.Callback = parserName
	return req
}

// Request_SetErrback will set callback name to Meta
func (req *Request) SetErrback(errbackName string) *Request {
	req.Callback = errbackName
	return req
}

// Request_SetPriority  will set callback name to Meta
func (req *Request) SetPriority(prior int32) *Request {
	req.Priority = prior
	return req
}

func (req *Request) DisableFilter() *Request {
	req.IgnoreDupe = true
	return req
}
