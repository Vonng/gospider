package gospider

import (
	"net/http"
)

/**************************************************************
* interface: Downloader
**************************************************************/
// Downloader is interface of downloader module
type Downloader interface {
	Module
	// Download receive Request and returns Response
	Download(req *Request) (*Response, error)
}

/**************************************************************
* struct: defaultDownloader
**************************************************************/
type defaultDownloader struct {
	ModuleInternal
	*http.Client
	Filter
}

// NewDefaultDownloader return a simple downloader without error
func NewDefaultDownloader() Downloader {
	return &defaultDownloader{
		ModuleInternal: NewModuleInternalFromType(ModuleTypeDownloader),
		Client:         &http.Client{},
		Filter:         nil,
	}
}

// NewDownloader will create a new downloader from given id
func NewDownloader(client *http.Client, filter Filter) (Downloader, error) {
	if client == nil {
		client = new(http.Client)
	}

	return &defaultDownloader{
		ModuleInternal: NewModuleInternalFromType(ModuleTypeDownloader),
		Client:         client,
		Filter:         filter,
	}, nil
}

// defaultDownloader_Download will download response from given request
func (self *defaultDownloader) Download(req *Request) (res *Response, err error) {
	self.ModuleInternal.Call()
	if req == nil || req.Request == nil {
		return nil, ErrNilRequest
	}

	// Dupe filter
	// by default filter is disabled when Download does not have a filter
	if self.Filter != nil {
		// by default, all request will pass dupe filter. except they explict set Meta["_filter"] = false
		if filter, ok := req.Meta[KeyFilter]; !(ok && filter.(bool) == false) {
			// need filter
			if self.Seen(req) {
				return nil, ErrDupeRequest
			}
		}
	}

	self.Doing()
	httpRes, err := self.Do(req.Request)

	if err != nil {
		self.Fail()
		return NewResponse(httpRes, req.Meta), err
	} else {
		self.Done()
		return NewResponse(httpRes, req.Meta), nil
	}
}
