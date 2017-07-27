package gospider

/**************************************************************
* interface: processor
**************************************************************/

// Processor is a function which take item as params
// when pass is true, error could be omitted
type Processor func(item Item) error

/**************************************************************
* interface: Pipeline
**************************************************************/
// Pipeline take Item in and handle it
type Pipeline interface {
	// Send will make an item go through pipeline
	// pipe will interrupt when ErrDropItem is returned by processor
	Send(item Item) []error
}

/**************************************************************
* defaultPipeline: Pipeline
**************************************************************/

// defaultPipeline is default implementation of interface Pipeline
type defaultPipeline struct {
	processors []Processor
}

// NewPipeline create a default pipeline
func NewPipeline(processors []Processor) (Pipeline, error) {
	if processors == nil || len(processors) == 0 {
		return nil, ErrNilProcessor
	}

	var list []Processor
	for _, processor := range processors {
		if processor == nil {
			return nil, ErrNilProcessor
		}
		list = append(list, processor)
	}

	return &defaultPipeline{
		processors: list,
	}, nil
}

// NewPipelineSolo create pipeline from a solo processor
// this constructor do not check processor == nil
func NewPipelineSolo(processor Processor) (Pipeline) {
	return &defaultPipeline{[]Processor{processor}}
}

// defaultPipeline_Send will put item into pipeline for handling
// nil item will not be checked
func (self *defaultPipeline) Send(item Item) []error {
	// normal errors will just be collected together except ErrDropItem
	var errs []error
	for _, processor := range self.processors {
		err := processor(item)
		if err != nil {
			errs = append(errs, err)
			if err == ErrDropItem {
				break
			}
		}
	}

	return errs
}
