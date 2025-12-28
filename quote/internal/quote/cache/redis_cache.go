package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{
		client: client,
	}
}

// GetTodayQuote retrieves today's quote from cache
func (rc *RedisCache) GetTodayQuote(ctx context.Context) (interface{}, error) {
	key := getTodayKey()
	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}

	return data, nil
}

// SetTodayQuote caches today's quote with 24-hour TTL
func (rc *RedisCache) SetTodayQuote(ctx context.Context, quote interface{}) error {
	key := getTodayKey()
	jsonData, err := json.Marshal(quote)
	if err != nil {
		return err
	}

	return rc.client.Set(ctx, key, jsonData, 24*time.Hour).Err()
}

// getTodayKey returns a cache key for today's date
func getTodayKey() string {
	return "quote:today:" + time.Now().Format("2006-01-02")
}

// Close closes the Redis connection
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}
