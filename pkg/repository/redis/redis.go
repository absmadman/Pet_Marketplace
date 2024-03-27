package redisPkg

import (
	"VK_Internship_Marketplace/config"
	"VK_Internship_Marketplace/internal/entities"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	redis           *redis.Client
	storingDuration time.Duration
}

// NewRedis констуктор для структуры Redis
func NewRedis(cfg *config.Config) *Redis {
	return &Redis{
		redis:           NewRedisConnection(cfg),
		storingDuration: time.Minute * time.Duration(cfg.RedisStoringDuration),
	}
}

// NewRedisConnection создает подключение к Redis хранилищу
func NewRedisConnection(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		/*
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		*/
		Addr:     "redis_db:6379",
		Password: "",
		DB:       0,
	})
}

// AppendCache сериализует структуру Advert и добавляет в кэш Redis
func (r *Redis) AppendCache(id int, a *entities.Advert) {
	bytes, err := a.MarshalBinary()
	if err != nil {
		return
	}
	r.redis.Set(context.Background(), string(id), bytes, r.storingDuration)
}

// GetFromCache вытаскивает данные из Redis по id и десеарилизует их в структуру Advert
func (r *Redis) GetFromCache(id int) (*entities.Advert, error) {
	var a *entities.Advert
	res := r.redis.Get(context.Background(), string(id))
	bytes, err := res.Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &a)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, errors.New("empty")
	}
	return a, nil
}

// RemoveFromCache удаляет данные из Redis по id
func (r *Redis) RemoveFromCache(id int) {
	r.redis.Del(context.Background(), string(id))
}
