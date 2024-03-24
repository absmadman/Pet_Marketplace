package server

import (
	"VK_Internship_Marketplace/pkg/repository/db"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strings"
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
	ctx.IndentedJSON(http.StatusCreated, &user)
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
	splitToken := strings.Split(header, " ")
	if len(splitToken) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "error header format"})
		return
	}
	claims, err := parseToken(splitToken[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "error parsing token"})
		return
	}
	ctx.Set("UserID", claims)
}

func (h *Handler) advList(ctx *gin.Context) { // to do
	authUserId := h.GetIdByTokenIfExist(ctx)
	var al db.AdvList
	var filter db.Filter
	page, err := getIntegerParam("page", ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid parameter"})
		return
	}
	err = ctx.BindJSON(&filter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "error filter parameters"})
		return
	}
	err = al.GetAdvList(page, h.psql, &filter, authUserId)
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

func (h *Handler) removeAdvert(ctx *gin.Context) {
	var adv db.Advert
	currUserId, ok := ctx.Get("UserID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	advertId, err := getIntegerParam("advert_id", ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid parameter"})
		return
	}
	err = adv.GetAdv(h.psql, advertId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "advert is not exist"})
		return
	}
	if adv.UserId != currUserId.(int) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "you dont have permissions"})
		return
	}
	err = adv.RemoveAdvert(h.psql)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.JSON(http.StatusOK, adv)
}

func (h *Handler) getAdvert(ctx *gin.Context) {
	var adv db.Advert
	authUserId := h.GetIdByTokenIfExist(ctx)
	advId, err := getIntegerParam("advert_id", ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid parameter"})
		return
	}
	err = adv.GetAdv(h.psql, advId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	if authUserId == adv.UserId {
		adv.ByThisUser = true
	}
	ctx.JSON(http.StatusOK, adv)
}

func (h *Handler) advFeed(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, gin.H{"message": "adv feeding"})
}

func (h *Handler) HttpServer() {
	auth := h.router.Group("/auth")
	{
		auth.POST("/register", h.signUp)
		auth.POST("/login", h.signIn)
	}
	h.router.GET("/api/feed", h.advList)
	h.router.GET("/api/advert", h.getAdvert)
	adv := h.router.Group("/api", h.checkAuth)
	{

		adv.POST("/advert", h.addAdvert)
		adv.DELETE("/advert", h.removeAdvert)
		//adv.GET("/feed", h.advList)
	}

	err := h.router.Run("localhost:8080")
	if err != nil {
		log.Println(err)
	}
}
