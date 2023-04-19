package service

import (
	"chatgpt-web-go/src/global"
	"chatgpt-web-go/src/global/enum"
	"chatgpt-web-go/src/utils"
	"chatgpt-web-go/src/utils/redis"
	"errors"
	"time"
)

func Login(username, password string) (token string, err error) {
	if verify, err := VerifyPassword(username, password); err != nil {
		return "", err
	} else if !verify {
		return "", errors.New("password error")
	}
	key := enum.CACHE_USER + "_" + username
	if redis.Exist(key) {
		redis.LikeDeletes(key)
	}
	token, _ = utils.GenerateToken(username, "", global.Cfg.App.JwtExpires)
	if err := redis.Set(key, token, time.Duration(global.Cfg.App.JwtExpires)*time.Hour); err != nil {
		return "", err
	}
	return token, nil
}
