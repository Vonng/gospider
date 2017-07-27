package wdj_app

import log "github.com/Sirupsen/logrus"

var redisURL = map[string]string{
	"dev":  "redis://localhost:6379/0",
	"prod": "redis://:myredis@localhost:6379/0",
}
var pgURL = map[string]string{
	"dev":  "postgres://vonng@localhost:5432/app?sslmode=disable",
	"prod": "postgres://app:appapp@localhost:5432/app?sslmode=disable",
}

const redisFilterKey = "wdj:app:seen"

func Run(Env string) error {
	engine := BuildEngine(redisURL[Env], pgURL[Env], redisFilterKey)
	generator, err := RequestGenerator(redisURL[Env])
	if err != nil {
		log.Error("build wdj app request generator failed!")
		return nil
	}

	for err := range engine.Run(generator) {
		log.Error(err)
	}
	return nil
}
