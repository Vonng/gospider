package gospider

import (
	"net/url"
	"sync"
	"github.com/go-redis/redis"
)

import . "github.com/PuerkitoBio/purell"

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
* struct: redisSetFilter
**************************************************************/

// redisSetFilter using HyperLogLog to determine whether seen a url
// It may loss some accuracy trade for memory usage
// little FP rate: some unseen page could be consider seen
type redisSetFilter struct {
	key    string
	client *redis.Client
}

// redisSetFilter create a new redis dupe filter using PFADD
func NewRedisSetFilter(redisURL string, key string) (Filter, error) {
	ops, err := redis.ParseURL(redisURL);
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(ops)

	if _, err = client.Ping().Result(); err != nil {
		return nil, err
	}

	return &redisSetFilter{key, client}, nil
}

// redisSetFilter_SeenURL using url directly instead of request
func (self *redisSetFilter) Seen(req *Request) bool {
	i, _ := self.client.SAdd(self.key, req.URL).Result()
	return i == 0
}

// redisSetFilter_Seen implemented with PFAdd
func (self *redisSetFilter) SeenURL(url string) bool {
	i, _ := self.client.PFAdd(self.key, url).Result()
	return i == 0
}

/**************************************************************
* struct: redisBloomFilter
**************************************************************/

// redisBloomFilter using HyperLogLog to determine whether seen a url
// It may loss some accuracy trade for memory usage
// little FP rate: some unseen page could be consider seen
type redisBloomFilter struct {
	key    string
	client *redis.Client
}

// redisBloomFilter create a new redis dupe filter using PFADD
func NewRedisBloomFilter(redisURL string, key string) (Filter, error) {
	ops, err := redis.ParseURL(redisURL);
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(ops)

	if _, err = client.Ping().Result(); err != nil {
		return nil, err
	}

	return &redisBloomFilter{key, client}, nil
}

// redisBloomFilter_SeenURL using url directly instead of request
func (self *redisBloomFilter) Seen(req *Request) bool {
	i, _ := self.client.PFAdd(self.key, req.URL).Result()
	return i == 0
}

// redisBloomFilter_Seen implemented with PFAdd
func (self *redisBloomFilter) SeenURL(url string) bool {
	i, _ := self.client.PFAdd(self.key, url).Result()
	return i == 0
}

/**************************************************************
* function: PureURL
**************************************************************/
var URLNormalizeFlag = FlagsSafe |
	FlagDecodeUnnecessaryEscapes |
	FlagEncodeNecessaryEscapes |
	FlagRemoveDotSegments |
	FlagRemoveFragment

// PureURL will normalize url with default flags
func PureURL(u *url.URL) string {
	return NormalizeURL(u, URLNormalizeFlag)
}

// PureURLString is like PureURL while url string may be invalid
// then an error may occur indicate such a parse error
func PureURLString(u string) (string, error) {
	return NormalizeURLString(u, URLNormalizeFlag);
}
