package redis_pkg

import "github.com/redis/go-redis/v9"

func NewRedisConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}
