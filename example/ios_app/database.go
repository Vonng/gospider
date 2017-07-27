package ios_app

import (
	"github.com/go-redis/redis"
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
	redisFilterKey    = "ios:app:seen"
	redisTodoKey      = "ios:app:todo"
	redisForceTodoKey = "ios:app:todo:force"
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
