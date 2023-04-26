package utils

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetRealIP(c *gin.Context) string {
	forwarded := c.Request.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ", ")
		if len(ips) > 1 {
			return ips[len(ips)-1]
		} else {
			return ips[0]
		}
	}
	return c.Request.RemoteAddr
}
