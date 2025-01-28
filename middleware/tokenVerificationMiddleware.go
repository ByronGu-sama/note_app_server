package middleware

import (
	"github.com/gin-gonic/gin"
	"note_app_server1/model"
	"note_app_server1/repository"
	"note_app_server1/response"
	"note_app_server1/service"
)

func TokenVerificationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("token")
		temp, err := service.ParseJWT(tokenStr)
		// 校验token有效性
		if tokenStr == "" || err != nil {
			response.RespondWithStatusOK(ctx, "无访问权限")
			return
		}
		claims := temp.(*model.JWT)
		// 验证uid
		uid := claims.Uid
		if uid == 0 {
			response.RespondWithStatusOK(ctx, "无访问权限")
			return
		}
		if _, err := repository.GetUserInfo(uid); err != nil {
			response.RespondWithStatusOK(ctx, "无访问权限")
			return
		}

		// uid写入上下文
		ctx.Set("uid", uid)
		ctx.Next()
	}
}
