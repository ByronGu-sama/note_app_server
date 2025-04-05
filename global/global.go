package global

import (
	green20220302 "github.com/alibabacloud-go/green-20220302/v2/client"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"sync"
)

var (
	Db                          *gorm.DB                   // Db mysql
	TokenRdb                    *redis.Client              // TokenRdb token redis
	SMSCaptchaRdb               *redis.Client              // SMSCaptchaRdb 验证码redis
	MsgRdb                      *redis.Client              // MsgRdb 私聊消息缓存
	NoteTrendingDB              *redis.Client              // NoteTrendingDB 私聊消息历史缓存
	NoteBufDB                   *redis.Client              // NoteBufDB 缓存热点笔记
	ThumbsUpRdbClient           *redis.Client              // ThumbsUpRdbClient ThumbsUpRdbClient缓存点赞数据
	UserLikedNotesRdbClient     *redis.Client              // UserLikedNotesRdbClient ThumbsUpRdbClient缓存点赞数据
	OssClientPool               *sync.Pool                 // OssClientPool oss客户端连接池
	CaptchaClientPool           *sync.Pool                 // CaptchaClientPool 验证码客户端连接池
	ESClient                    *elasticsearch.TypedClient // ESClient es客户端
	JWTKey                      []byte                     // JWTKey jwt密钥
	MongoClient                 *mongo.Client              // MongoClient mongoDB客户端
	ContentCheckClient          *green20220302.Client      // ContentCheckClient 内容审核器客户端
	CollectedCntRdbClient       *redis.Client              // CollectedCntRdbClient ThumbsUpRdbClient缓存点赞数据
	UserCollectedNotesRdbClient *redis.Client              // UserCollectedNotesRdbClient ThumbsUpRdbClient缓存点赞数据
)
