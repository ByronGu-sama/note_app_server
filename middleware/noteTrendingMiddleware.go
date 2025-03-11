package middleware

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"note_app_server/config"
	"note_app_server/global"
	"note_app_server/repository"
	"note_app_server/response"
	"strconv"
	"time"
)

// NoteTrendingMiddleware 笔记热度检测中间件
// 笔记访问热度超过阈值时将笔记数据缓存至redis，过期时间一小时
// 缓存的同时清空热度数据
func NoteTrendingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		nid := ctx.Param("nid")
		if nid == "" {
			response.RespondWithStatusBadRequest(ctx, "缺少关键信息")
			return
		}

		// 增加热度
		if val, err := global.NoteTrendingDB.Get(ctx, nid).Result(); errors.Is(err, redis.Nil) {
			global.NoteTrendingDB.Set(ctx, nid, 1, 30*time.Minute)
		} else {
			global.NoteTrendingDB.IncrBy(ctx, nid, 1)
			if i, err1 := strconv.Atoi(val); err1 != nil {
				log.Fatal(err1)
			} else {
				// 超过阈值
				if i > config.AC.App.NoteTrendingThreshold {
					// 重置热度
					global.NoteTrendingDB.Set(ctx, nid, 1, 30*time.Minute)
					// 缓存笔记数据
					if _, err = global.NoteBufDB.Get(ctx, nid).Result(); errors.Is(err, redis.Nil) {
						result, err2 := repository.GetNoteWithNid(nid)
						if err2 != nil {
							log.Fatal(err2)
						} else {
							bi, err3 := json.Marshal(result)
							if err3 != nil {
								log.Fatal(err3)
							} else {
								global.NoteBufDB.Set(ctx, nid, bi, time.Hour)
							}
						}
					} else {
						global.NoteBufDB.Expire(ctx, nid, time.Hour)
					}
				}
			}
		}
		ctx.Next()
	}
}
