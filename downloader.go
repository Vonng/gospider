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

// NewDownloader will create a new downloader from given id
func NewDownloader(client *http.Client, filter Filter) Downloader {
	if client == nil {
		client = new(http.Client)
	}

	return &defaultDownloader{
		ModuleInternal: NewModuleInternalFromType(ModuleTypeDownloader),
		Client:         client,
		Filter:         filter,
	}
}

// defaultDownloader_Download will download response from given request
func (self *defaultDownloader) Download(req *Request) (res *Response, err error) {
	self.ModuleInternal.Call()

	if req == nil || req.Request == nil {
		return nil, ErrNilRequest
	}

	if self.Filter != nil {
		// if bool field "filter" occurs and value is false
		var disableFilter = false
		if filter, ok := req.Meta["filter"]; ok {
			if filter.(bool) == false {
				disableFilter = true
			}
		}
		// by default: filter will check whether it's duplicate request
		if !disableFilter {
			if self.Seen(req) {
				log.Errorf("[$s] dupe request: %s", self.ID(), req.URL)
				return nil, ErrDupeRequest
			}
		}
	}

	self.Doing()
	log.Infof("[%s] URL: %s", self.ID(), req.URL)
	httpRes, err := self.Do(req.Request)

	if err != nil {
		self.Fail()
		return NewResponse(httpRes, req.Meta), err
	} else {
		self.Done()
		return NewResponse(httpRes, req.Meta), nil
	}
}
