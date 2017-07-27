package gospider

import "testing"

func TestNewAnalyzer(t *testing.T) {
	fakeURL := "https://www.google.com"
	fakeContent := "test content"
	// make a fake response
	res := FakeResponse(fakeURL, fakeContent)

	// NewAnalyzer will init a Analyzer will ContentReader as parser
	analyzer, err := NewAnalyzerSolo(BodyReader)
	if err != nil {
		t.Error(err)
	}

	datums, err := analyzer.Analyze(res)
	if err != nil {
		t.Error(err)
	}

	if len(datums) < 1 {
		t.Error("default Analyzer should generate at least one Item in datums contains http body")
	}

	item, ok := datums[0].(Item)
	if !ok {
		t.Error("default Analyzer should generate a item rather than other Data entry")
	}

	contentStr, ok := item[KeyBody]
	if !ok {
		t.Error(`default Analyzer's yield item should have field "body"`)
	}

	content, ok := contentStr.([]byte);
	if !ok {
		t.Error(`"body"" field should be type []byte rather than other types`)
	}

	if string(content) != fakeContent {
		t.Error(`string(content) is not equal original body`)
	}

}

func TestNewAnalyzerWithCallback(t *testing.T) {
	fakeURL := "https://www.google.com"
	fakeContent := "test content"
	// make a fake response
	res := FakeResponse(fakeURL, fakeContent)
	res.Request.SetCallback("urlparser")

	// NewAnalyzer will init a Analyzer will ContentReader as parser
	analyzer, err := NewAnalyzer(ParserMap{
		"urlparser": func(res *Response) ([]Data, error) {
			return Item{"url": res.Request.URL.String()}.DataList(), nil
		},
	})

	if err != nil {
		t.Error(err)
	}

	datums, err := analyzer.Analyze(res)
	if err != nil {
		t.Error(err)
	}

	if len(datums) < 1 {
		t.Error("Analyzer should generate at least one Item in datums contains http body")
	}

	item, ok := datums[0].(Item)
	if !ok {
		t.Error("Analyzer should generate a item rather than other Data entry")
	}

	urlStr, ok := item["url"]
	if !ok {
		t.Error(`Analyzer's yield item should have field "url"`)
	}

	url, ok := urlStr.(string);
	if !ok {
		t.Error(`default Analyzer's content field should be type string rather than other types`)
	}

	if string(url) != fakeURL {
		t.Error(`default Analyzer's string(content) is not equal original body`)
	}
}
