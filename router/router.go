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
		auth.POST("/verifyCaptcha", func(ctx *gin.Context) {
			controller.CheckCaptcha(ctx)
		})

		verification := auth.Group("").Use(middleware.TokenVerificationMiddleware())
		{
			verification.GET("/checkToken", func(ctx *gin.Context) {
				controller.CheckToken(ctx)
			})
			verification.GET("/logout", func(ctx *gin.Context) {
				controller.Logout(ctx)
			})
		}
	}

	userInfo := r.Group("/userInfo").Use(middleware.TokenVerificationMiddleware())
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
	}

	note := r.Group("/note")
	{
		note.GET("/pic/:nid/:fileName", func(ctx *gin.Context) {
			controller.GetNotePic(ctx)
		})
		verification := note.Group("").Use(middleware.TokenVerificationMiddleware())
		{
			verification.GET("/:nid", func(ctx *gin.Context) {
				controller.GetNote(ctx)
			})
			verification.PUT("", func(ctx *gin.Context) {
				controller.EditNote(ctx)
			})
			verification.DELETE("/:nid", func(ctx *gin.Context) {
				controller.DelNote(ctx)
			})
			verification.GET("/list", func(ctx *gin.Context) {
				controller.GetNoteList(ctx)
			})
			verification.GET("/myNotes", func(ctx *gin.Context) {
				controller.GetMyNotes(ctx)
			})
			verification.GET("/like/:nid", func(ctx *gin.Context) {
				controller.LikeNote(ctx)
			})
			verification.GET("/dislike/:nid", func(ctx *gin.Context) {
				controller.DislikeNote(ctx)
			})
			verification.GET("/collect/:nid", func(ctx *gin.Context) {
				controller.CollectNote(ctx)
			})
			verification.GET("/cancelCollect/:nid", func(ctx *gin.Context) {
				controller.CancelCollectNote(ctx)
			})
		}
		checkFileType := note.Group("").Use(middleware.DetectNotePicsTypeMiddleware())
		{
			checkFileType.POST("", func(ctx *gin.Context) {
				controller.NewNote(ctx)
			})
		}
	}

	comment := r.Group("/comment").Use(middleware.TokenVerificationMiddleware())
	{
		comment.POST("", func(ctx *gin.Context) {
			controller.NewComment(ctx)
		})
		comment.GET("/getList/:nid", func(ctx *gin.Context) {
			controller.GetCommentList(ctx)
		})
		comment.GET("/getSubList/:nid/:rootId", func(ctx *gin.Context) {
			controller.GetSubCommentList(ctx)
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

	style := r.Group("/style")
	{
		style.GET("/profileBanner/:bid", func(ctx *gin.Context) {
			controller.GetProfileBannerUrl(ctx)
		})
		verification := style.Group("").Use(middleware.TokenVerificationMiddleware())
		{
			verification.GET("", func(ctx *gin.Context) {
				controller.GetStyle(ctx)
			})
		}
		checkImageType := style.Group("").Use(middleware.DetectNormalImageTypeMiddleware())
		{
			checkImageType.POST("/updateProfileBanner", func(ctx *gin.Context) {
				controller.UpdateProfileBanner(ctx)
			})
		}
	}

	return r
}
