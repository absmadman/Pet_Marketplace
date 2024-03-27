package server

import (
	"VK_Internship_Marketplace/config"
	"VK_Internship_Marketplace/internal/entities"
	"VK_Internship_Marketplace/pkg/repository/db"
	redisPkg "VK_Internship_Marketplace/pkg/repository/redis"
	jwttoken "VK_Internship_Marketplace/pkg/repository/token"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Handler структура для работы эндпоинтов
type Handler struct {
	router *gin.Engine
	db     *db.Database
	redis  *redisPkg.Redis
}

// NewHandler констуктор для Handler
func NewHandler(engine *gin.Engine, db *db.Database, redis *redisPkg.Redis) *Handler {
	return &Handler{
		router: engine,
		db:     db,
		redis:  redis,
	}
}

// signUp метод выполняющий регистрацию пользователей последством занесения данных пользователя в базу данных
func (h *Handler) signUp(ctx *gin.Context) {
	var user entities.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	if !isValid(user.Password) || !isValid(user.Login) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid login or password"})
		return
	}
	err = h.db.CreateUser(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, &user)
}

// signIn метод выполняющий проверку полученных данных пользователя с данными из бызы данных
func (h *Handler) signIn(ctx *gin.Context) {
	var user entities.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	err = h.db.GetUser(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	token, err := jwttoken.NewTokenFromId(user.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.IndentedJSON(http.StatusOK, &token.Token)
}

// checkAuth метод проверяет авторизован ли пользователь и валидность jwt токена
func (h *Handler) checkAuth(ctx *gin.Context) {
	token, err := jwttoken.NewTokenFromCtx(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "error parsing token"})
	}
	id, err := token.GetId()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err})
	}
	ctx.Set("UserID", id)
}

// advList метод собирает список объявлений по заданным фильтрам
func (h *Handler) advList(ctx *gin.Context) {
	var al entities.AdvList
	var filter entities.Filter
	authUserId := h.GetIdByTokenIfExist(ctx)
	page, err := getIntegerParam("page", ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid parameter"})
		return
	}
	err = ctx.BindJSON(&filter)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error filter parameters"})
		return
	}
	err = h.db.GetAdvList(page, &al, &filter, authUserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error database get"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, al)
}

// addAdvert метод парсит полученный json и заносит данные в базу данных, а также добавляет данные в кeш
func (h *Handler) addAdvert(ctx *gin.Context) {
	var adv entities.Advert
	err := ctx.BindJSON(&adv)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error json"})
		return
	}
	ok := adv.ValidateAdvertData()
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid advert data"})
		return
	}
	id, ok := ctx.Get("UserID")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	adv.UserId = id.(int)
	if err = h.db.CreateAdvert(&adv); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "advert is not valid"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, adv)
}

// removeAdvert удаляет объявление по id из базы данных и кeша
func (h *Handler) removeAdvert(ctx *gin.Context) {
	var adv entities.Advert
	currUserId, ok := ctx.Get("UserID")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	advertId, err := getIntegerParam("advert_id", ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid parameter"})
		return
	}
	ok = h.checkAdvertOwnership(ctx, &adv, currUserId.(int), advertId)
	if !ok {
		return
	}
	err = h.db.RemoveAdvert(&adv)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	h.redis.RemoveFromCache(adv.Id)
	ctx.IndentedJSON(http.StatusOK, adv)
}

// getAdvert возвращает объявление из кеша если они там есть, иначе из базы данных и добавляет их в кеш
func (h *Handler) getAdvert(ctx *gin.Context) {
	var adv entities.Advert
	authUserId := h.GetIdByTokenIfExist(ctx)
	advId, err := getIntegerParam("advert_id", ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid parameter"})
		return
	}
	cachedAdv, err := h.redis.GetFromCache(advId)
	if err == nil {
		adv = *cachedAdv
	} else {
		err = h.db.GetAdv(&adv, advId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "advert is not exist"})
			return
		}
	}
	if authUserId == adv.UserId {
		adv.ByThisUser = true
	}
	h.redis.AppendCache(adv.Id, &adv)
	ctx.IndentedJSON(http.StatusOK, adv)
}

// updateAdvert обновляет объялвение в базе данных и в кеше
func (h *Handler) updateAdvert(ctx *gin.Context) {
	var adv entities.Advert
	currUserId, ok := ctx.Get("UserID")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	advertId, err := getIntegerParam("advert_id", ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid parameter"})
		return
	}
	ok = h.checkAdvertOwnership(ctx, &adv, currUserId.(int), advertId)
	if !ok {
		return
	}
	err = ctx.BindJSON(&adv)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid advert data"})
		return
	}
	adv.Id = advertId
	err = h.db.UpdateAdvert(&adv)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error updating database"})
		return
	}
	h.redis.RemoveFromCache(adv.Id)
	h.redis.AppendCache(adv.Id, &adv)
	ctx.IndentedJSON(http.StatusOK, adv)
}

// HttpServer описание эндпоинтов
func (h *Handler) HttpServer(cfg *config.Config) {
	auth := h.router.Group("/auth")
	{
		auth.POST("/register", h.signUp)
		auth.GET("/login", h.signIn)
	}
	h.router.GET("/api/feed", h.advList)
	h.router.GET("/api/advert", h.getAdvert)
	adv := h.router.Group("/api", h.checkAuth)
	{
		adv.POST("/advert", h.addAdvert)
		adv.DELETE("/advert", h.removeAdvert)
		adv.PUT("/advert", h.updateAdvert)
	}

	err := h.router.Run(cfg.ServerAddress + ":" + cfg.ServerPort)
	if err != nil {
		log.Println(err)
	}
}
