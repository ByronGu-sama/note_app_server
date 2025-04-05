package config

import (
	"github.com/redis/go-redis/v9"
	"note_app_server/global"
	"time"
)

func InitRedisConfig() {
	tokenRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.TokenDB,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	captchaRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.CaptchaDB,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	msgRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.MsgDB,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	noteTrending := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.NoteTrendingDB,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	noteBuf := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.NoteBufDB,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	thumbsUpRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.ThumbsUpRdb,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	userLikedNotesRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.UserLikedNotesRdb,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	collectedCntRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.CollectedCntRdb,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	userCollectedNotesRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.UserCollectedNotesRdb,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	global.ThumbsUpRdbClient = thumbsUpRdb
	global.UserLikedNotesRdbClient = userLikedNotesRdb
	global.TokenRdb = tokenRdb
	global.SMSCaptchaRdb = captchaRdb
	global.MsgRdb = msgRdb
	global.NoteTrendingDB = noteTrending
	global.NoteBufDB = noteBuf
	global.CollectedCntRdbClient = collectedCntRdb
	global.UserCollectedNotesRdbClient = userCollectedNotesRdb
}
