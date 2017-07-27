package wdj_app

import . "github.com/Vonng/gospider"
import log "github.com/Sirupsen/logrus"

func BuildEngine(redisURL, pgURL, dedupeKey string) Engine {
	analyzer, err := NewAnalyzerFromParser(ParseWdjApp)
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

	filter, err := NewRedisBloomFilter(redisURL, dedupeKey)
	if err != nil {
		log.Errorf("build wdj app redis bloom filter failed! %s", err.Error())
		return nil
	}
	args := EngineArgs{
		Filter:      filter,
		Downloader:  downloader,
		Analyzer:    analyzer,
		Pipeline:    pipeline,
		DWorkers:    20,
		ReqBufSize:  10,
		ResBufSize:  10000,
		ItemBufSize: 10000,
		ErrBufSize:  10000,
	}
	return NewEngine(&args)
}
