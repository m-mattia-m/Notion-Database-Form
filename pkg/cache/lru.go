package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// TODO: The LRU cache would still need a logic with the popularity and the last call date.
//  Add this corresponding logic.

type LRUCache struct {
	redisClient *redis.Client
	capacity    int
}

func NewLRUCache(redisAddr, redisPassword string, capacity int) *LRUCache {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	return &LRUCache{
		redisClient: redisClient,
		capacity:    capacity,
	}
}

func (c *LRUCache) Get(key string) (*string, error) {
	val, err := c.redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if val == "" {
		return nil, fmt.Errorf("key exist but is empty")
	}

	c.redisClient.Expire(context.Background(), key, 30*24*time.Hour)
	return &val, nil
}

func (c *LRUCache) Set(key string, value string) error {
	err := c.redisClient.Set(context.Background(), key, fmt.Sprintf("%v", value), 0).Err()
	if err != nil {
		return err
	}
	return nil
}
