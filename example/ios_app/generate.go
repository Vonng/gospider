package ios_app

import . "github.com/Vonng/gospider"

// GetGenerator will pull appid name from redis
func RequestGenerator(redisURL string) (<-chan Data, error) {
	InitRedis(redisURL)
	c := make(chan Data)
	go func(c chan<- Data) {
		for {
			res, err := Redis.BRPop(pollTimeout, redisTodoKey, redisForceTodoKey).Result()
			if err != nil {
				if err.Error() == "redis: nil" {
					continue
				} else {
					close(c)
				}
			}
			if len(res) == 2 {
				key := res[0]
				appid := res[1]
				var req *Request
				if key == redisForceTodoKey {
					req, err = NewRequest(
						"GET",
						PageURL(appid),
						nil,
						MetaMap{"id": appid },
					)
					if err == nil {
						c <- req.DisableFilter()
					}
				} else if key == redisTodoKey {
					req, err = NewRequest(
						"GET",
						PageURL(appid),
						nil,
						MetaMap{"id": appid},
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
