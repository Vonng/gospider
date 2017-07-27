package wdj_app

import (
	"fmt"
	"time"
	"bytes"
	"text/template"
	"database/sql"
	"strings"
)

import . "github.com/Vonng/gospider"

// WdjApp 包括了豌豆荚应用页上的有用信息，也对应着数据库中的相应表结构
type WdjApp struct {
	Apk         string `sql:",pk"`
	Name        sql.NullString
	PageURL     sql.NullString
	IconURL     sql.NullString
	DownURL     sql.NullString
	Size        sql.NullInt64
	InstallCnt  sql.NullInt64
	CommentCnt  sql.NullInt64
	FavorRate   sql.NullInt64
	LastMtime   time.Time
	LastChange  sql.NullString
	LastVersion sql.NullString
	Vendor      sql.NullString
	Subtitle    sql.NullString
	Review      sql.NullString
	Description sql.NullString
	System      sql.NullString
	Permissions []string `pg:",array"`
	Tags        []string `pg:",array"`
	Categories  []string `pg:",array"`
	Crumb       []string `pg:",array"`
	RelateApps  []string `pg:",array"`
	Screenshots []string `pg:",array"`
	Comments    []string `pg:",array"`
	CrawlTime   time.Time
	TableName   struct{} `json:"-" xml:"-" sql:"wdj_app"`
}

// WdjAppURLPrefix 豌豆荚的应用页存在模式：http://www.wandoujia.com/apps/<package_name>
const WdjAppPageURLPrefix = "http://www.wandoujia.com/apps/"
const WdjAppDownURLPattern = "http://www.wandoujia.com/apps/%s/binding"

func NewWdjApp(apk string) *WdjApp {
	return &WdjApp{Apk: apk    }
}

func ApkFromPageURL(url string) string {
	return strings.TrimLeft(url, WdjAppPageURLPrefix)
}

func PageURL(apk string) string {
	return WdjAppPageURLPrefix + apk
}

func DownURL(apk string) string {
	return fmt.Sprintf(WdjAppDownURLPattern, apk)
}

var wdjAppTmpl, _ = template.New("WdjAppEntity").Parse(`
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
┃ 豌豆荚应用: {{ .Apk }}
┃ 页面链接：  {{.PageURL}}
┣┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈
┃ Apk         ┆ {{ .Apk }}
┃ Name        ┆ {{if .Name.Valid        }}{{ .Name.String        }} {{else}}NULL{{end}}
┃ PageURL     ┆ {{if .PageURL.Valid     }}{{ .PageURL.String     }} {{else}}NULL{{end}}
┃ IconURL     ┆ {{if .IconURL.Valid     }}{{ .IconURL.String     }} {{else}}NULL{{end}}
┃ DownURL     ┆ {{if .DownURL.Valid     }}{{ .DownURL.String     }} {{else}}NULL{{end}}
┃ Size        ┆ {{if .Size.Valid        }}{{ .Size.Int64         }} {{else}}NULL{{end}}
┃ InstallCnt  ┆ {{if .InstallCnt.Valid  }}{{ .InstallCnt.Int64   }} {{else}}NULL{{end}}
┃ CommentCnt  ┆ {{if .CommentCnt.Valid  }}{{ .CommentCnt.Int64   }} {{else}}NULL{{end}}
┃ FavorRate   ┆ {{if .FavorRate.Valid   }}{{ .FavorRate.Int64    }} {{else}}NULL{{end}}
┃ LastMtime   ┆ {{ .LastMtime }}
┃ LastChange  ┆ {{if .LastChange.Valid  }}{{ .LastChange.String  }} {{else}}NULL{{end}}
┃ LastVersion ┆ {{if .LastVersion.Valid }}{{ .LastVersion.String }} {{else}}NULL{{end}}
┃ Vendor      ┆ {{if .Vendor.Valid      }}{{ .Vendor.String      }} {{else}}NULL{{end}}
┃ Subtitle    ┆ {{if .Subtitle.Valid    }}{{ .Subtitle.String    }} {{else}}NULL{{end}}
┃ Review      ┆ {{if .Review.Valid      }}{{ .Review.String      }} {{else}}NULL{{end}}
┃ Description ┆ {{if .Description.Valid }}{{ .Description.String }} {{else}}NULL{{end}}
┃ System      ┆ {{if .System.Valid      }}{{ .System.String      }} {{else}}NULL{{end}}
┃ Permissions ┆ {{ .Permissions  }}
┃ Tags        ┆ {{ .Tags         }}
┃ Categories  ┆ {{ .Categories   }}
┃ Crumb       ┆ {{ .Crumb        }}
┃ RelateApps  ┆ {{ .RelateApps   }}
┃ Screenshots ┆ {{ .Screenshots  }}
┃ Comments    ┆ {{ .Comments     }}
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`)

func (app *WdjApp) Valid() bool {
	return app != nil && app.Apk != "" && app.Name.Valid
}

// Print 打印出人类可读版本的应用信息
func (app *WdjApp) Print() {
	buf := new(bytes.Buffer)
	if err := wdjAppTmpl.Execute(buf, app); err != nil {
		panic(err)
		return
	}
	fmt.Println(buf.String())
}

func (app *WdjApp) Item() Item {
	m := make(Item, 2)
	m["data"] = app
	m["apk"] = app.Apk
	return m
}
