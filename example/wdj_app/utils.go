package wdj_app

import (
	"strings"
	"strconv"
	"unicode"
	"github.com/PuerkitoBio/goquery"
	"errors"
)

/**************************************************************\
* 辅助函数
***************************************************************/

// ChineseSuffixStringToInt 将形如 "1.28亿"转换为相应的整型值
func ChineseSuffixStringToInt(s string) (res int64, err error) {
	r := []rune(s)
	n := len(r)

	var mutiplier float64;
	switch r[n-1] {
	case rune('万'):
		mutiplier = 10000
		r = r[0:n-1]
	case rune('亿'):
		mutiplier = 100000000
		r = r[0:n-1]
	default:
		mutiplier = 1
	}

	numStr := string(r)
	if dotInd := strings.Index(numStr, "."); dotInd == -1 {
		// 没有小数点
		if i, err := strconv.Atoi(numStr); err != nil {
			return 0, err
		} else {
			return int64(float64(i) * mutiplier), nil
		}
	} else {
		// 有小数点,判断小数位数并移除小数点
		for i := 0; i < len(numStr)-dotInd-1; i++ {
			mutiplier /= 10
		}

		numStr = strings.Replace(numStr, ".", "", 1)
		if i, err := strconv.Atoi(numStr); err != nil {
			return 0, err
		} else {
			return int64(float64(i) * mutiplier), nil
		}
	}
}

// PrefixedBytesToInt 用于将形如"128k" 转换为相应的字节数
func PrefixedBytesToInt(s string) (res int64, err error) {
	var i, nFrac int
	var val int64
	var c byte
	var dot bool

	// parse numeric val (omit dot), and length of frac part
Loop:
	for i < len(s) {
		c = s[i]
		switch {
		case '0' <= c && c <= '9':
			val *= 10
			val += int64(c - '0')
			if dot {
				nFrac ++
			}
			i++
		case c == '.':
			dot = true
			i++
		default:
			break Loop
		}
	}
	unit := strings.ToUpper(strings.TrimSpace(s[i:]))

	switch unit {
	case "", "B":
	case "KB", "K":
		val <<= 10
	case "MB", "M":
		val <<= 20
	case "GB", "G":
		val <<= 30
	case "TB", "T":
		val <<= 40
	case "PB", "P":
		val <<= 50
	case "EB", "E":
		val <<= 60
	default:
		return 0, errors.New("parse multiplier failed")
	}

	// handle frac
	for j := 0; j < nFrac; j++ {
		val /= 10
	}

	return val, nil
}

// removeEmpty 辅助函数:将字符串数组中的空字符串移除
func removeEmpty(input []string) (output []string) {
	for _, str := range input {
		if str != "" {
			output = append(output, str)
		}
	}
	return
}

// squeezeTime 辅助函数:将"2015年05月10日"压缩为"20150510"的形式
func squeezeTime(s string) string {
	var nb []rune
	for _, ch := range s {
		if unicode.IsDigit(ch) {
			nb = append(nb, ch)
		}
	}
	return string(nb)
}

// changeBrToNewLine 将文本中的换行转换为\n
func changeBrToNewLine(s string) string {
	s = strings.Replace(s, "<br>", "\n", -1)
	s = strings.Replace(s, "<br/>", "\n", -1)
	return strings.TrimSpace(s)
}

// ParseWdjAppFromApk 会直接根据包名生成URL并爬取
func ParseWdjAppFromApk(apk string) (app *WdjApp, err error) {
	app = NewWdjApp(apk)
	doc, err := goquery.NewDocument(PageURL(apk))
	if err != nil {
		return
	}
	err = app.ParseFrom(doc)
	return
}
