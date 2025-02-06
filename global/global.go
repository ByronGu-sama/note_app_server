package global

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sync"
)

var (
	Db            *gorm.DB      // Db mysql
	TokenRdb      *redis.Client // TokenRdb token redis
	CaptchaRdb    *redis.Client // CaptchaRdb 验证码redis
	OssClientPool *sync.Pool    // OssClient oss客户端连接池
	JWTKey        []byte        // JWTKey jwt密钥
)
