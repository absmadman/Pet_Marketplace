package server

import (
	"VK_Internship_Marketplace/pkg/repository/db"
	redis_pkg "VK_Internship_Marketplace/pkg/repository/redis"
	jwttoken "VK_Internship_Marketplace/pkg/repository/token"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	router *gin.Engine
	psql   *db.Database
	redis  *redis_pkg.Redis
}

func NewHandler(engine *gin.Engine, db *db.Database, redis *redis_pkg.Redis) *Handler {
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	if !isValid(user.Password) || !isValid(user.Login) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid login or password"})
		return
	}
	err = h.psql.CreateUser(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, &user)
}

func (h *Handler) signIn(ctx *gin.Context) {
	var user db.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error"})
		return
	}
	err = h.psql.GetUser(&user)
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

func (h *Handler) advList(ctx *gin.Context) { // to do
	authUserId := h.GetIdByTokenIfExist(ctx)
	var al db.AdvList
	var filter db.Filter
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
	err = h.psql.GetAdvList(page, &al, &filter, authUserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error database get"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, al)
}

func (h *Handler) addAdvert(ctx *gin.Context) {
	var adv db.Advert
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
	if err = h.psql.CreateAdvert(&adv); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "advert is not valid"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, adv)
}

func (h *Handler) removeAdvert(ctx *gin.Context) {
	var adv db.Advert
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
	err = h.psql.RemoveAdvert(&adv)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	ctx.IndentedJSON(http.StatusOK, adv)
}

func (h *Handler) getAdvert(ctx *gin.Context) {
	var adv db.Advert
	authUserId := h.GetIdByTokenIfExist(ctx)
	advId, err := getIntegerParam("advert_id", ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid parameter"})
		return
	}
	err = h.psql.GetAdv(&adv, advId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	if authUserId == adv.UserId {
		adv.ByThisUser = true
	}
	ctx.IndentedJSON(http.StatusOK, adv)
}

func (h *Handler) updateAdvert(ctx *gin.Context) {
	var adv db.Advert
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
	err = h.psql.UpdateAdvert(&adv)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "error updating database"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, adv)
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
		adv.PUT("/advert", h.updateAdvert)
	}

	err := h.router.Run("localhost:8080")
	if err != nil {
		log.Println(err)
	}
}
