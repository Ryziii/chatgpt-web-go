package middleware

import (
	"chatgpt-web-go/src/global/enum"
	result "chatgpt-web-go/src/model/common/response"
	"chatgpt-web-go/src/utils"
	"chatgpt-web-go/src/utils/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/login" || c.Request.URL.Path == "/api/session" || c.Request.URL.Path == "/api/logout" {
			c.Next()
			return
		}
		token, _ := c.Cookie("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, result.NotAuth)
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, result.NotAuth)
			return
		}
		key := enum.CACHE_USER + claims.Username + ":"
		if tokenRedis, err := redis.Get(key); err != nil {
			c.JSON(http.StatusUnauthorized, result.NotAuth)
			return
		} else {
			if tokenRedis != token {
				c.JSON(http.StatusUnauthorized, result.NotAuth)
				return
			}
		}

		c.Set("claims", claims)
		c.Next()
	}
}
