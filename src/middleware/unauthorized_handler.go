package middleware

import (
	"github.com/gin-gonic/gin"
)

func UnauthorizedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Writer.Status() == 401 {
			c.SetCookie("token", "", -1, "/", "", false, false)
			c.Abort()
			return
		}
		c.Next()
	}
}
