package main

import (
	"VK_Internship_Marketplace/config"
	"VK_Internship_Marketplace/internal/server"
	"VK_Internship_Marketplace/pkg/repository/db"
	"VK_Internship_Marketplace/pkg/repository/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()
	if cfg == nil {
		return
	}
	psql := db.NewDatabase(cfg)
	redis := redisPkg.NewRedis(cfg)
	handler := server.NewHandler(gin.Default(), psql, redis)
	handler.HttpServer(cfg)
}
