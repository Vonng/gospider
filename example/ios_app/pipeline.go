package ios_app

import . "github.com/Vonng/gospider"
import (
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

const dataDir = "/var/data/ios/"

func GetiOSPipeline() Pipeline {
	return NewPipelineSolo(WriteItemToFile)
}

func WriteItemToFile(item Item) error {
	content := item.GetString(KeyBody)
	id := item.GetString("id")

	if content == "" || id == "" {
		log.Errorf("processing file failed : %v", item)
		return ErrNilItem
	}
	filename := dataDir + id + ".html"
	log.Infof("[WRITE] %s to file %s", id, filename)
	return ioutil.WriteFile(filename, []byte(content), 0666)
}
