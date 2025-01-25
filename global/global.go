package global

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Db         *gorm.DB      // Db mysql
	TokenRdb   *redis.Client // TokenRdb token redis
	CaptchaRdb *redis.Client // CaptchaRdb 验证码redis
	OssClient  *oss.Client   // OssClient oss客户端
)
