package routes

import (
	"chatgpt-web-go/src/api/auth"
	"chatgpt-web-go/src/api/user"
	"chatgpt-web-go/src/middleware"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes() *gin.Engine {
	g := gin.New()
	apiRoutes := g.Group("/api")
	//TODO 鉴权
	apiRoutes.Use(middleware.JWTAuth())
	{
		apiRoutes.POST("/login", auth.Login)
		apiRoutes.POST("/usage", user.GetUsage)
		apiRoutes.POST("/config", user.GetConfig)
		apiRoutes.POST("/session", auth.GetSession)
	}
	return g
}
