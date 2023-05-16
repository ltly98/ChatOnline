package routes

import (
	v1 "chat-online/api/v1"
	"chat-online/middleware"
	"chat-online/utils"

	"github.com/gin-gonic/gin"
)

func InitRouter() {

	gin.SetMode(utils.AppMode)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	auth := r.Group("/api/v1")
	auth.Use(middleware.JwtToken())
	{
		//websocket
		//auth.GET("/ws", v1.WebSocket)
	}

	//此处只有登陆、注册、获取验证码和忘记密码不需要认证
	router := r.Group("/api/v1")
	{
		router.GET("/ws", v1.WebSocket)
		router.POST("/login", v1.Login)
		router.POST("/register", v1.Register)
		router.POST("/forget", v1.Forget)
		router.POST("/email", v1.SendCode)
		router.PUT("/email", v1.ReplaceEmail)
	}

	r.Run()
}
