package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"note_app_server/global"
	"time"
)

// RequestLimiter 限流器
func RequestLimiter(ctx *gin.Context, rate, capacity int) (bool, error) {
	limiterLuaScript := redis.NewScript(`
			-- rate 每秒生成的令牌数
			local rate = tonumber(ARGV[1])
			-- capacity 令牌桶容量
			local capacity = tonumber(ARGV[2])
			-- 当前的时间戳（毫秒）
			local now = tonumber(ARGV[3])
			
			-- 上次的令牌数和更新时间
			local token_data = redis.call("HMGET", KEYS[1], "tokens", "timestamp")
			local tokens = tonumber(token_data[1])
			local last_time = tonumber(token_data[2])
			
			if tokens == nil then
				tokens = capacity
				last_time = now
			end
			
			-- 计算新生成的令牌数
			local delta = math.max(0, now - last_time)
			local added_tokens = math.floor(delta * rate / 1000)
			tokens = math.min(capacity, tokens + added_tokens)
			last_time = now
						
			-- 是否通过请求
			if tokens > 0 then
				tokens = tokens - 1
				redis.call("HMSET", KEYS[1], "tokens", tokens, "timestamp", now)
				redis.call("EXPIRE", KEYS[1], 60)
				return 1
			else
				redis.call("HMSET", KEYS[1], "tokens", tokens, "timestamp", now)
				redis.call("EXPIRE", KEYS[1], 60)
				return 0
			end`)
	var uid int64
	var ip string
	var key string

	// 用请求路径区分不同的请求限流器
	pathMark := ctx.Request.URL.Path
	tempUid, ok := ctx.Get("uid")
	if !ok {
		ip = ctx.ClientIP()
		key = fmt.Sprintf("rate_limit:%s:%s", pathMark, ip)
	} else {
		uid = tempUid.(int64)
		key = fmt.Sprintf("rate_limit:%s:%d", pathMark, uid)
	}

	now := time.Now().UnixMilli()

	pass, err := limiterLuaScript.Run(ctx, global.RequestLimitRdb, []string{key}, rate, capacity, now).Int()
	return pass == 1, err
}
