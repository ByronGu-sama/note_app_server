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
		auth.Use(middleware.RateLimiterMiddleware(1, 5)).POST("/register", func(ctx *gin.Context) {
			controller.Register(ctx)
		})
		auth.Use(middleware.RateLimiterMiddleware(1, 5)).POST("/login", func(ctx *gin.Context) {
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
		userInfo.GET("/follows", func(ctx *gin.Context) {
			controller.GetUserFollows(ctx)
		})
		userInfo.GET("/followers", func(ctx *gin.Context) {
			controller.GetUserFollowers(ctx)
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
			verification.Use(middleware.RateLimiterMiddleware(2, 50)).PUT("", func(ctx *gin.Context) {
				controller.EditNote(ctx)
			})
			verification.Use(middleware.RateLimiterMiddleware(3, 50)).DELETE("/:nid", func(ctx *gin.Context) {
				controller.DelNote(ctx)
			})
			verification.Use(middleware.RateLimiterMiddleware(5, 75)).GET("/list", func(ctx *gin.Context) {
				controller.GetNoteList(ctx)
			})
			verification.Use(middleware.RateLimiterMiddleware(5, 75)).GET("/myNotes", func(ctx *gin.Context) {
				controller.GetMyNotes(ctx)
			})
			verification.Use(middleware.RateLimiterMiddleware(2, 15)).GET("/like/:nid", func(ctx *gin.Context) {
				controller.LikeNote(ctx)
			})
			verification.Use(middleware.RateLimiterMiddleware(2, 15)).GET("/dislike/:nid", func(ctx *gin.Context) {
				controller.DislikeNote(ctx)
			})
			verification.Use(middleware.RateLimiterMiddleware(2, 15)).GET("/collect/:nid", func(ctx *gin.Context) {
				controller.CollectNote(ctx)
			})
			verification.Use(middleware.RateLimiterMiddleware(2, 15)).GET("/cancelCollect/:nid", func(ctx *gin.Context) {
				controller.CancelCollectNote(ctx)
			})
		}

		countHeat := note.Group("").Use(middleware.TokenVerificationMiddleware()).Use(middleware.NoteTrendingMiddleware()).Use(middleware.RateLimiterMiddleware(5, 75))
		{
			countHeat.GET("/:nid", func(ctx *gin.Context) {
				controller.GetNote(ctx)
			})
		}

		checkFileType := note.Group("").Use(middleware.TokenVerificationMiddleware()).Use(middleware.DetectNotePicsTypeMiddleware())
		{
			checkFileType.POST("", func(ctx *gin.Context) {
				controller.NewNote(ctx)
			})
		}

		keywordFilter := note.Group("").Use(middleware.TokenVerificationMiddleware()).Use(middleware.RateLimiterMiddleware(5, 50))
		{
			keywordFilter.GET("/search/:keyword", func(ctx *gin.Context) {
				controller.GetNotesListWithKeyword(ctx)
			})
		}
	}

	comment := r.Group("/comment").Use(middleware.TokenVerificationMiddleware())
	{
		comment.Use(middleware.RateLimiterMiddleware(3, 30)).POST("", func(ctx *gin.Context) {
			controller.NewComment(ctx)
		})
		comment.Use(middleware.RateLimiterMiddleware(10, 60)).GET("/getList/:nid", func(ctx *gin.Context) {
			controller.GetCommentList(ctx)
		})
		comment.Use(middleware.RateLimiterMiddleware(10, 60)).GET("/getSubList/:nid/:rootId", func(ctx *gin.Context) {
			controller.GetSubCommentList(ctx)
		})
		comment.Use(middleware.RateLimiterMiddleware(5, 50)).DELETE("/delComment/:cid", func(ctx *gin.Context) {
			controller.DelComment(ctx)
		})
		comment.Use(middleware.RateLimiterMiddleware(10, 60)).GET("/like/:cid", func(ctx *gin.Context) {
			controller.LikeComment(ctx)
		})
		comment.Use(middleware.RateLimiterMiddleware(10, 60)).GET("/dislike/:cid", func(ctx *gin.Context) {
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
		checkImageType := style.Group("").Use(middleware.TokenVerificationMiddleware()).Use(middleware.DetectNormalImageTypeMiddleware())
		{
			checkImageType.POST("/updateProfileBanner", func(ctx *gin.Context) {
				controller.UpdateProfileBanner(ctx)
			})
		}
	}

	message := r.Group("/message")
	{
		message.GET("/init", func(ctx *gin.Context) {
			controller.GetWS(ctx)
		})
	}

	return r
}
