package main

import . "github.com/Vonng/gospider/example/wdj"

func main() {
	Run("dev")
}

//var redisURL = map[string]string{
//	"dev":  "redis://localhost:6379/0",
//	"prod": "redis://:myredis@localhost:6379/0",
//}
//var pgURL = map[string]string{
//	"dev":  "postgres://vonng@localhost:5432/app?sslmode=disable",
//	"prod": "postgres://app:appapp@localhost:5432/app?sslmode=disable",
//}
//
//const redisFilterKey = "wdj:app:seen"
//
//
//func main() {
//	Env := "dev"
//	if os.Getenv("ENV") == "prod" {
//		Env = "prod"
//	}
//
//	log.Info("[INIT] WdjApp Spider Start!")
//	engine := BuildEngine(redisURL[Env], pgURL[Env])
//	if engine == nil {
//		log.Fatal("[INIT] build engine failed")
//	}
//
//	log.Info("[INIT] WdjApp Engine! begin running")
//	for err := range engine.Run() {
//		log.Errorf("[ERROR] %s", err.Error())
//	}
//}
