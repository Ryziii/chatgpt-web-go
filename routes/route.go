package routes

import (
	"chatgpt-web-go/api/auth"
	"chatgpt-web-go/api/gpt"
	"chatgpt-web-go/api/user"
	middleware2 "chatgpt-web-go/middleware"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes() *gin.Engine {
	g := gin.New()
	apiRoutes := g.Group("/api")
	apiRoutes.Use(middleware2.JWTAuth()).Use(middleware2.UnauthorizedHandler())
	{
		apiRoutes.POST("/login", auth.Login)
		apiRoutes.POST("/chat-process", gpt.ChatConversationProcess)
		apiRoutes.POST("/add-chat-room", gpt.AddChatRoom)
		apiRoutes.POST("/logout", auth.Logout)
		apiRoutes.POST("/usage", user.GetUsage)
		apiRoutes.POST("/config", user.GetConfig)
		apiRoutes.POST("/session", auth.GetSession)
	}
	return g
}
