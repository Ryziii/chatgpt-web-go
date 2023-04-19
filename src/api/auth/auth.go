package auth

import (
	"chatgpt-web-go/src/model/api"
	"chatgpt-web-go/src/model/api/user/request"
	"chatgpt-web-go/src/model/api/user/response"
	result "chatgpt-web-go/src/model/common/response"
	"chatgpt-web-go/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSession(c *gin.Context) {
	c.SetCookie("name", "123", 3600, "/", "localhost", false, false)
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
		c.JSON(http.StatusBadRequest, result.Fail.WithData(err.Error()))
		return
	}
	if token, err := service.Login(requestP.Username, requestP.Password); err != nil {
		c.JSON(http.StatusBadRequest, result.Fail.WithData(err.Error()))
	} else {
		c.JSON(http.StatusOK, result.OK.WithData(response.LoginResponse{Token: token}))
	}
}
