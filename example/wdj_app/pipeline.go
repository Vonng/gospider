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
// GetPipeline will build pipeline from pgURL
func GetPipeline(pgURL string) (Pipeline, error) {
	if err := InitPg(pgURL); err != nil {
		return nil, err
	}
	return NewPipelineSolo(Save), nil
}

// Save is the only processor that pipe use
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
