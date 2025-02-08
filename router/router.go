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
	note.GET("/pic/:nid/:fileName", func(ctx *gin.Context) {
		controller.GetNotePic(ctx)
	})
	note.Use(middleware.TokenVerificationMiddleware())
	{
		note.POST("", func(ctx *gin.Context) {
			controller.NewNote(ctx)
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

	comment := r.Group("/comment")
	comment.Use(middleware.TokenVerificationMiddleware())
	{
		comment.POST("", func(ctx *gin.Context) {
			controller.NewComment(ctx)
		})
		comment.GET("/getList/:nid", func(ctx *gin.Context) {
			controller.GetCommentList(ctx)
		})
		comment.DELETE("/delComment/:cid", func(ctx *gin.Context) {
			controller.DelComment(ctx)
		})
		comment.GET("/like/:cid", func(ctx *gin.Context) {
			controller.LikeComment(ctx)
		})
		comment.GET("/dislike/:cid", func(ctx *gin.Context) {
			controller.CancelLikeComment(ctx)
		})
	}

	return r
}
