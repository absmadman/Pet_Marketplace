package server

import (
	"VK_Internship_Marketplace/pkg/repository/db"
	jwttoken "VK_Internship_Marketplace/pkg/repository/token"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-passwd/validator"
	"net/http"
	"strconv"
)

func getIntegerParam(paramName string, ctx *gin.Context) (int, error) {
	param := ctx.Query(paramName)
	if param == "" {
		return 0, errors.New("empty param")
	}
	intParam, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return intParam, nil
}

func (h *Handler) GetIdByTokenIfExist(ctx *gin.Context) int {
	/*
		var token *jwttoken.Token
		header := ctx.GetHeader("Authorization")
		if header == "" {
			return 0
		}
		splitToken := strings.Split(header, " ")
		if len(splitToken) != 2 {
			return 0
		}
		claims, err := token.ParseToken(splitToken[1])
	*/
	token, err := jwttoken.NewTokenFromCtx(ctx)
	if err != nil {
		return 0
	}
	id, err := token.GetId()
	if err != nil {
		return 0
	}
	return id
}

func isValid(pass string) bool {
	valid := validator.New(validator.MinLength(8, nil),
		validator.MaxLength(36, nil),
		validator.ContainsOnly(db.AllowedSymbols, nil))
	err := valid.Validate(pass)
	if err != nil {
		return false
	}
	return true
}

func (h *Handler) checkAdvertOwnership(ctx *gin.Context, adv *db.Advert, userId int, advertId int) bool {
	err := h.psql.GetAdv(adv, advertId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "advert is not exist"})
		return false
	}
	if adv.UserId != userId {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "you dont have permissions"})
		return false
	}
	return true
}
