package server

import (
	"VK_Internship_Marketplace/pkg/repository/db"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-passwd/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func isValid(pass string) bool {
	valid := validator.New(validator.MinLength(8, nil), validator.MaxLength(36, nil), validator.ContainsOnly("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!_", nil))
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
	if !isValid(user.Password) || !isValid(user.Login) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid login or password"})
		return
	}
	err = user.CreateUser(h.psql)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, &user)
}

func createToken(u db.User) (string, error) {
	tokenCfg := jwt.NewWithClaims(jwt.SigningMethodHS256, &Token{jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	},
		u.Id,
	})
	token, err := tokenCfg.SignedString([]byte("dsajkfashfaklajhf13"))
	if err != nil {
		return "", err
	}
	return token, nil
}

func parseToken(token string) (int, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte("dsajkfashfaklajhf13"), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := parsedToken.Claims.(*Token)
	if !ok {
		return 0, errors.New("invalid token")
	}
	return claims.Id, nil
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
	token, err := createToken(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, &token)
}

func (h *Handler) checkAuth(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "empty header"})
		return
	}
	splitted := strings.Split(header, " ")
	if len(splitted) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "error header format"})
		return
	}
	claims, err := parseToken(splitted[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "error parsing token"})
		return
	}
	ctx.Set("UserID", claims)
}

func (h *Handler) AdvList(ctx *gin.Context) { // to do
	//id, _ := ctx.Get("UserID")
	page := ctx.Query("page")
	pageid, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error parameter page"})
		return
	}
	var al db.AdvList
	err = al.GetAdvList(pageid, h.psql, true, 0, 1000)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error database get"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, al)
}

func (h *Handler) addAdvert(ctx *gin.Context) {
	var adv db.Advert
	err := ctx.BindJSON(&adv)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error json"})
		return
	}
	ok := adv.ValidateAdvertData()
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid advert data"})
		return
	}
	id, ok := ctx.Get("UserID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	adv.UserId = id.(int)
	if err = adv.CreateAdvert(h.psql); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "advert is not valid"})
		return
	}
	ctx.JSON(http.StatusOK, adv)
}

func (h *Handler) HttpServer() {
	auth := h.router.Group("/auth")
	{
		auth.POST("/register", h.signUp)
		auth.POST("/login", h.signIn)
	}
	adv := h.router.Group("/api", h.checkAuth)
	{
		adv.POST("/advert", h.addAdvert)
		adv.POST("/feed", h.AdvList)
		/*
			list := h.router.Group(":page/list")
			{

			}
		*/
	}

	err := h.router.Run("localhost:8080")
	if err != nil {
		log.Println(err)
	}
}
