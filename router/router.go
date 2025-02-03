package router

import (
	"github.com/gin-gonic/gin"
	"note_app_server/controller"
	"note_app_server/middleware"
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
		auth.Use(middleware.TokenVerificationMiddleware()).POST("/logout", func(context *gin.Context) {
			controller.Logout(context)
		})
	}

	userInfo := r.Group("/userInfo")
	userInfo.Use(middleware.TokenVerificationMiddleware())
	{
		userInfo.GET("", func(context *gin.Context) {
			controller.GetUserInfo(context)
		})
		userInfo.POST("/update", func(context *gin.Context) {
			controller.UpdateUserInfo(context)
		})
	}

	avatar := r.Group("/avatar")
	{
		avatar.GET("/:fileName", func(context *gin.Context) {
			controller.GetAvatarUrl(context)
		})
		avatar.Use(middleware.TokenVerificationMiddleware()).POST("/upload", func(context *gin.Context) {
			controller.UploadAvatar(context)
		})
	}

	return r
}
