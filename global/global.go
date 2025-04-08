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
	Db                 *gorm.DB                   // Db mysql
	AuthRdb            *redis.Client              // AuthRdb token redis
	MsgRdb             *redis.Client              // MsgRdb 私聊消息缓存
	BoomNoteDB         *redis.Client              // BoomNoteDB 私聊消息历史缓存
	NoteNormalRdb      *redis.Client              // NoteNormalRdb 缓存热点笔记
	CommentNormalRdb   *redis.Client              // CommentNormalRdb 缓存评论相关的附加数据
	OssClientPool      *sync.Pool                 // OssClientPool oss客户端连接池
	CaptchaClientPool  *sync.Pool                 // CaptchaClientPool 验证码客户端连接池
	ESClient           *elasticsearch.TypedClient // ESClient es客户端
	JWTKey             []byte                     // JWTKey jwt密钥
	MongoClient        *mongo.Client              // MongoClient mongoDB客户端
	ContentCheckClient *green20220302.Client      // ContentCheckClient 内容审核器客户端
)
