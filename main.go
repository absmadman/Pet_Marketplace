package main

import (
	"VK_Internship_Marketplace/db"
	redis_pkg "VK_Internship_Marketplace/redis"
	"VK_Internship_Marketplace/server"
	"github.com/gin-gonic/gin"
)

func main() {
	db := db.NewDBConnection()
	redis := redis_pkg.NewRedisConnection()
	handler := server.NewHandler(gin.Default(), db, redis)
	handler.HttpServer()
}
