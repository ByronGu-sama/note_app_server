package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"note_app_server/config"
	"note_app_server/global"
	"note_app_server/repository"
	"note_app_server/response"
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

		// 笔记热度值
		trendCnt := nid + ":Trend"
		// 笔记数据缓存
		noteBuf := nid + ":Buf"

		var incrWithTimeoutLuaScript = redis.NewScript(`
			if redis.call('EXISTS', KEYS[1]) == 0 then
				redis.call('SET', KEYS[1], 1)
				redis.call('EXPIRE', KEYS[1], ARGV[1])
				return 1
			else
				return redis.call('INCR', KEYS[1])
			end
		`)

		keys := []string{trendCnt}
		args := []interface{}{30 * 60}

		result, err := incrWithTimeoutLuaScript.Run(ctx, global.BoomNoteDB, keys, args).Result()
		if err != nil {
			for _ = range 3 {
				_, err = incrWithTimeoutLuaScript.Run(ctx, global.BoomNoteDB, keys, args).Result()
				if err == nil {
					break
				}
			}
			if err != nil {
				log.Println(err)
			}
		}
		if err == nil && result.(int64) >= int64(config.AC.App.NoteTrendingThreshold) {
			// 重置热度
			global.BoomNoteDB.Set(ctx, trendCnt, 1, 30*time.Minute)
			// 缓存笔记数据
			if _, err = global.BoomNoteDB.Get(ctx, noteBuf).Result(); err != nil {
				result, err2 := repository.GetNoteWithNid(nid)
				if err2 != nil {
					log.Fatal(err2)
				} else {
					bi, err3 := json.Marshal(result)
					if err3 != nil {
						log.Fatal(err3)
					} else {
						global.BoomNoteDB.Set(ctx, noteBuf, bi, time.Hour)
					}
				}
			} else {
				// 已有缓存时重置过期时间
				global.BoomNoteDB.Expire(ctx, noteBuf, time.Hour)
			}
		}
		ctx.Next()
	}
}
