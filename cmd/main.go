package main

import (
	"VK_Internship_Marketplace/internal/server"
	"VK_Internship_Marketplace/pkg/repository/db"
	"VK_Internship_Marketplace/pkg/repository/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	psql := db.NewDatabase()
	redis := redisPkg.NewRedis()
	handler := server.NewHandler(gin.Default(), psql, redis)
	handler.HttpServer()
}
