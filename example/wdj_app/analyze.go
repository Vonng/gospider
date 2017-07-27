package wdj_app

import (
	"strings"
	"bytes"
	"strconv"
	"time"
)

import "github.com/PuerkitoBio/goquery"
import . "github.com/Vonng/gospider"

func ParseWdjApp(res *Response) ([]Data, error) {
	item := make(Item, 2)
	apk := ApkFromPageURL(res.Request.URL.String())
	if apk == "" {
		return nil, ErrParse
	}
	item["apk"] = apk

	app := NewWdjApp(apk)
	doc, err := goquery.NewDocumentFromResponse(res.Response);
	if err != nil {
		return nil, err
	}

	if err = app.ParseFrom(doc); err != nil {
		return nil, err
	}
	item["data"] = app

	return item.DataList(), nil
}

const chineseTimeFormat = "2006年01月02日"

// WdjApp_ParseFrom will parse wandoujia App from goquery document.
func (app *WdjApp) ParseFrom(doc *goquery.Document) error {
	// app.Apk: 包名，必需存在
	if apk, exist := doc.Find("body").Attr("data-pn"); exist && apk != "" {
		app.Apk = apk
	} else {
		return ErrParse
	}

	// app.Name: 应用名称，从标题解析
	if appName := strings.TrimSpace(doc.Find("p.app-name span.title").Text()); appName != "" {
		app.Name.String = appName
		app.Name.Valid = true
	}

	// app.PageURL: 应用下载页面URL
	app.PageURL.String = PageURL(app.Apk)
	app.PageURL.Valid = true

	if icon, exist := doc.Find("div.app-icon img").Attr("src"); exist {
		app.IconURL.String = icon
		app.IconURL.Valid = true
	}

	// app.IconURL: 应用图标URL
	if icon, exist := doc.Find("div.app-icon img").Attr("src"); exist {
		app.IconURL.String = icon
		app.IconURL.Valid = true
	}

	// app.DownURL: 应用APK下载链接
	if url, exist := doc.Find("a.install-btn").Attr("href"); exist {
		app.DownURL.String = url
		app.DownURL.Valid = true
	}

	infosList := doc.Find("dl.infos-list") // 右侧包含了分类、标签的信息栏
	numList := doc.Find("div.num-list")    // 右侧上方数值列表

	// app.Size: 应用安装包大小
	if sizeStr, exist := infosList.Find("meta[itemprop=fileSize]").Attr("content"); exist {
		if sz, err := PrefixedBytesToInt(sizeStr); err == nil {
			app.Size.Int64 = sz
			app.Size.Valid = true
		}
	}

	// app.InstallCnt: 应用安装数目
	if icStr := strings.TrimSpace(numList.Find("i[itemprop=interactionCount]").Text()); icStr != "" {
		if cnt, err := ChineseSuffixStringToInt(icStr); err == nil {
			app.InstallCnt.Int64 = cnt
			app.InstallCnt.Valid = true
		}
	}

	// app.CommentCnt: 应用评论数
	if ccStr := strings.TrimSpace(numList.Find("a.comment-open i").Text()); ccStr != "" {
		if cc, err := strconv.Atoi(ccStr); err == nil {
			app.CommentCnt.Int64 = int64(cc)
			app.CommentCnt.Valid = true
		}
	}

	// app.FavorRate: 应用好评率
	if frStr := strings.TrimSpace(numList.Find("span.love i").Text()); frStr != "" && frStr != "暂无" {
		if strings.HasSuffix(frStr, "%") {
			if f, err := strconv.ParseFloat(strings.Trim(frStr, "%"), 32); err == nil {
				app.FavorRate.Int64 = int64(f)
				app.FavorRate.Valid = true
			}
		}
	}

	// app.LastMtime: 应用上次更新时间
	if lmtStr, exist := infosList.Find("#baidu_time").Attr("datetime"); exist {
		if t, err := time.Parse(chineseTimeFormat, lmtStr); err == nil {
			app.LastMtime = t
		}
	}

	// app.LastChange: 应用上次更新内容
	if lcStr, err := doc.Find("div.change-info div").Html(); err == nil && lcStr != "" {
		app.LastChange.String = changeBrToNewLine(lcStr)
		app.LastChange.Valid = true
	}

	// app.LastVersion: 最新版本
	if lvStr := strings.TrimSpace(infosList.Find("dd:nth-last-of-type(3)").Text()); lvStr != "" {
		app.LastVersion.String = lvStr
		app.LastVersion.Valid = true
	}

	// app.Vendor: App厂商
	if vendor := strings.TrimSpace(infosList.Find("span.dev-sites").Text()); vendor != "" {
		app.Vendor.String = vendor
		app.Vendor.Valid = true
	}

	// app.Subtitle: 副标题，又称为Tagline
	if subtitle := strings.TrimSpace(doc.Find("p.tagline").Text()); subtitle != "" {
		app.Subtitle.String = subtitle
		app.Subtitle.Valid = true
	}

	// app.Review: 编辑评论
	if review := strings.TrimSpace(doc.Find("div.editorComment div.con").Text()); review != "" {
		app.Review.String = review
		app.Review.Valid = true
	}

	// app.Description: 应用的描述信息
	if desc, err := doc.Find("div.desc-info div.con").Html(); err == nil && desc != "" {
		app.Description.String = changeBrToNewLine(desc)
		app.Description.Valid = true
	}

	// app.System: 应用的系统要求
	if requireNodes := infosList.Find("dd.perms").Nodes; len(requireNodes) > 0 {
		if len(requireNodes) > 0 && requireNodes[0].FirstChild != nil {
			if sysStr := strings.TrimSpace(requireNodes[0].FirstChild.Data); strings.HasPrefix(sysStr, "Android ") {
				app.System.String = strings.Trim(strings.TrimLeft(sysStr, "Android "), " 以上")
				app.System.Valid = true
			}
		}
	}

	// app.Permissions: 应用要求的权限
	app.Permissions = removeEmpty(doc.Find("span.perms").Map(func(ind int, s *goquery.Selection) string {
		return strings.TrimSpace(s.Text())
	}))

	// app.Tags: 应用标签
	app.Categories = removeEmpty(infosList.Find("dd.tag-box a").Map(func(ind int, s *goquery.Selection) string {
		return strings.TrimSpace(s.Text())
	}))

	// app.Categories: 应用分类(多个)
	app.Tags = removeEmpty(infosList.Find("div.tag-box a").Map(func(ind int, s *goquery.Selection) string {
		return strings.TrimSpace(s.Text())
	}))

	// app.Crumb: 面包屑
	app.Crumb = removeEmpty(doc.Find("div.crumb a").Map(func(ind int, s *goquery.Selection) string {
		return strings.TrimSpace(s.Text())
	}))
	// 丢弃面包屑第一项：“应用首页”
	if len(app.Crumb) > 1 {
		app.Crumb = app.Crumb[1:]
	}

	// app.RelateApps: 相关应用,以Apk标识
	app.RelateApps = removeEmpty(doc.Find("ul.relative-download li a.d-btn").Map(func(ind int, s *goquery.Selection) string {
		pname, _ := s.Attr("data-app-pname")
		return strings.TrimSpace(pname)
	}))

	// app.Screenshots: 应用截图
	app.Screenshots = removeEmpty(doc.Find("img.screenshot-img").Map(func(ind int, s *goquery.Selection) string {
		imgSrc, _ := s.Attr("src")
		return imgSrc
	}))

	// app.Comments: 评论
	app.Comments = removeEmpty(doc.Find("ul.comments-list li.normal-li").Map(func(ind int, s *goquery.Selection) string {
		user := strings.TrimSpace(s.Find("p.first span.name").Text())
		ts := squeezeTime(strings.TrimSpace(s.Find("p.first span:last-of-type").Text()))
		content := strings.TrimSpace(s.Find("p.cmt-content span").Text())
		if user == "" || ts == "" || len(ts) != 8 {
			return ""
		}
		var b bytes.Buffer
		b.WriteString(ts)
		b.WriteByte(',')
		b.WriteString(user)
		b.WriteByte(',')
		b.WriteString(content)
		return b.String()
	}))
	return nil
}
