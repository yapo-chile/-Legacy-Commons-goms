package infrastructure

import (
	"crypto/md5" //nolint: gosec
	"encoding/json"
	"fmt"

	"github.com/Yapo/goutils"
	"github.com/anevsky/cachego/memory"
)

type RequestCache struct {
	enabled  bool
	cache    memory.CACHE
	cacheTTL int
}

func (rc *RequestCache) getHash(data interface{}) string {
	sum := md5.Sum([]byte(fmt.Sprintf("%v", data))) //nolint: gosec
	return fmt.Sprintf("%x", sum)
}

func (rc *RequestCache) GetCache(input interface{}) (*goutils.Response, error) {
	var response goutils.Response
	if rc.enabled {
		hash := rc.getHash(input)
		if stringResponse, err := rc.cache.Get(hash); err == nil {
			err = json.Unmarshal([]byte(stringResponse.(string)), &response)
			return &response, err
		}
	}
	return &response, fmt.Errorf("cache disabled")
}

func (rc *RequestCache) SetCache(input interface{}, response *goutils.Response) error {
	if rc.enabled {
		hash := rc.getHash(input)
		stringResponse, err := json.Marshal(response)
		if err != nil {
			return err
		}

		if err := rc.cache.SetString(hash, string(stringResponse)); err != nil {
			return err
		}
		return rc.cache.SetTTL(hash, rc.cacheTTL)
	}
	return fmt.Errorf("cache disabled")
}

func NewRequestCacheHandler(ttl int) *RequestCache {
	return &RequestCache{
		cache:    memory.Alloc(),
		cacheTTL: ttl,
		enabled:  true,
	}
}
