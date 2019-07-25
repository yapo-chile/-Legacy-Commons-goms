package infrastructure

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"

	"github.mpi-internal.com/Yapo/goms/pkg/interfaces/handlers"
)

// cache is a handler to get cached request responses using redis
type cache struct {
	handler RedisHandler
	prefix  string
	maxAge  time.Duration
}

// NewCacheHandler returns an instance of cache handler
func NewCacheHandler(handler RedisHandler, prefix string, maxAge time.Duration) handlers.CacheHandler {
	return &cache{
		handler: handler,
		prefix:  prefix,
		maxAge:  maxAge,
	}
}

// makeRedisKey generates sha1 key taking request input data
func (c *cache) makeRedisKey(input interface{}) (string, error) {
	inputRaw, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	h := sha1.New()
	h.Write(inputRaw) // nolint
	key := fmt.Sprintf("%x", h.Sum(nil))
	return c.prefix + key, nil
}

// GetCache returns cached response for input request
func (c *cache) GetCache(input interface{}) (json.RawMessage, error) {
	var response json.RawMessage
	key, err := c.makeRedisKey(input)
	if err != nil {
		return response, err
	}
	res, err := c.handler.Get(key)
	if err == nil {
		response, err = res.Bytes()
	}
	return response, err
}

// SetCache saves the response for input request
func (c *cache) SetCache(input interface{}, response json.RawMessage) error {
	key, err := c.makeRedisKey(input)
	if err != nil {
		return err
	}
	return c.handler.Set(key, string(response), c.maxAge)
}
