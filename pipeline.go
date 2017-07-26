package gospider

import "fmt"

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
	Module
	// Send : Send Item to pipeline。
	Send(item Item) []error
}

/**************************************************************
* defaultPipeline: Pipeline
**************************************************************/

// defaultPipeline is default implementation of interface Pipeline
type defaultPipeline struct {
	ModuleInternal
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
		ModuleInternal: NewModuleInternalFromType(ModuleTypePipeline),
		processors:     list,
	}, nil
}

// NewPipelineFromProcessor create pipeline from single processor
func NewPipelineFromProcessor(processor Processor) (Pipeline, error) {
	if processor == nil {
		return nil, ErrNilProcessor
	}
	return &defaultPipeline{
		ModuleInternal: NewModuleInternalFromType(ModuleTypePipeline),
		processors:     []Processor{processor},
	}, nil
}

// defaultPipeline_Send will put item into pipeline for handling
func (self *defaultPipeline) Send(item Item) []error {
	self.ModuleInternal.Call()

	var errs []error
	if item == nil {
		errs = append(errs, ErrNilItem)
		return errs
	}
	self.ModuleInternal.Doing()

	for _, processor := range self.processors {
		err := processor(item)
		if err != nil {
			errs = append(errs, err)
			// if returns a ErrDropItem then break directly
			if err == ErrDropItem {
				break
			}
		}
	}

	if len(errs) == 0 {
		self.Done()
	} else {
		self.Fail()
	}

	return errs
}

// PrintItem is simplest pipeline
func PrintItem(item Item) error {
	fmt.Println("%+v", item)
	return nil
}
