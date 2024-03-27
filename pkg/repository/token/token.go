package jwttoken

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type ClaimsToken struct {
	jwt.RegisteredClaims
	Id int
}

// Token структура для хранения jwt токена в формате строки
type Token struct {
	Token string
}

// NewTokenFromCtx получает токен из хедера и возвращает структуру Token
func NewTokenFromCtx(ctx *gin.Context) (*Token, error) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return nil, errors.New("empty header")
	}
	splitToken := strings.Split(header, " ")
	if len(splitToken) != 2 {
		return nil, errors.New("empty token")
	}
	return &Token{Token: splitToken[1]}, nil
}

// NewTokenFromId создает токен из id и возвращает структуру Token
func NewTokenFromId(id int) (*Token, error) {
	tokenCfg := jwt.NewWithClaims(jwt.SigningMethodHS256, &ClaimsToken{jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	},
		id,
	})
	token, err := tokenCfg.SignedString([]byte(signedString))
	if err != nil {
		return nil, err
	}
	return &Token{Token: token}, nil
}

// GetId возвращает токен из структуры Token
func (t *Token) GetId() (int, error) {
	parsedToken, err := jwt.ParseWithClaims(t.Token, &ClaimsToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(signedString), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := parsedToken.Claims.(*ClaimsToken)
	if !ok {
		return 0, errors.New("invalid token")
	}
	return claims.Id, nil
}
