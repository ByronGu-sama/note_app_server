package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sync"
)

var (
	Db                *gorm.DB                   // Db mysql
	TokenRdb          *redis.Client              // TokenRdb token redis
	SMSCaptchaRdb     *redis.Client              // SMSCaptchaRdb 验证码redis
	MsgRdb            *redis.Client              // MsgRdb 私聊消息缓存
	MsgHistoryRdb     *redis.Client              // MsgHistoryRdb 私聊消息历史缓存
	OssClientPool     *sync.Pool                 // OssClientPool oss客户端连接池
	CaptchaClientPool *sync.Pool                 // CaptchaClientPool 验证码客户端连接池
	ESClient          *elasticsearch.TypedClient // ESClient es客户端
	JWTKey            []byte                     // JWTKey jwt密钥
)
