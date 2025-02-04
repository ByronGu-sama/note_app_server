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

	note := r.Group("/note")
	note.GET("/list", func(c *gin.Context) {
		controller.GetNoteList(c)
	})
	note.Use(middleware.TokenVerificationMiddleware())
	{
		note.POST("", func(c *gin.Context) {
			controller.NewNote(c)
		})
		note.GET("/:nid", func(c *gin.Context) {
			controller.GetNote(c)
		})
		note.PUT("", func(c *gin.Context) {
			controller.EditNote(c)
		})
		note.DELETE("/:nid", func(c *gin.Context) {
			controller.DelNote(c)
		})

		note.GET("/like/:nid", func(c *gin.Context) {
			controller.LikeNote(c)
		})
		note.GET("/dislike/:nid", func(c *gin.Context) {
			controller.DislikeNote(c)
		})
		note.GET("/collect/:nid", func(c *gin.Context) {
			controller.CollectNote(c)
		})
		note.GET("/cancelCollect/:nid", func(c *gin.Context) {
			controller.CancelCollectNote(c)
		})
	}
	return r
}
