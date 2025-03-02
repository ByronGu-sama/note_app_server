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
	msgHistoryRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.MsgHistoryDB,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	global.TokenRdb = tokenRdb
	global.SMSCaptchaRdb = captchaRdb
	global.MsgRdb = msgRdb
	global.MsgHistoryRdb = msgHistoryRdb
}
