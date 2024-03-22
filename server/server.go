package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewHandler(engine *gin.Engine, db *sql.DB, redis *redis.Client) *Handler {
	return &Handler{
		router: engine,
		psql:   db,
		redis:  redis,
	}
}

func (handler *Handler) signUp(ctx *gin.Context) {

}

func (handler *Handler) signIn(ctx *gin.Context) {

}

func (handler *Handler) HttpServer() {

	handler.router.POST("/register", handler.signUp)

	err := handler.router.Run("localhost:8080")
	if err != nil {
		log.Println(err)
	}
}
