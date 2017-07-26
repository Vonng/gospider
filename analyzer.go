package gospider

import "io/ioutil"

/**************************************************************
* interface: Analyzer
**************************************************************/

// Parser is function which takes Response and yields New Requests/Items
type Parser func(res *Response) ([]Data, error)

// Parser map is a map holding parsers
type ParserMap map[string]Parser

// Analyzer is interface for module Analyzer
// Analyzer take Response as input and yield multiple Request or Item
type Analyzer interface {
	Module
	// Get parser according to parser's name
	GetParser(name string) Parser
	// Analyze will invoke specific parser for parsing according to res.Meta
	Analyze(res *Response) ([]Data, error)
}

/**************************************************************
* struct: defaultAnalyzer
**************************************************************/

// defaultAnalyzer is default implementation of interface Analyzer
type defaultAnalyzer struct {
	ModuleInternal
	// parsers Contains named parsers, Read only
	parsers       ParserMap
	defaultParser Parser
}

// NewAnalyzer will init a new Analyzer will a parserMap
// if parsers == nil, then BodyReader is used as default analyzer
func NewAnalyzer(parsers ParserMap) (Analyzer, error) {
	// if map contains a parser named "default", then it is the default parser
	var defaultParser Parser
	// if no parser is provided, use content reader as default parser and new a ParserMap
	if parsers == nil {
		return nil, ErrNilParser
	}

	if p, ok := parsers[KeyDefault]; ok {
		defaultParser = p
	} else if len(parsers) == 1 {
		// if map contains only one parser, whatever it called, it is the default parser
		for _, v := range parsers {
			if v != nil {
				defaultParser = v
			} else {
				return nil, ErrNilParser
			}
		}
	}

	return &defaultAnalyzer{
		ModuleInternal: NewModuleInternalFromType(ModuleTypeAnalyzer),
		parsers:        parsers,
		defaultParser:  defaultParser,
	}, nil
}

// NewAnalyzerFromParser is one-parser only version of Analyzer constructor
func NewAnalyzerFromParser(parser Parser) (Analyzer, error) {
	if parser == nil {
		return nil, ErrNilParser
	}
	return &defaultAnalyzer{
		ModuleInternal: NewModuleInternalFromType(ModuleTypeAnalyzer),
		parsers:        ParserMap{KeyDefault: parser},
		defaultParser:  parser,
	}, nil
}

// defaultAnalyzer_GetParser will get parser by name
func (self *defaultAnalyzer) GetParser(name string) Parser {
	if p, ok := self.parsers[name]; ok {
		return p
	}
	return nil
}

// defaultAnalyzer_Analyze will parse response and yield request & items
func (self *defaultAnalyzer) Analyze(res *Response) ([]Data, error) {
	// use default parser by default
	parser := self.defaultParser
	// if callback is manually set, then use corresponding parser
	if callback, ok := res.Meta[KeyCallback]; ok {
		callbackName, ok := callback.(string)
		if !ok {
			return nil, ErrInvalidCallback
		}
		if p := self.GetParser(callbackName); p == nil {
			return nil, ErrCallbackNotFount
		} else {
			parser = p
		}
	}

	return parser(res)
}

// BodyReader is a naive parser that read html body([]byte) into item["content"]
// and also copy all kv in request's meta to item (cautious: do not use "content" as key)
// This can be used when no analyzer is given
func BodyReader(res *Response) ([]Data, error) {
	item := make(Item, len(res.Meta)+1)
	if body, err := ioutil.ReadAll(res.Response.Body); err != nil {
		return nil, err
	} else {
		// copy request meta to item
		for k, v := range res.Meta {
			item[k] = v
		}

		// copy content([]byte) to item["content"]. it will overwrite meta's content (if set)
		item[KeyBody] = body
		return []Data{item}, nil
	}
}
