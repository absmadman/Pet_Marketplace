package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	router *gin.Engine
	psql   *sql.DB
	redis  *redis.Client
}

type RegisterInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
