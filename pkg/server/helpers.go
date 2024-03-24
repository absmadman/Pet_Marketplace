package server

import (
	"VK_Internship_Marketplace/pkg/repository/db"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-passwd/validator"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
	"time"
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
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return 0
	}
	splitToken := strings.Split(header, " ")
	if len(splitToken) != 2 {
		return 0
	}
	claims, err := parseToken(splitToken[1])
	if err != nil {
		return 0
	}
	return claims
}

func createToken(u db.User) (string, error) {
	tokenCfg := jwt.NewWithClaims(jwt.SigningMethodHS256, &Token{jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	},
		u.Id,
	})
	token, err := tokenCfg.SignedString([]byte(signedString))
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
		return []byte(signedString), nil
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
