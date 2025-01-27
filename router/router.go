package router

import (
	"github.com/gin-gonic/gin"
	"note_app_server1/controller"
	"note_app_server1/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	auth := r.Group("/auth")
	{
		auth.POST("/register", func(context *gin.Context) {
			controller.Register(context)
		})
		auth.POST("/login", func(context *gin.Context) {
			controller.Login(context)
		})
		auth.POST("/logout", func(context *gin.Context) {
			controller.Logout(context)
		})
	}
	avatar := r.Group("/avatar")
	{
		avatar.GET("/:fileName", func(context *gin.Context) {
			controller.GetAvatarUrl(context)
		})
	}
	return r
}
