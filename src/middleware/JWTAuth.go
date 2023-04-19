package middleware

import (
	result "chatgpt-web-go/src/model/common/response"
	"chatgpt-web-go/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/login" || c.Request.URL.Path == "/api/session" {
			c.Next()
			return
		}
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, result.NotAuth)
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, result.NotAuth)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
