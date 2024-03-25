package main

import (
	"VK_Internship_Marketplace/internal/server"
	"VK_Internship_Marketplace/pkg/repository/db"
	"VK_Internship_Marketplace/pkg/repository/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	db := db.NewDBConnection()
	redis := redis_pkg.NewRedisConnection()
	handler := server.NewHandler(gin.Default(), db, redis)
	handler.HttpServer()
}
