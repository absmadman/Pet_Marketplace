package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	PGDBName             string
	PGLogin              string
	PGPassword           string
	RedisAddress         string
	RedisPassword        string
	RedisDB              int
	RedisStoringDuration int
	ServerAddress        string
	ServerPort           string
}

func NewConfig() *Config {
	duration, err := strconv.Atoi(os.Getenv("REDIS_STORING_DURATION"))
	if err != nil {
		log.Println("Invalid parameter REDIS_STORING_DURATION")
		return nil
	}
	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Println("Invalid parameter REDIS_DB")
		return nil
	}
	return &Config{
		PGDBName:             os.Getenv("PG_DB_NAME"),
		PGLogin:              os.Getenv("PG_LOGIN"),
		PGPassword:           os.Getenv("PG_PASSWORD"),
		RedisAddress:         os.Getenv("REDIS_ADDRESS"),
		RedisPassword:        os.Getenv("REDIS_PASSWORD"),
		RedisDB:              redisDb,
		RedisStoringDuration: duration,
		ServerAddress:        os.Getenv("SERVER_ADDRESS"),
		ServerPort:           os.Getenv("SERVER_PORT"),
	}
}
