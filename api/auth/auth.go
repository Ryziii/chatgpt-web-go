package auth

import (
	"chatgpt-web-go/global"
	"chatgpt-web-go/model/api"
	"chatgpt-web-go/model/api/auth"
	"chatgpt-web-go/model/api/user/request"
	"chatgpt-web-go/model/api/user/response"
	result "chatgpt-web-go/model/common/response"
	"chatgpt-web-go/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSession(c *gin.Context) {
	c.JSON(http.StatusOK, result.OK.WithData(
		api.Session{
			Auth:  false,
			Model: "ChatGPTAPI",
		},
	))
}

func Login(c *gin.Context) {
	requestP := request.User{}
	if err := c.ShouldBindJSON(&requestP); err != nil {
		c.JSON(http.StatusOK, result.Fail.WithData(err.Error()))
		return
	}
	if token, err := service.Login(requestP.Username, requestP.Password); err != nil {
		switch err.(type) {
		case *global.SystemError:
			c.JSON(http.StatusOK, result.Fail.WithMessage("系统错误"))
		default:
			c.JSON(http.StatusOK, result.Fail.WithMessage("登录失败, 用户名或密码错误"))
		}
	} else {
		c.JSON(http.StatusOK, result.OK.WithData(response.LoginResponse{Token: token}))
	}
}

func Logout(c *gin.Context) {
	requestP := auth.Logout{}
	c.SetCookie("token", "", -1, "/", "", false, false)
	if err := c.ShouldBindJSON(&requestP); err != nil {
		c.JSON(http.StatusOK, result.Fail.WithMessage(err.Error()))
		return
	}
	if _, err := service.Logout(requestP.Token); err != nil {
		switch err.(type) {
		case *global.SystemError:
			c.JSON(http.StatusOK, result.Fail.WithMessage("系统错误"))
		default:
			c.JSON(http.StatusOK, result.Fail.WithMessage("登出失败"))
		}
	} else {
		c.JSON(http.StatusOK, result.OK.WithMessage("登出成功"))
	}
}
