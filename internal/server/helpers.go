package server

import (
	"VK_Internship_Marketplace/internal/entities"
	jwttoken "VK_Internship_Marketplace/pkg/repository/token"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-passwd/validator"
	"net/http"
	"strconv"
)

// getIntegerParam возвращает целочисленный параметр по имени параметра из контекста
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

// GetIdByTokenIfExist если Token содержит id возвращает id иначе 0
func (h *Handler) GetIdByTokenIfExist(ctx *gin.Context) int {
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

// isValid проверяет пароль на валидность
func isValid(pass string) bool {
	valid := validator.New(validator.MinLength(8, nil),
		validator.MaxLength(36, nil),
		validator.ContainsOnly(entities.AllowedSymbols, nil))
	err := valid.Validate(pass)
	if err != nil {
		return false
	}
	return true
}

// checkAdvertOwnership проверяет владеет ли текущий пользователь объявлением
func (h *Handler) checkAdvertOwnership(ctx *gin.Context, adv *entities.Advert, userId int, advertId int) bool {
	err := h.db.GetAdv(adv, advertId)
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
