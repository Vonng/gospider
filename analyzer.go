package gospider

/**************************************************************
* interface: Analyzer
**************************************************************/

// Parser is function which takes Response and yields New Requests/Items
type Parser func(res *Response) ([]Data, error)

// ContentReader is a naive parser that read html body into item["content"]
func ContentReader(res *Response) ([]Data, error) {
	if content, err := res.Content(); err != nil {
		return nil, err
	} else {
		item := Item{"content": content.String()}
		// copy request meta to item
		for k, v := range res.Meta {
			item[k] = v
		}
		return []Data{item}, nil
	}
}

// Analyzer is interface for module Analyzer
// Analyzer take Response as input, gives Requests & Items
type Analyzer interface {
	Module
	// Get parser according to parser's name
	GetParser(name string) Parser
	// Analyze 会根据Request与Response中的元数据调用指定的Parser解析
	Analyze(res *Response) ([]Data, error)
}

/**************************************************************
* struct: defaultAnalyzer
**************************************************************/

type defaultAnalyzer struct {
	ModuleInternal
	// Contains named parsers, Read only, No-concurrent support
	parsers       map[string]Parser
	defaultParser Parser
}

func NewAnalyzer(parsers map[string]Parser) Analyzer {
	// if map contains a parser named "default", then it is the default parser
	var defaultParser Parser
	if p, ok := parsers["default"]; ok {
		defaultParser = p
	} else if len(parsers) == 1 {
		// if map contains only one parser, then whatever it called, it's the default parser
		for _, v := range parsers {
			defaultParser = v
		}
	}

	return &defaultAnalyzer{
		ModuleInternal: NewModuleInternalFromType(ModuleTypeAnalyzer),
		parsers:        parsers,
		defaultParser:  defaultParser,
	}

	//
}

// defaultAnalyzer_GetParser will get parser by name
func (self *defaultAnalyzer) GetParser(name string) Parser {
	if p, ok := self.parsers[name]; ok {
		return p
	}
	return nil
}

// defaultAnalyzer_SetParser should only be used when init
func (self *defaultAnalyzer) SetParser(name string, parser Parser) {
	self.parsers[name] = parser
}

// defaultAnalyzer_Analyze will parse response and yield request & items
func (self *defaultAnalyzer) Analyze(res *Response) ([]Data, error) {
	// use default parser by default
	parser := self.defaultParser
	// if callback is manually set, then use corresponding parser
	if callback, ok := res.Meta["callback"]; ok {
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
