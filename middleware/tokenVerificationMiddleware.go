package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"note_app_server/global"
	"note_app_server/model/appModel"
	"note_app_server/model/userModel"
	"note_app_server/repository"
	"note_app_server/response"
	"note_app_server/service"
	"strconv"
)

func TokenVerificationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("token")
		temp, err := service.ParseJWT(tokenStr)
		// 校验token有效性
		if tokenStr == "" || err != nil {
			response.RespondWithUnauthorized(ctx, "无访问权限")
			ctx.Abort()
			return
		}
		claims := temp.(*appModel.JWT)
		// 验证uid
		uid := claims.Uid
		if uid == 0 {
			response.RespondWithUnauthorized(ctx, "无访问权限")
			ctx.Abort()
			return
		}

		rCtx := context.Background()
		_, err = global.AuthRdb.Get(rCtx, strconv.Itoa(int(uid))).Result()
		if errors.Is(err, redis.Nil) {
			response.RespondWithStatusBadRequest(ctx, "登陆已过期")
			ctx.Abort()
			return
		} else if err != nil {
			response.RespondWithStatusInternalServerError(ctx, "服务器内部错误")
			ctx.Abort()
			return
		}

		userInfo := &userModel.UserInfo{}
		if userInfo, err = repository.GetUserInfo(uid); err != nil {
			response.RespondWithUnauthorized(ctx, "无访问权限")
			ctx.Abort()
			return
		}

		// uid写入上下文
		ctx.Set("uid", uid)
		ctx.Set("username", userInfo.Username)
		ctx.Set("avatarUrl", userInfo.AvatarUrl)
		ctx.Next()
	}
}
