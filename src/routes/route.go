package routes

import (
	"chatgpt-web-go/src/api/auth"
	"chatgpt-web-go/src/api/gpt"
	"chatgpt-web-go/src/api/user"
	"chatgpt-web-go/src/middleware"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes() *gin.Engine {
	g := gin.New()
	apiRoutes := g.Group("/api")
	apiRoutes.Use(middleware.JWTAuth()).Use(middleware.UnauthorizedHandler())
	{
		apiRoutes.POST("/login", auth.Login)
		apiRoutes.POST("/chat-process", gpt.GPT)
		apiRoutes.POST("/logout", auth.Logout)
		apiRoutes.POST("/usage", user.GetUsage)
		apiRoutes.POST("/config", user.GetConfig)
		apiRoutes.POST("/session", auth.GetSession)
	}
	return g
}
