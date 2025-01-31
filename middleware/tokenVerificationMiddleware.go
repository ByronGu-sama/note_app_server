package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"note_app_server1/global"
	"note_app_server1/model"
	"note_app_server1/repository"
	"note_app_server1/response"
	"note_app_server1/service"
	"strconv"
)

func TokenVerificationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("token")
		temp, err := service.ParseJWT(tokenStr)
		// 校验token有效性
		if tokenStr == "" || err != nil {
			response.RespondWithUnauthorized(ctx, "无访问权限")
			return
		}
		claims := temp.(*model.JWT)
		// 验证uid
		uid := claims.Uid
		if uid == 0 {
			response.RespondWithUnauthorized(ctx, "无访问权限")
			return
		}

		rCtx := context.Background()
		_, err = global.TokenRdb.Get(rCtx, strconv.Itoa(int(uid))).Result()
		if errors.Is(err, redis.Nil) {
			response.RespondWithStatusBadRequest(ctx, "登陆已过期")
			return
		} else if err != nil {
			response.RespondWithStatusInternalServerError(ctx, "服务器内部错误")
			return
		}

		if _, err := repository.GetUserInfo(uid); err != nil {
			response.RespondWithUnauthorized(ctx, "无访问权限")
			return
		}

		// uid写入上下文
		ctx.Set("uid", uid)
		ctx.Next()
	}
}
