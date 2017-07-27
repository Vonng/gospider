package gospider

/**************************************************************
* interface: Scheduler
**************************************************************/
// Scheduler arrange request , remove duplicate request
type Scheduler interface {
	Seen(req *Request) bool

	// These should be block method
	PutRequest(req *Request) bool
	GetRequest() *Request
	LenRequests() int

	GetResponse() *Response
	PutResponse(res *Response)
	LenResponses() int

	GetItem() Item
	PutItem(item Item)
	LenItems() int

	SendData(datum Data)
	SendDataList(datum []Data)

	Pull(<-chan Data)

	Idle() bool
}

/**************************************************************
* struct: myScheduler
**************************************************************/
// myScheduler is default implement of interface Scheduler
type myScheduler struct {
	Filter
	Requests  chan *Request
	Responses chan *Response
	Items     chan Item
	Errors    chan error
}

// NewScheduler will create a new scheduler from given id
func NewScheduler(reqBufSz, resBufSz, itemBufSz uint, filter Filter) Scheduler {
	return &myScheduler{
		Filter:    filter,
		Requests:  make(chan *Request, reqBufSz),
		Responses: make(chan *Response, resBufSz),
		Items:     make(chan Item, itemBufSz),

	}
}

// myScheduler_PutRequest will download response from given request
// it will check duplicate accroding to req.Meta
// return value indicate whether this request is enqueued
func (self *myScheduler) PutRequest(req *Request) bool {
	if !req.IgnoreDupe && self.Filter != nil && self.Seen(req) {
		self.Errors <- ErrDupeRequest
		return false
	}
	self.Requests <- req
	return true
}

// myScheduler_GetRequest will fetch a request from chan
// block method
func (self *myScheduler) GetRequest() *Request {
	return <-self.Requests
}

func (self *myScheduler) LenRequests() int {
	return len(self.Requests)
}

// myScheduler_PutResponse will download Response from given request
func (self *myScheduler) PutResponse(res *Response) {
	if res != nil {
		self.Responses <- res
	}
}

// myScheduler_GetResponse will fetch a Response from chan
// block method
func (self *myScheduler) GetResponse() *Response {
	return <-self.Responses
}

func (self *myScheduler) LenResponses() int {
	return len(self.Responses)
}

// myScheduler_PutItem will download Item from given request
func (self *myScheduler) PutItem(item Item) {
	if item != nil {
		self.Items <- item
	}
}

// myScheduler_GetItem will fetch a Item from chan
// block method
func (self *myScheduler) GetItem() Item {
	return <-self.Items
}

func (self *myScheduler) LenItems() int {
	return len(self.Items)
}

func (self *myScheduler) SendData(datum Data) {
	switch v := datum.(type) {
	case *Request:
		self.PutRequest(v)
	case *Response:
		self.PutResponse(v)
	case Item:
		self.PutItem(v)
	}
}

func (self *myScheduler) SendDataList(datum []Data) {
	if datum != nil {
		for _, data := range datum {
			self.SendData(data)
		}
	}
}

func (self *myScheduler) Pull(generator <-chan Data) {
	go func() {
		for datum := range generator {
			self.SendData(datum)
		}
	}()
}

func (self *myScheduler) Idle() bool {
	return self.LenItems() == 0 && self.LenRequests() == 0 && self.LenResponses() == 0
}
