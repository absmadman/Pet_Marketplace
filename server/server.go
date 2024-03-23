package server

import (
	"VK_Internship_Marketplace/db"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-passwd/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

type Token struct {
	jwt.RegisteredClaims
	Id int
}

type Handler struct {
	router *gin.Engine
	psql   *sql.DB
	redis  *redis.Client
}

func NewHandler(engine *gin.Engine, db *sql.DB, redis *redis.Client) *Handler {
	return &Handler{
		router: engine,
		psql:   db,
		redis:  redis,
	}
}

func isPasswordValid(pass string) bool {
	valid := validator.New(validator.MinLength(8, nil), validator.MaxLength(36, nil), validator.ContainsOnly("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", nil))
	err := valid.Validate(pass)
	if err != nil {
		return false
	}
	return true
}

func (h *Handler) signUp(ctx *gin.Context) {
	var user db.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	if !isPasswordValid(user.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
		return
	}
	err = user.CreateUser(h.psql)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, &user)
}

func (h *Handler) signIn(ctx *gin.Context) {
	var user db.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	err = user.GetUser(h.psql)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	tokenCfg := jwt.NewWithClaims(jwt.SigningMethodHS256, &Token{jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	},
		user.Id,
	})
	token, err := tokenCfg.SignedString([]byte("dsajkfashfaklajhf13"))
	// переписать этот кусок кода
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.IndentedJSON(http.StatusOK, &token)
}

func (h *Handler) HttpServer() {

	h.router.POST("/register", h.signUp)
	h.router.POST("/login", h.signIn)

	err := h.router.Run("localhost:8080")
	if err != nil {
		log.Println(err)
	}
}
