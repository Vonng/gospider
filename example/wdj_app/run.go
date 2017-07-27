package wdj_app

import . "github.com/Vonng/gospider"
import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func BuildEngine(redisURL, pgURL string) Engine {
	analyzer, err := NewAnalyzerSolo(ParseWdjApp)
	if err != nil {
		log.Errorf("build wdj app analyzer failed!", err.Error())
		return nil
	}

	downloader, err := NewDownloader(nil)
	if err != nil {
		log.Errorf("build wdj app downloader failed! %s", err.Error())
		return nil
	}

	pipeline, err := GetPipeline(pgURL)
	if err != nil {
		log.Errorf("build wdj app pipeline failed! %s", err.Error())
		return nil
	}

	filter, err := NewRedisBloomFilter(redisURL, redisFilterKey)
	if err != nil {
		log.Errorf("build wdj app redis bloom filter failed! %s", err.Error())
		return nil
	}

	args := EngineArgs{
		Filter:      filter,
		Downloader:  downloader,
		Analyzer:    analyzer,
		Pipeline:    pipeline,
		DWorkers:    10,
		ReqBufSize:  0,
		ResBufSize:  10000,
		ItemBufSize: 10000,
		ErrBufSize:  10000,
	}

	return NewEngine(&args)
}

func Run() {
	Env := "dev"
	if os.Getenv("ENV") == "prod" {
		Env = "prod"
	}

	engine := BuildEngine(redisURL[Env], pgURL[Env])
	if engine == nil {
		log.Error("Build engine failed")
	}

	generator, err := RequestGenerator(redisURL[Env])
	if err != nil {
		log.Errorf("build wdj app request generator failed! %s", err.Error())
	}

	for err := range engine.Run(generator) {
		log.Error(err)
	}
}
