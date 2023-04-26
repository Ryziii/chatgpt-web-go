package utils

import (
	"chatgpt-web-go/global"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"time"
)

var jwtSecret []byte

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username, password string, exHour int) (string, error) {
	jwtSecret = []byte(global.Cfg.App.JwtSecret)
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(exHour) * time.Hour)

	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "chat",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	jwtSecret = []byte(global.Cfg.App.JwtSecret)
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func IsTokenExpired(token string) bool {
	claims, err := ParseToken(token)
	if err != nil {
		global.Gzap.Error("IsTokenExpired", zap.Error(err))
		return false
	}
	return claims.ExpiresAt < time.Now().Unix()
}
