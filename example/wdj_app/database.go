package wdj_app

import (
	"github.com/go-redis/redis"
	"github.com/go-pg/pg"
	"time"
)

/**************************************************************
* REDIS
**************************************************************/

// apk from forceTodoList won't be filtered
var (
	Redis    *redis.Client
	redisURL = map[string]string{
		"dev":  "redis://localhost:6379/0",
		"prod": "redis://:myredis@localhost:6379/0",
	}
	pollTimeout       = time.Minute
	redisFilterKey    = "wdj:app:seen"
	redisTodoKey      = "wdj:app:todo"
	redisForceTodoKey = "wdj:app:todo:force"
)

// InitRedis will init global redis instance will given url
func InitRedis(redisURL string) error {
	// parse redis url
	redisOption, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}

	// test redis is available
	Redis = redis.NewClient(redisOption)
	_, err = Redis.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

/**************************************************************
* POSTGRESQL
**************************************************************/

// Postgres Instance & Prepared statements
var (
	Pg    *pg.DB
	pgURL = map[string]string{
		"dev":  "postgres://vonng@localhost:5432/app?sslmode=disable",
		"prod": "postgres://app:appapp@localhost:5432/app?sslmode=disable",
	}
	setStatusStmt *pg.Stmt
	upsertStmt    *pg.Stmt
	setStatusSQL  = `INSERT INTO android(apk,wdj) VALUES ($1,$2)
ON CONFLICT(apk) DO UPDATE SET wdj = EXCLUDED.wdj;`

	upsertSQL = `INSERT INTO wdj_app (apk, name, page_url, icon_url, down_url, size, install_cnt, comment_cnt,
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
)

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
