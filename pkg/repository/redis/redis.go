package redis_pkg

import (
	"VK_Internship_Marketplace/internal/entities"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

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
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func (r *Redis) AppendCache(id int, a *entities.Advert) {
	r.redis.Set(context.Background(), string(id), a, 0)
}

func (r *Redis) GetFromCache(id int) (*entities.Advert, error) {
	var a *entities.Advert
	err := r.redis.Get(context.Background(), string(id)).Scan(a)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, errors.New("empty")
	}
	return a, nil
}
