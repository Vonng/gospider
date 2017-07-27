package wdj_app

import (
	"time"
	"errors"
)

import (
	. "github.com/Vonng/gospider"
	"github.com/go-pg/pg"
	log "github.com/Sirupsen/logrus"
)

/**************************************************************
* Pipeline integration
**************************************************************/
func GetPipeline(pgURL string) (Pipeline, error) {
	if err := InitPg(pgURL); err != nil {
		return nil, err
	}
	return NewPipelineSolo(Save), nil
}

func Save(item Item) (err error) {
	var app *WdjApp
	if wdjApp, ok := item["data"]; ok {
		if v, ok := wdjApp.(*WdjApp); ok {
			app = v
		}
	}

	if app == nil {
		return ErrNilItem
	}

	if app.Valid() {
		log.Infof("[PIPE] item done %s status=done", app.Apk, )
		if err = app.Save(); err != nil {
			return
		}
		if err = app.SaveStatus(StatusDone); err != nil {
			return
		}
	} else {
		log.Infof("[PIPE] item done %s status=fail", app.Apk, )
		return app.SaveStatus(StatusFail)
	}
	return nil
}

/**************************************************************
* Postgresql related works
**************************************************************/
// Postgres Instance & Prepared statements
var (
	Pg            *pg.DB
	setStatusStmt *pg.Stmt
	upsertStmt    *pg.Stmt
)

const setStatusSQL = `INSERT INTO android(apk,wdj) VALUES ($1,$2)
ON CONFLICT(apk) DO UPDATE SET wdj = EXCLUDED.wdj;`

const upsertSQL = `INSERT INTO wdj_app (apk, name, page_url, icon_url, down_url, size, install_cnt, comment_cnt,
                     favor_rate, last_mtime, last_change, last_version, vendor, subtitle, review,
                     description, system, permissions, tags, categories, crumb, relate_apps, screenshots,
                     comments, crawl_time) VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
ON CONFLICT (apk)
  DO UPDATE SET
    name         = EXCLUDED.name,
    page_url     = EXCLUDED.page_url,
    icon_url     = EXCLUDED.icon_url,
    down_url     = EXCLUDED.down_url,
    size         = EXCLUDED.size,
    install_cnt  = EXCLUDED.install_cnt,
    comment_cnt  = EXCLUDED.comment_cnt,
    favor_rate   = EXCLUDED.favor_rate,
    last_mtime   = EXCLUDED.last_mtime,
    last_change  = EXCLUDED.last_change,
    last_version = EXCLUDED.last_version,
    vendor       = EXCLUDED.vendor,
    subtitle     = EXCLUDED.subtitle,
    review       = EXCLUDED.review,
    description  = EXCLUDED.description,
    system       = EXCLUDED.system,
    permissions  = EXCLUDED.permissions,
    tags         = EXCLUDED.tags,
    categories   = EXCLUDED.categories,
    crumb        = EXCLUDED.crumb,
    relate_apps  = EXCLUDED.relate_apps,
    screenshots  = EXCLUDED.screenshots,
    comments     = EXCLUDED.comments,
    crawl_time   = EXCLUDED.crawl_time;`

func InitPg(pgURL string) error {
	pgOption, err := pg.ParseURL(pgURL)
	if err != nil {
		return err
	}

	Pg = pg.Connect(pgOption)
	setStatusStmt, err = Pg.Prepare(setStatusSQL)
	if err != nil {
		return err
	}

	upsertStmt, err = Pg.Prepare(upsertSQL)
	if err != nil {
		return err
	}
	return nil
}

// WdjApp_Save will save app to corresponding postgres tables `wdj_app`
func (app *WdjApp) Save() (err error) {
	app.CrawlTime = time.Now()
	_, err = upsertStmt.Exec(
		app.Apk,
		app.Name,
		app.PageURL,
		app.IconURL,
		app.DownURL,
		app.Size,
		app.InstallCnt,
		app.CommentCnt,
		app.FavorRate,
		app.LastMtime,
		app.LastChange,
		app.LastVersion,
		app.Vendor,
		app.Subtitle,
		app.Review,
		app.Description,
		app.System,
		pg.Array(app.Permissions),
		pg.Array(app.Tags),
		pg.Array(app.Categories),
		pg.Array(app.Crumb),
		pg.Array(app.RelateApps),
		pg.Array(app.Screenshots),
		pg.Array(app.Comments),
		app.CrawlTime,
	)
	return err
}

// StatusXXXX 定义了页面的状态
const (
	StatusTodo  = 0
	StatusDoing = 1
	StatusFail  = 2
	StatusDone  = 3
)

// SaveStatus will change the database table `android`
func (app *WdjApp) SaveStatus(status int16) (err error) {
	if status < StatusTodo || status > StatusDone {
		return errors.New("invalid status code")
	}
	_, err = setStatusStmt.Exec(app.Apk, status)
	return
}
