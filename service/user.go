package service

import (
	"chatgpt-web-go/global"
	"chatgpt-web-go/model/api/user"
	"chatgpt-web-go/utils"
	"errors"
)

func VerifyPassword(username, password string) (bool, error) {
	user := user.User{}
	userDB := global.Gdb.Model(&user)
	e := userDB.Where("username = ?", username).First(&user).Error
	if e != nil {
		return false, e
	}
	if user.ID == 0 {
		return false, errors.New("user not found")
	}
	if utils.MD5(utils.MD5(password)) != user.Password {
		return false, errors.New("password error")
	}
	return true, nil
}
