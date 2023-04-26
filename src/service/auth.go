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
	key := enum.CACHE_USER + username + ":"
	if redis.Exist(key) {
		redis.LikeDeletes(key)
	}
	token, _ = utils.GenerateToken(username, "", global.Cfg.App.JwtExpires)
	if err := redis.Set(key, token, time.Duration(global.Cfg.App.JwtExpires)*time.Hour); err != nil {
		return "", &global.SystemError{Msg: err.Error()}
	}
	return token, nil
}

func Logout(token string) (bool, error) {
	parseToken, err := utils.ParseToken(token)
	if err != nil {
		return false, &global.SystemError{Msg: err.Error()}
	}
	username := parseToken.Username
	key := enum.CACHE_USER + username + ":"
	if redis.Exist(key) {
		if err := redis.LikeDeletes(key); err != nil {
			return false, &global.SystemError{Msg: err.Error()}
		}
	}
	return true, nil
}
