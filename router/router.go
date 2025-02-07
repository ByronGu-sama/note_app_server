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
		auth.POST("/register", func(ctx *gin.Context) {
			controller.Register(ctx)
		})
		auth.POST("/login", func(ctx *gin.Context) {
			controller.Login(ctx)
		})
		auth.Use(middleware.TokenVerificationMiddleware()).GET("/checkToken", func(ctx *gin.Context) {
			controller.CheckToken(ctx)
		})
		auth.Use(middleware.TokenVerificationMiddleware()).POST("/logout", func(ctx *gin.Context) {
			controller.Logout(ctx)
		})
	}

	userInfo := r.Group("/userInfo")
	userInfo.Use(middleware.TokenVerificationMiddleware())
	{
		userInfo.GET("", func(ctx *gin.Context) {
			controller.GetUserInfo(ctx)
		})
		userInfo.POST("/update", func(ctx *gin.Context) {
			controller.UpdateUserInfo(ctx)
		})
	}

	avatar := r.Group("/avatar")
	{
		avatar.GET("/:fileName", func(ctx *gin.Context) {
			controller.GetAvatarUrl(ctx)
		})
		avatar.Use(middleware.TokenVerificationMiddleware()).POST("/upload", func(ctx *gin.Context) {
			controller.UploadUserAvatar(ctx)
		})
		avatar.Use(middleware.TokenVerificationMiddleware()).POST("/change", func(ctx *gin.Context) {
			controller.ChangeAvatar(ctx)
		})
	}

	note := r.Group("/note")
	note.Use(middleware.TokenVerificationMiddleware())
	{
		note.POST("", func(ctx *gin.Context) {
			controller.NewNote(ctx)
		})
		note.POST("/uploadPics", func(ctx *gin.Context) {
			controller.UploadNotePics(ctx)
		})
		note.GET("/:nid", func(ctx *gin.Context) {
			controller.GetNote(ctx)
		})
		note.PUT("", func(ctx *gin.Context) {
			controller.EditNote(ctx)
		})
		note.DELETE("/:nid", func(ctx *gin.Context) {
			controller.DelNote(ctx)
		})
		note.GET("/list", func(ctx *gin.Context) {
			controller.GetNoteList(ctx)
		})
		note.GET("/pic/:nid/:fileName", func(ctx *gin.Context) {
			controller.GetNotePic(ctx)
		})
		note.GET("/myNotes", func(ctx *gin.Context) {
			controller.GetMyNotes(ctx)
		})
		note.GET("/like/:nid", func(ctx *gin.Context) {
			controller.LikeNote(ctx)
		})
		note.GET("/dislike/:nid", func(ctx *gin.Context) {
			controller.DislikeNote(ctx)
		})
		note.GET("/collect/:nid", func(ctx *gin.Context) {
			controller.CollectNote(ctx)
		})
		note.GET("/cancelCollect/:nid", func(ctx *gin.Context) {
			controller.CancelCollectNote(ctx)
		})
	}
	return r
}
