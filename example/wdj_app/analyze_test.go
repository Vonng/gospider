package wdj_app

import (
	"fmt"
	"testing"
)

func TestParseWdjAppFromPkgName(t *testing.T) {
	app, err := ParseWdjAppFromApk("com.autonavi.minimap")
	if err != nil {
		t.Error(err)
	}
	app.Print()
}

func TestChineseSuffixStringToInt(t *testing.T) {
	testCase := []struct {
		Input  string
		Expect int64
	}{
		{"2.56万", 25600},
		{"256万", 2560000},
		{"256.万", 2560000},
		{"25.6万", 256000},
		{"0.256万", 2560},
		{"1.256万", 12560},
		{"1.2567万", 12567},
		{"1.25678万", 12567},
		{"1.25678亿", 125678000},
		{"1.25678901亿", 125678901},
		{"1.256789019亿", 125678901},
		{"0.00001亿", 1000},
	}

	for _, c := range testCase {
		output, err := ChineseSuffixStringToInt(c.Input)
		if err != nil {
			t.Error(err)
		}

		if output != c.Expect {
			t.Errorf("Input[%s] Expect[%d] Got[%d]\n", c.Input, c.Expect, output)
		}
		fmt.Printf("Input[%s] Expect[%d] = Got[%d]\n", c.Input, c.Expect, output)
	}
}

func TestPrefixedBytesToInt(t *testing.T) {
	testCase := []struct {
		Input  string
		Expect int64
	}{
		{"256", 256},
		{"256.128", 256},
		{"2KB", 2048},
		{"2.56KB", 2621},
		{"1024K", 1048576},
		{"2M", 2097152},
		{"2.5M", 2621440},
		{"2.5432M", 2666738},
	}

	for _, c := range testCase {
		output, err := PrefixedBytesToInt(c.Input)
		if err != nil {
			t.Error(err)
		}

		if output != c.Expect {
			t.Errorf("Input[%s] Expect[%d] Got[%d]\n", c.Input, c.Expect, output)
		}
		fmt.Printf("Input[%s] Expect[%d] = Got[%d]\n", c.Input, c.Expect, output)
	}
}
