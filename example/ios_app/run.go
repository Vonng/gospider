package ios_app

import . "github.com/Vonng/gospider"
import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func BuildEngine() Engine {
	analyzer, err := NewAnalyzerSolo(BodyReader)
	if err != nil {
		log.Errorf("build naive ios analyzer failed!", err.Error())
		return nil
	}

	downloader, err := NewDownloader(nil)
	if err != nil {
		log.Errorf("build naive ios downloader failed! %s", err.Error())
		return nil
	}

	pipeline := GetiOSPipeline()
	if err != nil {
		log.Errorf("build naive ios app pipeline failed! %s", err.Error())
		return nil
	}

	args := EngineArgs{
		Filter:      nil,
		Downloader:  downloader,
		Analyzer:    analyzer,
		Pipeline:    pipeline,
		DWorkers:    20,
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

	engine := BuildEngine()
	if engine == nil {
		log.Error("Build ios app engine failed")
	}

	generator, err := RequestGenerator(redisURL[Env])
	if err != nil {
		log.Errorf("build ios app request generator failed! %s", err.Error())
	}

	for err := range engine.Run(generator) {
		log.Error(err)
	}
}
