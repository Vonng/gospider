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
	GetParser(name string) Parser
	// Analyze will invoke specific parser for parsing according to res.Meta
	Analyze(res *Response) ([]Data, error)
}

/**************************************************************
* struct: myAnalyzer
**************************************************************/

// myAnalyzer is default implementation of interface Analyzer
type myAnalyzer struct {
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

	return &myAnalyzer{
		parsers:       parsers,
		defaultParser: defaultParser,
	}, nil
}

// NewAnalyzerSolo is one-parser only version of Analyzer constructor
func NewAnalyzerSolo(parser Parser) (Analyzer, error) {
	if parser == nil {
		return nil, ErrNilParser
	}
	return &myAnalyzer{
		parsers:       ParserMap{KeyDefault: parser},
		defaultParser: parser,
	}, nil
}

// myAnalyzer_GetParser will get parser by name
func (self *myAnalyzer) GetParser(name string) Parser {
	if p, ok := self.parsers[name]; ok {
		return p
	}
	return nil
}

// myAnalyzer_Analyze will parse response and yield request & items
func (self *myAnalyzer) Analyze(res *Response) ([]Data, error) {
	// use default parser by default
	callback := self.defaultParser

	// if callback is manually set and could be found, then use it
	if res.Request != nil && res.Request.Callback != "" {
		if callback = self.GetParser(res.Request.Callback); callback == nil {
			callback = self.defaultParser
		}
	}

	return callback(res)
}

/**************************************************************
* Parser: BodyReader
**************************************************************/
// BodyReader is a naive parser that read html body([]byte) into item["content"]
// and also copy all kv in request's meta to item (cautious: do not use "content" as key)
// This can be used when no analyzer is given
func BodyReader(res *Response) ([]Data, error) {
	item := make(Item, len(res.Request.Meta)+1)
	if body, err := ioutil.ReadAll(res.Response.Body); err != nil {
		return nil, err
	} else {
		// copy request meta to item
		for k, v := range res.Request.Meta {
			item[k] = v
		}

		// copy content([]byte) to item["content"]. it will overwrite meta's content (if set)
		item[KeyBody] = body
		return []Data{item}, nil
	}
}
