package infrastructure

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// RedisHandler handler to connect with redis
type RedisHandler struct {
	Client *redis.Client
}

// RedisResult allows operation over redis result
type RedisResult interface {
	Bytes() ([]byte, error)
}

// NewRedisHandler creates a new instance of redis handler
func NewRedisHandler(address, password string, db int) RedisHandler {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	return RedisHandler{
		Client: client,
	}
}

// Get gets the result of a GET command with the given key
func (r *RedisHandler) Get(key string) (RedisResult, error) {
	result := r.Client.Get(key)
	err := result.Err()
	if err == redis.Nil {
		return result, fmt.Errorf("KEY_NOT_FOUND: %s", key)
	}
	return result, err
}

// Set sets a value in redis with the given key
func (r *RedisHandler) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(key, value, expiration).Err()
}

// Del deletes the given key in redis
func (r *RedisHandler) Del(key string) error {
	return r.Client.Del(key).Err()
}
