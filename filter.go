package gospider

import (
	"sync"
	"github.com/go-redis/redis"
)

/**************************************************************
* interface: Filter
**************************************************************/

// Filter : interface for eliminate duplicate request。
type Filter interface {
	// Seen indicate whether s is in seen list and update it
	Seen(req *Request) (bool)
}

/**************************************************************
* struct: defaultFilter
**************************************************************/

// defaultFilter implemented with sync.Map
type defaultFilter struct {
	seen sync.Map
}

// NewMapFilter create a default dupe filter
func NewMapFilter() Filter {
	return &defaultFilter{}
}

// Seen: Caller must guarantee req is not nil
func (self *defaultFilter) Seen(req *Request) (bool) {
	_, seen := self.seen.Load(req.URL)
	self.seen.Store(req.URL, nil)
	return seen
}

// defaultFilter_SeenURL using url directly instead of request
func (self *defaultFilter) SeenURL(url string) bool {
	_, seen := self.seen.Load(url)
	self.seen.Store(url, nil)
	return seen
}

/**************************************************************
* struct: redisFilter
**************************************************************/

// redisFilter using HyperLogLog to determine whether seen a url
// It may loss some accuracy trade for memory usage
// little FP rate: some unseen page could be consider seen
type redisFilter struct {
	key    string
	client *redis.Client
}

// redisFilter create a new redis dupe filter using PFADD
func NewRedisFilter(redisURL string, key string) (Filter, error) {
	ops, err := redis.ParseURL(redisURL);
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(ops)

	if _, err = client.Ping().Result(); err != nil {
		return nil, err
	}

	return &redisFilter{key, client}, nil
}

// redisFilter_SeenURL using url directly instead of request
func (self *redisFilter) Seen(req *Request) bool {
	i, _ := self.client.PFAdd(self.key, req.URL).Result()
	return i == 0
}

// redisFilter_Seen implemented with PFAdd
func (self *redisFilter) SeenURL(url string) bool {
	i, _ := self.client.PFAdd(self.key, url).Result()
	return i == 0
}
