package cache

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Keys(context.Context, string) *redis.StringSliceCmd
	Del(context.Context, ...string) *redis.IntCmd
}

func NewRedisClient() *redis.Client {
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379" // Default Redis internal port
	}
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + redisPort,
		Password: "",
		DB:       0,
	})
}
