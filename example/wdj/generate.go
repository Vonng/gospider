package wdj_app

import "time"
import . "github.com/Vonng/gospider"
import "github.com/go-redis/redis"

const pollTimeout = time.Minute
const todoListKey = "wdj:app:todo"
const forceTodoListKey = "wdj:app:todo:force"

// apk from forceTodoList won't be filtered

// GetGenerator will pull apk name from redis
func RequestGenerator(redisURL string) (<-chan Data, error) {
	client, err := NewRedisInstance(redisURL)
	if err != nil {
		return nil, err
	}

	c := make(chan Data)
	go func(c chan<- Data) {
		for {
			res, err := client.BRPop(pollTimeout, todoListKey, forceTodoListKey).Result()
			if err != nil {
				if err.Error() == "redis: nil" {
					continue
				} else {
					close(c)
				}
			}
			if len(res) == 2 {
				key := res[0]
				apk := res[1]
				var req *Request
				if key == forceTodoListKey {
					req, err = NewRequest(
						"GET",
						PageURL(apk),
						nil,
						MetaMap{"filter": false},
					)
					if err == nil {
						c <- req.DisableFilter()
					}
				} else if key == todoListKey {
					req, err = NewRequest(
						"GET",
						PageURL(apk),
						nil,
						nil,
					)
					if err == nil {
						c <- req
					}
				}
			}
		}
	}(c)

	return (<-chan Data)(c), nil
}

func NewRedisInstance(redisURL string) (*redis.Client, error) {
	// parse redis url
	redisOption, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	// test redis is available
	client := redis.NewClient(redisOption)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
