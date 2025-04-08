package config

import (
	"github.com/redis/go-redis/v9"
	"note_app_server/global"
	"time"
)

func InitRedisConfig() {
	authDB := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.AuthDB,
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
	boomNoteDB := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.BoomNoteDB,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	noteNormalRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.NoteNormalRdb,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	commentNormalRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.CommentNormalRdb,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	requestLimitRdb := redis.NewClient(&redis.Options{
		Addr:            AC.Redis.Host + AC.Redis.Port,
		DB:              AC.Redis.RequestLimitRdb,
		Password:        AC.Redis.Password,
		DialTimeout:     AC.Redis.Timeout * time.Millisecond,
		PoolSize:        AC.Redis.Pool.MaxActive,
		MaxIdleConns:    AC.Redis.Pool.MaxIdle,
		MinIdleConns:    AC.Redis.Pool.MinIdle,
		ConnMaxLifetime: AC.Redis.Pool.MaxWait * time.Millisecond,
	})
	global.AuthRdb = authDB
	global.MsgRdb = msgRdb
	global.BoomNoteDB = boomNoteDB
	global.NoteNormalRdb = noteNormalRdb
	global.CommentNormalRdb = commentNormalRdb
	global.RequestLimitRdb = requestLimitRdb
}
