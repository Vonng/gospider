package gospider

/**************************************************************
* interface: Generator
**************************************************************/
// Generator is interface for ModuleTypeGenerator
// It could generate Request/Response/Item at any time ,
// for hooking into spider's control flow.

type Generator interface {
	Module
	// Generator_Gen will return a <-channel Data  and start a new routine sending values to it.
	// Generator may close channel, which means no more Data will be generated
	Gen() (<-chan Data)
	// It should return self as Generator
	Generator() Generator
}

/**************************************************************
* struct: defaultGenerator
**************************************************************/

// defaultGenerator will send request to bus
type defaultGenerator struct {
	ModuleInternal
	Ch chan Data
}

// NewGenerator will create a default generator
// if bufSize = 0 then it's an unbuffered channel underneath
func NewGenerator(bufSize uint) *defaultGenerator {
	if bufSize == 0 {
		// channel without buffer
		return &defaultGenerator{
			ModuleInternal: NewModuleInternalFromType(ModuleTypeGenerator),
			Ch:             make(chan Data),
		}
	} else {
		return &defaultGenerator{
			ModuleInternal: NewModuleInternalFromType(ModuleTypeGenerator),
			Ch:             make(chan Data, bufSize),
		}
	}
}

// NewGeneratorFromChan build new generator from existing channel
func NewGeneratorFromChan(c chan Data) *defaultGenerator {
	return &defaultGenerator{
		ModuleInternal: NewModuleInternalFromType(ModuleTypeGenerator),
		Ch:             c,
	}
}

// defaultGenerator_Generator wrap itself as Generator
func (self *defaultGenerator) Generator() Generator {
	return self
}

// defaultGenerator_Gen implement interface Generator
// it returns a chan, and can be used as :
// for data,ok := range generator.Gen() {}
func (self *defaultGenerator) Gen() (<-chan Data) {
	return (<-chan Data)(self.Ch)
}

// defaultGenerator_SendRequest will send raw request to bus
func (self *defaultGenerator) SendRequest(req *Request) {
	go func(c chan<- Data) {
		self.Ch <- req
	}(self.Ch)
}

// defaultGenerator_SendResponse will inject a response into spider
func (self *defaultGenerator) SendResponse(res *Response) {
	go func(c chan<- Data) {
		self.Ch <- res
	}(self.Ch)
}

// defaultGenerator_SendItem will inject an item into spider
func (self *defaultGenerator) SendItem(item Item) {
	go func(c chan<- Data) {
		self.Ch <- item
	}(self.Ch)
}

// defaultGenerator_SendURL will send url non-block to inner-channel
func (self *defaultGenerator) SendURL(url string) {
	req, _ := NewRequest("GET", url, nil, nil)
	self.SendRequest(req)
}

// defaultGenerator_SendURLs will start a new goroutine sending url as request to spider
func (self *defaultGenerator) SendURLs(urls []string) {
	go func(c chan<- Data) {
		for _, url := range urls {
			req, _ := NewRequest("GET", url, nil, nil)
			self.Ch <- req
		}
	}(self.Ch)
}

// defaultGenerator_SendFakeResponseWithURL will make fake response from content & url
func (self *defaultGenerator) SendFakeResponseString(url, content string) {
	self.SendResponse(FakeResponseString(url, content))
}

// defaultGenerator_SendFakeResponse will mock a response with given content
func (self *defaultGenerator) SendFakeResponse(url string, content []byte) {
	self.SendResponse(FakeResponse(url, content))
}

// defaultGenerator_Relay will concat Data chan
func (self *defaultGenerator) Relay(input <-chan Data) {
	for data := range input {
		self.Ch <- data
	}
}

// defaultGenerator_Close will close generator channel. VERY DANGEROUS
// Warning: writing to closed channel or close a closed channel will cause runtime panic
// So this function should only be called be producer only once after all payloads are send
func (self *defaultGenerator) Close() {
	close(self.Ch)
}
