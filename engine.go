package gospider

/**************************************************************
* struct:  EngineArgs
**************************************************************/

// EngineArgs holds args of engine
type EngineArgs struct {
	Generator       Generator
	Analyzer        Analyzer
	Downloader      Downloader
	Pipeline        Pipeline
	Workers         uint32
	RequestBufSize  uint32
	ResponseBufSize uint32
	ItemBufSize     uint32
	ErrorBufSize    uint32
}

// Engine: core of spider
type Engine struct {
	generator  Generator
	analyzer   Analyzer
	downloader Downloader
	pipeline   Pipeline
	Workers    uint32
	Requests   chan *Request
	Responses  chan *Response
	Items      chan Item
	Errors     chan error
}

// NewEngine build a engine from given args
func NewEngine(args *EngineArgs) *Engine {
	engine := &Engine{
		generator:  args.Generator,
		analyzer:   args.Analyzer,
		downloader: args.Downloader,
		pipeline:   args.Pipeline,
		Workers:    args.Workers,
		Requests:   make(chan *Request, args.RequestBufSize),
		Responses:  make(chan *Response, args.ResponseBufSize),
		Items:      make(chan Item, args.ItemBufSize),
		Errors:     make(chan error, args.ErrorBufSize),
	}
	return engine
}

// Engine_Run will start engine. returns a ErrorChannel
func (self *Engine) Run() <-chan error {
	log.Info("[INIT] engine health check...")

	log.Info("[INIT] prepare downloader")
	if self.downloader != nil {
		self.download()
		log.Info("[INIT] downloader ready")
	} else {
		log.Warn("[INIT] nil downloader")
	}

	log.Info("[INIT] prepare analyzer")
	if self.analyzer != nil {
		self.analyze()
		log.Warn("[INIT] analyzer ready")
	} else {
		log.Warn("[INIT] nil analyzer")
	}

	if self.pipeline != nil {
		log.Info("[INIT] prepare pipeline")
		self.pipework()
		log.Warn("[INIT] analyzer ready")
	} else {
		log.Warn("[INIT] nil pipeline. set a pipeline with PrintItem")
		self.pipeline, _ = NewPipelineFromProcessor(PrintItem)
	}

	if self.generator != nil {
		log.Info("[INIT] prepare generator")
		self.generate()
		log.Info("[INIT] generator ready")
	} else {
		log.Warn("[INIT] nil generator")
	}

	log.Info("[INIT] engine build complete success")
	return (<-chan error)(self.Errors)
}

// Engine_download will start download loop with n worker
func (self *Engine) download() {
	for i := 0; uint32(i) < self.Workers; i ++ {
		log.Infof("[INIT] downloader worker [%d] init", i)
		go func(worker int) {
			for {
				self.downloadOne()
			}
		}(i)
	}
}

// Engine_downloadOne will pick a request and yield a response
func (self *Engine) downloadOne() {
	// Get one request from pool
	req, ok := <-self.Requests
	if !ok {
		self.Errors <- ErrTrashInRequestPool
	}

	// send it to downloader
	res, err := self.downloader.Download(req)
	if err != nil || res == nil {
		self.Errors <- err
	} else {
		self.Responses <- res
	}
}

// Engine_analyzer will begin a analyze loop
func (self *Engine) analyze() {
	go func() {
		for res := range self.Responses {
			if res.Valid() {
				go self.parseOne(res)
			}
		}
	}()
}

// Engine_parseOne will take one item from item chan and parse it
func (self *Engine) parseOne(res *Response) {
	data, err := self.analyzer.Analyze(res)
	if data != nil && len(data) > 0 {
		for _, itemOrReq := range data {
			switch v := itemOrReq.(type) {
			case Item:
				self.Items <- v
			case *Request:
				self.Requests <- v
			case *Response:
				self.Errors <- ErrResponseFromAnalyzer
			}
		}
	}

	if err != nil {
		self.Errors <- err
	}
}

// Engine_pipework will start item processing loop
func (self *Engine) pipework() {
	go func() {
		for item := range self.Items {
			go self.pickOne(item)
		}
	}()
}

// Engine_pickOne will get an item from item chan and handle it
func (self *Engine) pickOne(item Item) {
	go func() {
		errs := self.pipeline.Send(item)
		if len(errs) > 0 {
			for _, err := range errs {
				self.Errors <- err
			}
		}
	}()
}

// Engine_generate: start working loop:
// Data -> item/req/res pool
func (self *Engine) generate() {
	go func() {
		for datum := range self.generator.Generator().Gen() {
			switch v := datum.(type) {
			case Item:
				self.Items <- v
				log.Debugf("[GEN] new item %+v", v)
			case *Request:
				self.Requests <- v
				log.Debugf("[GEN] new request %+v", v)
			case *Response:
				self.Responses <- v
				log.Debugf("[GEN] new response %+v", v)
			default:
				self.Errors <- ErrGenerateInvalidType
			}
		} // quit when generator is closed
	}()
}
