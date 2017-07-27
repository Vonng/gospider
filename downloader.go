package gospider

import "net/http"


/**************************************************************
* interface: Downloader
**************************************************************/
// Downloader is interface of downloader module
type Downloader interface {
	Download(req *Request) (*Response, error)
}

/**************************************************************
* struct: myDownloader
**************************************************************/
// myDownloader is default implement of interface Downloader
type myDownloader struct {
	*http.Client
}

// NewDownloader will create a new downloader from given id
func NewDownloader(client *http.Client) (Downloader, error) {
	if client == nil {
		client = new(http.Client)
	}

	return &myDownloader{
		Client: client,
	}, nil
}

// myDownloader_Download will download response from given request
func (self *myDownloader) Download(req *Request) (res *Response, err error) {
	if req == nil && req.URL == nil {
		return nil, ErrNilRequest
	}

	httpRes, err := self.Do(req.Request)
	if err != nil {
		return NewResponse(httpRes, req), err
	} else {
		return NewResponse(httpRes, req), nil
	}
}

// ParallelWork start a series number of worker
func ParallelWork(n uint32, worker func(id int)) {
	for i := uint32(1); i <= n; i ++ {
		go worker(int(i))
	}
}
