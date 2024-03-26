package redis_pkg

import "github.com/redis/go-redis/v9"

type Redis struct {
	redis *redis.Client
}

func NewRedis() *Redis {
	return &Redis{
		redis: NewRedisConnection(),
	}
}

func NewRedisConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}
