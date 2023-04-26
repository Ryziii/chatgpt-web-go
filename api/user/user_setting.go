package user

import (
	"chatgpt-web-go/model/api/user/setting"
	result "chatgpt-web-go/model/common/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetConfig(c *gin.Context) {
	// TODO 从数据库中获取用户配置
	c.JSON(http.StatusOK, result.OK.WithData(
		setting.Usage{
			ApiModel:     "12312312321",
			Usage:        "1231231312",
			ReverseProxy: "123",
			TimeoutMs:    123,
			SocksProxy:   "",
			HttpsProxy:   "",
		}),
	)
}

func GetUsage(c *gin.Context) {
	// TODO 从数据库中获取用户配置
	c.JSON(http.StatusOK, result.OK.WithData(
		setting.Usage{
			ApiModel:     "12312312321",
			Usage:        "1231231312",
			ReverseProxy: "123",
			TimeoutMs:    123,
			SocksProxy:   "",
			HttpsProxy:   "",
		}),
	)
}
