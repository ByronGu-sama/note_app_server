package router

import (
	"github.com/gin-gonic/gin"
	"note_app_server1/controller"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	auth := r.Group("/auth")
	{
		auth.POST("/register", func(context *gin.Context) {
			controller.Register(context)
		})
		auth.POST("/login", func(context *gin.Context) {
			controller.Login(context)
		})
	}
	return r
}
