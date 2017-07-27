package gospider

import log "github.com/Sirupsen/logrus"

/**************************************************************
* struct:  EngineArgs
**************************************************************/

// EngineArgs holds args of engine
type EngineArgs struct {
	// Spider name
	Name string

	// Request Dupe Filter. Set to nil to disable
	Filter Filter

	// Analyzer instance
	Analyzer Analyzer

	// Downloader instance
	Downloader Downloader

	// Pipeline instance
	Pipeline Pipeline

	// downloader work count. should not be zero
	DWorkers uint32

	// ReqBufSize sets request chan buffer size. set to a small number is ok
	ReqBufSize uint32

	// ResBufSize is response chan buffer size. set to zero to be on-demand-spawn
	ResBufSize uint32

	// ResBufSize is item chan buffer size. set to zero to be on-demand-spawn
	ItemBufSize uint32

	// ErrBufSize could be set to a proper number like 1000
	ErrBufSize uint32
}

// Default presets
func NewEngineArgs() *EngineArgs {
	return &EngineArgs{
		Filter:      NewMapFilter(),
		DWorkers:    20,
		ReqBufSize:  10,
		ResBufSize:  10000,
		ItemBufSize: 10000,
		ErrBufSize:  10000,
	}
}

/**************************************************************
* interface:  Engine
**************************************************************/
// Engine: core of spider

type Engine interface {
	Scheduler
	Run(<-chan Data) <-chan error
	Summary() string
}

type myEngine struct {
	myScheduler
	Name       string
	Args       *EngineArgs
	Analyzer   Analyzer
	Downloader Downloader
	Pipeline   Pipeline
}

func NewEngine(args *EngineArgs) Engine {
	engine := &myEngine{
		myScheduler: myScheduler{
			Filter:    args.Filter,
			Requests:  make(chan *Request, args.ReqBufSize),
			Responses: make(chan *Response, args.ResBufSize),
			Items:     make(chan Item, args.ItemBufSize),
			Errors:    make(chan error, args.ErrBufSize),
		},
		Args:       args,
		Analyzer:   args.Analyzer,
		Downloader: args.Downloader,
		Pipeline:   args.Pipeline,
	}
	return engine
}

func (self *myEngine) Run(generator <-chan Data) <-chan error {
	log.Info("[INIT] engine starting...")
	self.analyze()
	self.pipeline()
	self.download()

	if generator != nil {
		self.Pull(generator)
	}

	return (<-chan error)(self.Errors)
}

func (self *myEngine) Stop([]Data) error {
	log.Info("[INIT] engine stopping...[Not implemented]")
	return nil
}

func (self *myEngine) Summary() string {
	return "not implemented"
}

func (self *myEngine) download() {
	var n uint32
	if n = self.Args.DWorkers; n == 0 {
		// means on-demand-spawn goroutines // dangerous
		log.Infof("[INIT] DWorker = 0. spawn goroutine for each request.")
		for req := range self.Requests {
			go self.downloadReq(req)
		}

	} else {
		log.Infof("[INIT] DWorker = %d. spawn %d download goroutine", n, n)
		ParallelWork(n, self.downloadLoop)
	}
	log.Infof("[INIT] Downloader init complete")
}

func (self *myEngine) downloadLoop(id int) {
	log.Infof("[INIT] Downloader[id:%d] routine init", id)
	for req := range self.Requests {
		log.Debugf("[DOWN]-[%d] fetch %s ", id, req.URL)
		res, err := self.Downloader.Download(req)
		if err != nil {
			self.Errors <- err
		} else {
			self.Responses <- res
			log.Infof("[DOWN][%d] done %s ", id, req.URL)
		}
	}
}

// myEngine_downloadReq : on-demand-spawn
func (self *myEngine) downloadReq(req *Request) {

	log.Infof("[DOWN] %s begin", req.URL)
	res, err := self.Downloader.Download(req)
	if err != nil {
		self.Errors <- err
	}
	// Put Response
	self.Responses <- res
	log.Infof("[DOWN] %s complete", req.URL)
}

// analyze start an analyze loop (ODS: on demand spawn)
func (self *myEngine) analyze() {
	log.Infof("[INIT] Analyzer init begin")
	go func() {
		for res := range self.Responses {
			if res == nil {
				self.Errors <- ErrNilResponse
			} else {
				go self.parseOne(res)
			}
		}
	}()
	log.Infof("[INIT] Analyzer init complete")
}

func (self *myEngine) parseOne(res *Response) {
	log.Info("[ANAY] parser one item")
	data, err := self.Analyzer.Analyze(res)
	if len(data) > 0 {
		self.SendDataList(data)
	} else {
		log.Warn("[ANAY] parse with no yield")
	}

	if err != nil {
		self.Errors <- err
	}
}

func (self *myEngine) pipeline() {
	log.Infof("[INIT] Pipeline init begin")
	go func() {
		for item := range self.Items {
			if item == nil {
				self.Errors <- ErrNilItem
			} else {
				go self.pickOne(item)
			}
		}
	}()
	log.Infof("[INIT] Pipeline init complete")
}

func (self *myEngine) pickOne(item Item) {
	log.Info("[PIPE] pick item %s", item.Repr())
	errs := self.Pipeline.Send(item)
	if len(errs) > 0 {
		for _, err := range errs {
			self.Errors <- err
		}
	}
}
