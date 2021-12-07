package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinsoft/it-ku/user-service/model"
	"github.com/jinsoft/it-ku/user-service/repo"
	"time"
)

var key = []byte("ikUserTokenKeySecret")

type CustomClaims struct {
	User *model.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *model.User) (string, error)
}

type TokenService struct {
	Repo repo.Repository
}

func (srv *TokenService) Decode(tokenString string) (*CustomClaims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	// 验证令牌并返回自定义声明
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (srv *TokenService) Encode(user *model.User) (string, error) {
	// 配置化
	expireToken := time.Now().Add(time.Hour * 72).Unix()

	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "ik.user.service",
		},
	}
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// token 签名并返回
	return token.SignedString(key)
}
