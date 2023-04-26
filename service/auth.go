package service

import (
	global2 "chatgpt-web-go/global"
	"chatgpt-web-go/global/enum"
	"chatgpt-web-go/utils"
	"chatgpt-web-go/utils/redis"
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
	token, _ = utils.GenerateToken(username, "", global2.Cfg.App.JwtExpires)
	if err := redis.Set(key, token, time.Duration(global2.Cfg.App.JwtExpires)*time.Hour); err != nil {
		return "", &global2.SystemError{Msg: err.Error()}
	}
	return token, nil
}

func Logout(token string) (bool, error) {
	parseToken, err := utils.ParseToken(token)
	if err != nil {
		return false, &global2.SystemError{Msg: err.Error()}
	}
	username := parseToken.Username
	key := enum.CACHE_USER + username + ":"
	if redis.Exist(key) {
		if err := redis.LikeDeletes(key); err != nil {
			return false, &global2.SystemError{Msg: err.Error()}
		}
	}
	return true, nil
}
